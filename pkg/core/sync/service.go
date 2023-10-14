package sync

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/quic-s/quics-protocol/pkg/types/fileinfo"
	"github.com/quic-s/quics/pkg/core/history"
	"github.com/quic-s/quics/pkg/core/registration"
	"github.com/quic-s/quics/pkg/types"
	"github.com/quic-s/quics/pkg/utils"
	"golang.org/x/exp/slices"
)

type SyncService struct {
	cancel                 map[string]context.CancelFunc
	registrationRepository registration.Repository
	historyRepository      history.Repository
	syncRepository         Repository
	networkAdapter         NetworkAdapter
	syncDirAdapter         SyncDirAdapter
}

func NewService(registrationRepository registration.Repository, historyRepository history.Repository, syncRepository Repository, networkAdapter NetworkAdapter, syncDirAdpater SyncDirAdapter) *SyncService {
	cancel := make(map[string]context.CancelFunc)
	return &SyncService{
		cancel:                 cancel,
		registrationRepository: registrationRepository,
		historyRepository:      historyRepository,
		syncRepository:         syncRepository,
		networkAdapter:         networkAdapter,
		syncDirAdapter:         syncDirAdpater,
	}
}

// SyncRootDir syncs root directory to other client from owner client
func (ss *SyncService) SyncRootDir(request *types.RootDirRegisterReq) (*types.RootDirRegisterRes, error) {
	client, err := ss.registrationRepository.GetClientByUUID(request.UUID)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	rootDir, err := ss.registrationRepository.GetRootDirByPath(request.AfterPath)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	// password check
	if rootDir.Password != request.RootDirPassword {
		return nil, errors.New("quics: Root directory password is not correct")
	}

	if !slices.Contains[[]string, string](rootDir.UUIDs, client.UUID) {
		// add client UUID to root directory
		rootDir.UUIDs = append(rootDir.UUIDs, client.UUID)
	}
	err = ss.registrationRepository.SaveRootDir(rootDir.AfterPath, rootDir)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	if !slices.ContainsFunc[[]types.RootDirectory, types.RootDirectory](client.Root, func(clientRoot types.RootDirectory) bool {
		return reflect.DeepEqual(rootDir, clientRoot)
	}) {
		// add root directory to client
		client.Root = append(client.Root, *rootDir)
	}

	// save updated client entity with new root directory
	err = ss.registrationRepository.SaveClient(client.UUID, client)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	response := &types.RootDirRegisterRes{
		UUID: request.UUID,
	}
	return response, nil
}

// GetFileMetadataForPleaseSync returns file metadata by path
func (ss *SyncService) GetFileMetadataForPleaseSync(pleaseFileMetaReq *types.PleaseFileMetaReq) (*types.PleaseFileMetaRes, error) {
	afterPath := pleaseFileMetaReq.AfterPath

	isExistFile, err := ss.syncRepository.IsExistFileByPath(afterPath)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	if !isExistFile {
		// If file not exist, then create the file information to database
		fileMetadata := types.FileMetadata{
			Name:    "",
			Size:    0,
			Mode:    os.FileMode(0), // FIXME: initialize with proper value
			ModTime: time.Now(),     // FIXME: initialize with proper value
			IsDir:   false,
		}

		file := &types.File{
			BeforePath:          "",
			AfterPath:           afterPath,
			RootDirKey:          "",
			LatestHash:          "",
			LatestSyncTimestamp: 0,
			ContentsExisted:     false,
			Metadata:            fileMetadata,
		}

		err := ss.syncRepository.SaveFileByPath(afterPath, file)
		if err != nil {
			log.Println("quics: ", err)
			return nil, err
		}
	}

	file, err := ss.syncRepository.GetFileByPath(afterPath)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	pleaseFileMetaRes := &types.PleaseFileMetaRes{
		UUID:                pleaseFileMetaReq.UUID,
		AfterPath:           pleaseFileMetaReq.AfterPath,
		LatestHash:          file.LatestHash,
		LatestSyncTimestamp: file.LatestSyncTimestamp,
		ModifiedDate:        file.Metadata.ModTime.Local().Format("yyyy-MM-dd"),
	}

	return pleaseFileMetaRes, nil
}

// UpdateFileWithoutContents updates file (ContentExisted = false)
func (ss *SyncService) UpdateFileWithoutContents(pleaseSyncReq *types.PleaseSyncReq) (*types.PleaseSyncRes, error) {
	isExistFile, _ := ss.syncRepository.IsExistFileByPath(pleaseSyncReq.AfterPath)

	if !isExistFile {
		// If file not exist, then create the file information to database
		fileMetadata := types.FileMetadata{
			Name:    "",
			Size:    0,
			Mode:    os.FileMode(0), // FIXME: initialize with proper value
			ModTime: time.Now(),     // FIXME: initialize with proper value
			IsDir:   false,
		}
		rootDirName, _ := utils.GetNamesByAfterPath(pleaseSyncReq.AfterPath)
		rootDirKey := utils.GetQuicsRootDirPath(rootDirName)

		file := &types.File{
			BeforePath:          "",
			AfterPath:           pleaseSyncReq.AfterPath,
			RootDirKey:          rootDirKey,
			LatestHash:          "",
			LatestSyncTimestamp: 0,
			LatestEditClient:    pleaseSyncReq.UUID,
			Conflict:            types.Conflict{},
			ContentsExisted:     false,
			Metadata:            fileMetadata,
		}

		err := ss.syncRepository.SaveFileByPath(pleaseSyncReq.AfterPath, file)
		if err != nil {
			log.Println("quics: ", err)
			return nil, err
		}
	}

	file, err := ss.syncRepository.GetFileByPath(pleaseSyncReq.AfterPath)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	log.Println("quics: pleaseSyncReq: ", pleaseSyncReq)
	// check file is already conflicted
	if !reflect.ValueOf(file.Conflict).IsZero() {
		log.Println("quics: file is already conflicted")
		return nil, errors.New("quics: file is already conflicted")
	}

	if file.LatestHash == pleaseSyncReq.LastUpdateHash {
		log.Println("quics: file is already updated")
		return nil, errors.New("quics: file is already updated")
	}

	// check file is coflict
	// conflict case LastestSyncTimestamp < LastUpdateTimestamp && LastestSyncHash == LastSyncHash
	if reflect.ValueOf(file.Conflict).IsZero() && file.LatestSyncTimestamp < pleaseSyncReq.LastUpdateTimestamp && file.LatestHash == pleaseSyncReq.LastSyncHash {
		// check event type
		if pleaseSyncReq.Event == "REMOVE" {
			// if event type is REMOVE then set empty file metadata
			if pleaseSyncReq.LastUpdateHash != "" {
				return nil, errors.New("quics: LastUpdateHash must be empty")
			}

			file.LatestHash = pleaseSyncReq.LastUpdateHash
			file.LatestSyncTimestamp = pleaseSyncReq.LastUpdateTimestamp
			file.LatestEditClient = pleaseSyncReq.UUID
			file.Metadata = types.FileMetadata{}
			file.ContentsExisted = false
		} else {
			// if event type is not REMOVE then set file metadata
			file.LatestHash = pleaseSyncReq.LastUpdateHash
			file.LatestSyncTimestamp = pleaseSyncReq.LastUpdateTimestamp
			file.LatestEditClient = pleaseSyncReq.UUID
			file.Metadata = pleaseSyncReq.Metadata
			file.ContentsExisted = false
		}

		err = ss.syncRepository.UpdateFile(file)
		if err != nil {
			log.Println("quics: ", err)
			return nil, err
		}

		// create file history entity
		fileHistory := &types.FileHistory{
			Date:       time.Now().String(),
			UUID:       pleaseSyncReq.UUID,
			BeforePath: file.BeforePath,
			AfterPath:  file.AfterPath,
			Timestamp:  file.LatestSyncTimestamp,
			Hash:       file.LatestHash,
			File:       file.Metadata,
		}
		err = ss.historyRepository.SaveNewFileHistory(utils.GetHistoryFileNameByAfterPath(fileHistory.AfterPath, fileHistory.Timestamp), fileHistory)
		if err != nil {
			log.Println("quics: ", err)
			return nil, err
		}

		// update sync file
		pleaseSyncRes := &types.PleaseSyncRes{
			UUID:      pleaseSyncReq.UUID,
			AfterPath: pleaseSyncReq.AfterPath,
		}

		return pleaseSyncRes, nil
	} else {
		// TODO: handle conflict
		if reflect.ValueOf(file.Conflict).IsZero() {
			file.Conflict = types.Conflict{
				AfterPath:    file.AfterPath,
				StagingFiles: make(map[string]types.FileHistory),
			}
			err = ss.syncRepository.UpdateFile(file)
			if err != nil {
				return nil, err
			}
			err = ss.syncRepository.UpdateConflict(file.AfterPath, &file.Conflict)
			if err != nil {
				return nil, err
			}
		}

		// update sync file
		pleaseSyncRes := &types.PleaseSyncRes{
			UUID:      pleaseSyncReq.UUID,
			AfterPath: pleaseSyncReq.AfterPath,
		}
		return pleaseSyncRes, nil
	}
}

// UpdateFileWithContents updates file (ContentExisted = true)
func (ss *SyncService) UpdateFileWithContents(pleaseTakeReq *types.PleaseTakeReq, fileInfo *fileinfo.FileInfo, fileContent io.Reader) (*types.PleaseTakeRes, error) {
	file, err := ss.syncRepository.GetFileByPath(pleaseTakeReq.AfterPath)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	log.Println("quics: pleaseTakeReq: ", pleaseTakeReq)
	// check file is coflicted
	if reflect.ValueOf(file.Conflict).IsZero() {
		// if file is not conflicted then update file
		// save latest file to {rootDir}

		err = ss.syncDirAdapter.SaveFileToHistoryDir(file.AfterPath, file.LatestSyncTimestamp, fileInfo, fileContent)
		if err != nil {
			log.Println("quics: ", err)
			return nil, err
		}

		// check file is deleted
		if file.LatestHash == "" {
			// if file is deleted then remove file from {rootDir}
			ss.syncDirAdapter.DeleteFileFromLatestDir(file.AfterPath)
		} else {
			// if file is not deleted then save file to {rootDir}
			err = ss.syncDirAdapter.CopyHistoryFileToLatestDir(file.AfterPath, file.LatestSyncTimestamp, fileInfo)
			if err != nil {
				log.Println("quics: ", err)
				return nil, err
			}

		}

		file.ContentsExisted = true
		err = ss.syncRepository.UpdateFile(file)
		if err != nil {
			log.Println("quics: ", err)
			return nil, err
		}

		// TODO: call must sync
		// -> must sync transaction with goroutine (and end please transaction)

		go func() {
			// extract root directory of this file
			rootDirName, _ := utils.GetNamesByAfterPath(file.AfterPath)
			rootDir, err := ss.registrationRepository.GetRootDirByPath("/" + rootDirName)
			if err != nil {
				log.Println("quics: ", err)
				return
			}

			// remove the UUID of this client that requested previous please sync
			UUIDs := rootDir.UUIDs
			for i, UUID := range UUIDs {
				if UUID == pleaseTakeReq.UUID {
					UUIDs = append(UUIDs[:i], UUIDs[i+1:]...)
				}
			}

			err = ss.CallMustSync(file.AfterPath, UUIDs)
			if err != nil {
				log.Println("quics: ", err)
				return
			}
		}()

		// <- must sync transaction with goroutine (and end please transaction)

		pleaseTakeRes := &types.PleaseTakeRes{
			UUID:      pleaseTakeReq.UUID,
			AfterPath: pleaseTakeReq.AfterPath,
		}
		return pleaseTakeRes, nil
	} else {
		// if file is conflicted then save file to {rootDir}.conflict

		// save file to {rootDir}.conflict
		err = ss.syncDirAdapter.SaveFileToConflictDir(file.AfterPath, fileInfo, fileContent)
		if err != nil {
			log.Println("quics: ", err)
			return nil, err
		}

		file.Conflict.StagingFiles[pleaseTakeReq.UUID] = types.FileHistory{
			Date:       time.Now().String(),
			UUID:       pleaseTakeReq.UUID,
			BeforePath: file.BeforePath,
			AfterPath:  file.AfterPath,
			Timestamp:  file.LatestSyncTimestamp,
			Hash:       file.LatestHash,
			File:       file.Metadata,
		}

		err = ss.syncRepository.UpdateFile(file)
		if err != nil {
			log.Println("quics: ", err)
			return nil, err
		}

		err = ss.syncRepository.UpdateConflict(file.AfterPath, &file.Conflict)
		if err != nil {
			return nil, err
		}

		// update sync file
		pleaseTakeRes := &types.PleaseTakeRes{
			UUID:      pleaseTakeReq.UUID,
			AfterPath: pleaseTakeReq.AfterPath,
		}

		return pleaseTakeRes, nil
	}
}

// CallMustSync calls must sync transaction
func (ss *SyncService) CallMustSync(filePath string, UUIDs []string) error {
	if _, exists := ss.cancel[filePath]; exists {
		log.Println("quics: Cancel MUSTSYNC to ", filePath)
		ss.cancel[filePath]()
	}

	ctx, cancel := context.WithCancel(context.Background())
	ss.cancel[filePath] = cancel

	defer func() {
		delete(ss.cancel, filePath)
	}()
	for _, UUID := range UUIDs {
		transaction, err := ss.networkAdapter.OpenMustSyncTransaction(UUID)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}
		log.Println("quics: MUSTSYNC to ", UUID)

		// -> must sync

		go func() {
			defer func() {
				err = transaction.Close()
				if err != nil {
					log.Println("quics: ", err)
					return
				}
			}()
			file, err := ss.syncRepository.GetFileByPath(filePath)
			if err != nil {
				log.Println("quics: ", err)
				return
			}
			if ctx.Err() != nil {
				log.Println("quics: ", err)
				return
			}

			mustSyncReq := &types.MustSyncReq{
				LatestHash:          file.LatestHash,
				LatestSyncTimestamp: file.LatestSyncTimestamp,
				BeforePath:          file.BeforePath,
				AfterPath:           file.AfterPath,
			}
			if ctx.Err() != nil {
				log.Println("quics: ", err)
				return
			}

			mustSyncRes, err := transaction.RequestMustSync(mustSyncReq)
			if err != nil {
				log.Println("quics: ", err)
				return
			}
			if ctx.Err() != nil {
				log.Println("quics: ", err)
				return
			}

			// <- must sync

			// -> give file

			giveYouReq := &types.GiveYouReq{
				UUID:      mustSyncRes.UUID,
				AfterPath: mustSyncRes.AfterPath,
			}
			if ctx.Err() != nil {
				log.Println("quics: ", err)
				return
			}

			historyFilePath := utils.GetHistoryFileNameByAfterPath(mustSyncRes.AfterPath, mustSyncRes.LatestSyncTimestamp)
			giveYouRes, err := transaction.RequestGiveYou(giveYouReq, historyFilePath)
			if err != nil {
				log.Println("quics: ", err)
				return
			}
			if ctx.Err() != nil {
				log.Println("quics: ", err)
				return
			}

			file, err = ss.syncRepository.GetFileByPath(giveYouRes.AfterPath)
			if err != nil {
				log.Println("quics: ", err)
				return
			}
			if ctx.Err() != nil {
				log.Println("quics: ", err)
				return
			}

			err = validateGiveYouTransaction(file, giveYouRes)
			if err != nil {
				log.Println("quics: ", err)
				return
			}
			// <- give file
		}()
	}
	return nil
}

func (ss *SyncService) GetConflictList(request *types.AskConflictListReq) (*types.AskConflictListRes, error) {
	client, err := ss.registrationRepository.GetClientByUUID(request.UUID)
	if err != nil {
		return nil, err
	}

	rootDirs := []string{}
	for _, r := range client.Root {
		rootDirs = append(rootDirs, r.AfterPath)
	}
	conflicts, err := ss.syncRepository.GetConflictList(rootDirs)
	if err != nil {
		return nil, err
	}

	return &types.AskConflictListRes{
		UUID:      request.UUID,
		Conflicts: conflicts,
	}, nil
}

// GetFilesByRootDir returns files by root directory path
func (ss *SyncService) GetFilesByRootDir(rootDirPath string) []*types.File {
	filesByRootDir := make([]*types.File, 0)
	files := ss.syncRepository.GetAllFiles()

	for _, file := range files {
		if file.RootDirKey == rootDirPath {
			filesByRootDir = append(files, file)
		}
	}

	return filesByRootDir
}

// GetFiles returns all files in database
func (ss *SyncService) GetFiles() []*types.File {
	return ss.syncRepository.GetAllFiles()
}

// GetFileByPath returns file entity by path
func (ss *SyncService) GetFileByPath(path string) (*types.File, error) {
	return ss.syncRepository.GetFileByPath(path)
}

// ********************************************************************************
//                                  Private Logic
// ********************************************************************************

func validateGiveYouTransaction(file *types.File, giveYouRes *types.GiveYouRes) error {
	if file.LatestSyncTimestamp != giveYouRes.LastSyncTimestamp && file.LatestHash != giveYouRes.LastHash {
		err := errors.New("not equals hash and timestamp")
		if err != nil {
			return err
		}
	}

	if file.LatestSyncTimestamp != giveYouRes.LastSyncTimestamp {
		err := errors.New("not equals timestamp")
		if err != nil {
			return err
		}
	}

	if file.LatestHash != giveYouRes.LastHash {
		err := errors.New("not equals hash")
		if err != nil {
			return err
		}
	}

	return nil
}
