package sync

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/quic-s/quics/pkg/core/history"
	"github.com/quic-s/quics/pkg/core/registration"
	"github.com/quic-s/quics/pkg/types"
	"github.com/quic-s/quics/pkg/utils"
	"golang.org/x/exp/slices"
)

type SyncService struct {
	cancelMut              sync.RWMutex
	cancel                 map[string]context.CancelFunc
	FSTrigger              chan string
	registrationRepository registration.Repository
	historyRepository      history.Repository
	syncRepository         Repository
	networkAdapter         NetworkAdapter
	syncDirAdapter         SyncDirAdapter
}

func NewService(registrationRepository registration.Repository, historyRepository history.Repository, syncRepository Repository, networkAdapter NetworkAdapter, syncDirAdpater SyncDirAdapter) Service {
	return &SyncService{
		cancelMut:              sync.RWMutex{},
		cancel:                 map[string]context.CancelFunc{},
		FSTrigger:              make(chan string),
		registrationRepository: registrationRepository,
		historyRepository:      historyRepository,
		syncRepository:         syncRepository,
		networkAdapter:         networkAdapter,
		syncDirAdapter:         syncDirAdpater,
	}
}

// RegisterRootDir registers initial root directory to client database
func (ss *SyncService) RegisterRootDir(request *types.RootDirRegisterReq) (*types.RootDirRegisterRes, error) {
	_, err := ss.syncRepository.GetRootDirByPath(request.AfterPath)
	if err == nil {
		return nil, errors.New("root dir is already exists")
	} else if err != ss.syncRepository.ErrKeyNotFound() && err != nil {
		return nil, err
	}

	// get client entity by uuid in request data
	client, err := ss.registrationRepository.GetClientByUUID(request.UUID)
	if err != nil {
		return nil, err
	}

	UUIDs := make([]string, 0)
	UUIDs = append(UUIDs, request.UUID)

	// create root directory entity
	rootDir := &types.RootDirectory{
		BeforePath: utils.GetQuicsSyncDirPath(),
		AfterPath:  request.AfterPath,
		Owner:      client.UUID,
		Password:   request.RootDirPassword,
		UUIDs:      UUIDs,
	}
	rootDirs := append(client.Root, *rootDir)
	client.Root = rootDirs

	// save updated client entity
	err = ss.registrationRepository.SaveClient(client.UUID, client)
	if err != nil {
		return nil, err
	}

	// save requested root directory
	err = ss.syncRepository.SaveRootDir(request.AfterPath, rootDir)
	if err != nil {
		return nil, err
	}

	return &types.RootDirRegisterRes{
		UUID: request.UUID,
	}, nil
}

// SyncRootDir syncs root directory to other client from owner client
func (ss *SyncService) SyncRootDir(request *types.RootDirRegisterReq) (*types.RootDirRegisterRes, error) {
	client, err := ss.registrationRepository.GetClientByUUID(request.UUID)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	rootDir, err := ss.syncRepository.GetRootDirByPath(request.AfterPath)
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
	err = ss.syncRepository.SaveRootDir(rootDir.AfterPath, rootDir)
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

// GetRootDirList gets root directory list of client
func (ss *SyncService) GetRootDirList() (*types.AskRootDirRes, error) {
	rootDirs, err := ss.syncRepository.GetAllRootDir()
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	rootDirNames := []string{}
	for _, rootDir := range rootDirs {
		rootDirNames = append(rootDirNames, rootDir.AfterPath)
	}
	askRootDirRes := &types.AskRootDirRes{
		RootDirList: rootDirNames,
	}

	return askRootDirRes, err
}

// GetRootDirByPath gets root directory by path
func (ss *SyncService) GetRootDirByPath(path string) (*types.RootDirectory, error) {
	rootDir, err := ss.syncRepository.GetRootDirByPath(path)
	if err != nil {
		log.Println("quics: ", err)
	}

	return rootDir, nil
}

func (ss *SyncService) DisconnectRootDir(request *types.DisconnectRootDirReq) (*types.DisconnectRootDirRes, error) {
	client, err := ss.registrationRepository.GetClientByUUID(request.UUID)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}
	rootDir, err := ss.syncRepository.GetRootDirByPath(request.AfterPath)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	// find client's uuid and delete it
	for i := 0; i < len(rootDir.UUIDs); i++ {
		if client.UUID == rootDir.UUIDs[i] {
			rootDir.UUIDs = append(rootDir.UUIDs[:i], rootDir.UUIDs[i+1:]...)
			i--
		}
	}
	err = ss.syncRepository.SaveRootDir(rootDir.AfterPath, rootDir)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	// find rootDir from client's rootDir list and delete it
	for i := 0; i < len(client.Root); i++ {
		if reflect.DeepEqual(client.Root[i], rootDir) {
			client.Root = append(client.Root[:i], client.Root[i+1:]...)
			i--
		}
	}
	// save updated client entity with new root directory
	err = ss.registrationRepository.SaveClient(client.UUID, client)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	response := &types.DisconnectRootDirRes{
		UUID:      client.UUID,
		AfterPath: rootDir.AfterPath,
	}
	return response, nil
}

// UpdateFileWithoutContents updates file (ContentExisted = false)
func (ss *SyncService) UpdateFileWithoutContents(pleaseSyncReq *types.PleaseSyncReq) (*types.PleaseSyncRes, error) {
	file, err := ss.syncRepository.GetFileByPath(pleaseSyncReq.AfterPath)
	if err == ss.syncRepository.ErrKeyNotFound() {
		// If file not exist, then create the file information to database
		fileMetadata := types.FileMetadata{
			Name:    "",
			Size:    0,
			Mode:    os.FileMode(0), // FIXME: initialize with proper value
			ModTime: time.Now(),     // FIXME: initialize with proper value
			IsDir:   false,
		}

		rootDirName, _ := utils.GetNamesByAfterPath(pleaseSyncReq.AfterPath)
		file = &types.File{
			BeforePath:          "",
			AfterPath:           pleaseSyncReq.AfterPath,
			RootDirKey:          "/" + rootDirName,
			LatestHash:          "",
			LatestSyncTimestamp: 0,
			LatestEditClient:    pleaseSyncReq.UUID,
			ContentsExisted:     false,
			NeedForceSync:       false,
			Conflict:            types.Conflict{},
			Metadata:            fileMetadata,
		}
	} else if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	log.Println("quics: pleaseSyncReq: ", pleaseSyncReq)

	switch {
	// check request type is remove and file is not exist
	case pleaseSyncReq.LastUpdateHash == "" && file.LatestSyncTimestamp == 0:
		// if file is deleted then remove file from {rootDir}
		err = ss.syncDirAdapter.DeleteFileFromLatestDir(file.AfterPath)
		if err != nil && !os.IsNotExist(err) {
			log.Println("quics: ", err)
			return nil, err
		}
		// update sync file
		pleaseSyncRes := &types.PleaseSyncRes{
			UUID:      pleaseSyncReq.UUID,
			AfterPath: pleaseSyncReq.AfterPath,
		}

		return pleaseSyncRes, nil

	// check file has been updated
	case file.LatestHash == pleaseSyncReq.LastUpdateHash:
		if !file.ContentsExisted {
			// update sync file
			pleaseSyncRes := &types.PleaseSyncRes{
				UUID:      pleaseSyncReq.UUID,
				AfterPath: pleaseSyncReq.AfterPath,
			}

			return pleaseSyncRes, nil
		}
		log.Println("quics: file is already updated")
		return nil, errors.New("quics: file is already updated")

	// check file is coflict
	// conflict case LastestSyncTimestamp < LastUpdateTimestamp && LastestSyncHash == LastSyncHash
	case reflect.ValueOf(file.Conflict).IsZero() && file.LatestSyncTimestamp < pleaseSyncReq.LastUpdateTimestamp && file.LatestHash == pleaseSyncReq.LastSyncHash:
		// check event type
		if pleaseSyncReq.LastUpdateHash == "" {
			// if event type is REMOVE then set empty file metadata
			file.LatestHash = pleaseSyncReq.LastUpdateHash
			file.LatestSyncTimestamp = pleaseSyncReq.LastUpdateTimestamp
			file.LatestEditClient = pleaseSyncReq.UUID
			file.Metadata = types.FileMetadata{}
			file.ContentsExisted = false
			file.NeedForceSync = false
		} else {
			// if event type is not REMOVE then set file metadata
			file.LatestHash = pleaseSyncReq.LastUpdateHash
			file.LatestSyncTimestamp = pleaseSyncReq.LastUpdateTimestamp
			file.LatestEditClient = pleaseSyncReq.UUID
			file.Metadata = pleaseSyncReq.Metadata
			file.ContentsExisted = false
			file.NeedForceSync = false
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
		err = ss.historyRepository.SaveNewFileHistory(fileHistory.AfterPath, fileHistory)
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

	// otherwise, file is conflicted
	default:
		// handle conflict
		if reflect.ValueOf(file.Conflict).IsZero() {
			file.Conflict = types.Conflict{
				AfterPath:    file.AfterPath,
				StagingFiles: map[string]types.FileHistory{},
			}
			latestFileHistory, err := ss.historyRepository.GetFileHistory(file.AfterPath, file.LatestSyncTimestamp)
			if err != nil {
				log.Println("quics: ", err)
				return nil, err
			}

			file.Conflict.StagingFiles["server"] = *latestFileHistory
		}

		file.Conflict.StagingFiles[pleaseSyncReq.UUID] = types.FileHistory{
			Date:      time.Now().String(),
			UUID:      pleaseSyncReq.UUID,
			AfterPath: pleaseSyncReq.AfterPath,
			Timestamp: pleaseSyncReq.LastUpdateTimestamp,
			Hash:      pleaseSyncReq.LastUpdateHash,
			File:      pleaseSyncReq.Metadata,
		}

		err = ss.syncRepository.UpdateFile(file)
		if err != nil {
			return nil, err
		}
		err = ss.syncRepository.UpdateConflict(file.AfterPath, &file.Conflict)
		if err != nil {
			return nil, err
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
func (ss *SyncService) UpdateFileWithContents(pleaseTakeReq *types.PleaseTakeReq, fileMetadata *types.FileMetadata, fileContent io.Reader) (*types.PleaseTakeRes, error) {
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

		err = ss.syncDirAdapter.SaveFileToHistoryDir(file.AfterPath, file.LatestSyncTimestamp, fileMetadata, fileContent)
		if err != nil {
			log.Println("quics: ", err)
			return nil, err
		}

		// check file is deleted
		if file.LatestHash == "" {
			// if file is deleted then remove file from {rootDir}
			err = ss.syncDirAdapter.DeleteFileFromLatestDir(file.AfterPath)
			if err != nil && !os.IsNotExist(err) {
				log.Println("quics: ", err)
				return nil, err
			}
		} else {
			// check file hash is correct
			fileInfo, err := ss.syncDirAdapter.GetFileInfoFromHistoryDir(file.AfterPath, file.LatestSyncTimestamp)
			if err != nil {
				log.Println("quics: ", err)
				return nil, err
			}
			downloadedHash := utils.MakeHashFromFileMetadata(file.AfterPath, fileInfo)

			if downloadedHash != file.LatestHash {
				// if file hash is not correct then return error
				log.Println("quics: ", err)
				return nil, errors.New("quics: file hash is not correct")
			}

			// if file is not deleted then save file to {rootDir}
			fileMetadata, fileContent, err = ss.syncDirAdapter.GetFileFromHistoryDir(file.AfterPath, file.LatestSyncTimestamp)
			if err != nil {
				log.Println("quics: ", err)
				return nil, err
			}
			err = ss.syncDirAdapter.SaveFileToLatestDir(file.AfterPath, fileMetadata, fileContent)
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
			rootDir, err := ss.syncRepository.GetRootDirByPath(file.RootDirKey)
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
		err = ss.syncDirAdapter.SaveFileToConflictDir(pleaseTakeReq.UUID, file.AfterPath, fileMetadata, fileContent)
		if err != nil {
			// delete staging file info from conflict info when error occurred
			delete(file.Conflict.StagingFiles, pleaseTakeReq.UUID)
			ss.syncRepository.UpdateFile(file)
			ss.syncRepository.UpdateConflict(file.AfterPath, &file.Conflict)
			return nil, err
		}

		// check file hash is correct
		fileInfo, err := ss.syncDirAdapter.GetFileInfoFromConflictDir(file.AfterPath, pleaseTakeReq.UUID)
		if err != nil {
			log.Println("quics: ", err)
			return nil, err
		}
		downloadedHash := utils.MakeHashFromFileMetadata(file.AfterPath, fileInfo)
		if file.LatestHash != "" && downloadedHash != file.Conflict.StagingFiles[pleaseTakeReq.UUID].Hash {
			// delete staging file info from conflict info when error occurred
			delete(file.Conflict.StagingFiles, pleaseTakeReq.UUID)
			ss.syncRepository.UpdateFile(file)
			ss.syncRepository.UpdateConflict(file.AfterPath, &file.Conflict)
			return nil, errors.New("quics: file hash is not correct")
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
	ss.cancelMut.Lock()
	if _, exists := ss.cancel[filePath]; exists {
		log.Println("quics: Cancel MUSTSYNC to ", filePath)
		ss.cancel[filePath]()
	}
	ss.cancelMut.Unlock()

	ctx, cancel := context.WithCancel(context.Background())

	// use cancelMut mutex for atomic to cancel map
	ss.cancelMut.Lock()
	ss.cancel[filePath] = cancel
	ss.cancelMut.Unlock()

	defer func() {
		ss.cancelMut.Lock()
		delete(ss.cancel, filePath)
		ss.cancelMut.Unlock()
	}()

	for _, UUID := range UUIDs {
		transaction, err := ss.networkAdapter.OpenTransaction(types.MUSTSYNC, UUID)
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
			if mustSyncRes.AfterPath == "" {
				log.Println("quics: ", errors.New("quics: mustSyncRes.AfterPath is empty"))
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

func (ss *SyncService) ChooseOne(request *types.PleaseFileReq) (*types.PleaseFileRes, error) {
	client, err := ss.registrationRepository.GetClientByUUID(request.UUID)
	if err != nil {
		return nil, err
	}

	file, err := ss.syncRepository.GetFileByPath(request.AfterPath)
	if err != nil {
		return nil, err
	}
	if !slices.ContainsFunc[[]types.RootDirectory, types.RootDirectory](client.Root, func(root types.RootDirectory) bool {
		return root.AfterPath == file.RootDirKey
	}) {
		return nil, errors.New("quics: root directory is not registered")
	}

	if reflect.ValueOf(file.Conflict).IsZero() {
		return nil, errors.New("quics: file is not conflicted")
	}

	if _, exists := file.Conflict.StagingFiles[request.Side]; !exists {
		return nil, errors.New("quics: side is not exists")
	}

	// save file to {rootDir}
	if request.Side == "server" {
		fileMetadata, fileContent := &types.FileMetadata{}, io.Reader(nil)
		if file.ContentsExisted {
			fileMetadata, fileContent, err = ss.syncDirAdapter.GetFileFromHistoryDir(file.AfterPath, file.LatestSyncTimestamp)
			if err != nil {
				return nil, err
			}
		}
		// when selected side is server
		// save server file as new file to {rootDir}
		file.LatestSyncTimestamp = file.LatestSyncTimestamp + 1
		file.LatestEditClient = request.UUID
		file.NeedForceSync = true
		file.Conflict = types.Conflict{}

		// save file as new history when contents existed
		if file.ContentsExisted {
			err = ss.syncDirAdapter.SaveFileToHistoryDir(file.AfterPath, file.LatestSyncTimestamp, fileMetadata, fileContent)
			if err != nil {
				return nil, err
			}
		}

		err = ss.syncDirAdapter.DeleteFilesFromConflictDir(file.AfterPath)
		if err != nil {
			return nil, err
		}

		err = ss.syncRepository.DeleteConflict(file.AfterPath)
		if err != nil {
			return nil, err
		}

		err = ss.syncRepository.UpdateFile(file)
		if err != nil {
			return nil, err
		}
	} else {
		// when selected side is client
		// save client file as new file to {rootDir}
		selectedConflictFile := file.Conflict.StagingFiles[request.Side]
		file.LatestHash = selectedConflictFile.Hash
		file.LatestSyncTimestamp = file.LatestSyncTimestamp + 1
		file.LatestEditClient = selectedConflictFile.UUID
		file.ContentsExisted = true
		file.NeedForceSync = true
		file.Conflict = types.Conflict{}

		fileMetadata, fileContent, err := ss.syncDirAdapter.GetFileFromConflictDir(file.AfterPath, selectedConflictFile.UUID)
		if err != nil {
			return nil, err
		}

		err = ss.syncDirAdapter.SaveFileToHistoryDir(file.AfterPath, file.LatestSyncTimestamp, fileMetadata, fileContent)
		if err != nil {
			return nil, err
		}

		fileMetadata, fileContent, err = ss.syncDirAdapter.GetFileFromHistoryDir(file.AfterPath, file.LatestSyncTimestamp)
		if err != nil {
			return nil, err
		}

		err = ss.syncDirAdapter.SaveFileToLatestDir(file.AfterPath, fileMetadata, fileContent)
		if err != nil {
			return nil, err
		}

		err = ss.syncDirAdapter.DeleteFilesFromConflictDir(file.AfterPath)
		if err != nil {
			return nil, err
		}

		err = ss.syncRepository.DeleteConflict(file.AfterPath)
		if err != nil {
			return nil, err
		}

		err = ss.syncRepository.UpdateFile(file)
		if err != nil {
			return nil, err
		}
	}

	// TODO: call force sync
	// -> force sync transaction with goroutine (and end please transaction)

	if file.ContentsExisted {
		go func() {
			// extract root directory of this file
			rootDir, err := ss.syncRepository.GetRootDirByPath(file.RootDirKey)
			if err != nil {
				log.Println("quics: ", err)
				return
			}

			err = ss.CallForceSync(file.AfterPath, rootDir.UUIDs)
			if err != nil {
				log.Println("quics: ", err)
				return
			}
		}()
	}

	response := &types.PleaseFileRes{
		UUID:      request.UUID,
		AfterPath: request.AfterPath,
	}

	return response, nil
}

func (ss *SyncService) CallForceSync(filePath string, UUIDs []string) error {
	if _, exists := ss.cancel[filePath]; exists {
		log.Println("quics: Cancel FORCESYNC to ", filePath)
		ss.cancel[filePath]()
	}

	ctx, cancel := context.WithCancel(context.Background())
	ss.cancel[filePath] = cancel

	defer func() {
		delete(ss.cancel, filePath)
	}()
	for _, UUID := range UUIDs {
		transaction, err := ss.networkAdapter.OpenTransaction(types.FORCESYNC, UUID)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}
		log.Println("quics: FORCESYNC to ", UUID)

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

			historyFilePath := utils.GetHistoryFileNameByAfterPath(mustSyncReq.AfterPath, mustSyncReq.LatestSyncTimestamp)
			mustSyncRes, err := transaction.RequestForceSync(mustSyncReq, historyFilePath)
			if err != nil {
				log.Println("quics: ", err)
				return
			}
			if ctx.Err() != nil {
				log.Println("quics: ", err)
				return
			}

			if mustSyncReq.LatestHash != mustSyncRes.LatestSyncHash {
				log.Println("quics: hash is not correct; fail to send file")
			}
		}()
	}
	return nil
}

func (ss *SyncService) FullScan(uuid string) error {
	client, err := ss.registrationRepository.GetClientByUUID(uuid)
	if err != nil {
		return err
	}

	transaction, err := ss.networkAdapter.OpenTransaction(types.FULLSCAN, uuid)
	if err != nil {
		return err
	}

	askAllMetaReq := &types.AskAllMetaReq{
		UUID: uuid,
	}

	askAllMetaRes, err := transaction.RequestAskAllMeta(askAllMetaReq)
	if err != nil {
		return err
	}

	if askAllMetaRes.UUID != uuid {
		return errors.New("quics: UUID is not equal")
	}

	for _, rootDir := range client.Root {
		allFiles, err := ss.syncRepository.GetAllFiles(rootDir.AfterPath)
		if err != nil {
			return err
		}
		for i, file := range allFiles {
			if !file.ContentsExisted && file.LatestEditClient == uuid {
				err := ss.CallNeedContent(&allFiles[i])
				if err != nil {
					log.Println("quics: ", err)
				}
			}
			if !reflect.ValueOf(file.Conflict).IsZero() {
				continue
			}

			exist := false
			for _, clientFile := range askAllMetaRes.SyncMetaList {
				if file.AfterPath == clientFile.AfterPath {
					exist = true
					if clientFile.LastUpdateTimestamp == clientFile.LastSyncTimestamp && file.LatestSyncTimestamp > clientFile.LastUpdateTimestamp {
						// need must synce
						if file.NeedForceSync {
							err = ss.CallForceSync(file.AfterPath, []string{uuid})
							if err != nil {
								log.Println("quics: ", err)
								break
							}
						} else {
							err = ss.CallMustSync(file.AfterPath, []string{uuid})
							if err != nil {
								log.Println("quics: ", err)
								break
							}
						}
					}
					break
				}
			}
			if !exist {
				// need must sync
				if file.NeedForceSync {
					err = ss.CallForceSync(file.AfterPath, []string{uuid})
					if err != nil {
						log.Println("quics: ", err)
						continue
					}
				} else {
					err = ss.CallMustSync(file.AfterPath, []string{uuid})
					if err != nil {
						log.Println("quics: ", err)
						continue
					}
				}
			}

		}
	}

	return nil
}

func (ss *SyncService) BackgroundFullScan(secInterval uint64) error {
	go func() {
		for {
			time.Sleep(time.Duration(secInterval) * time.Second)
			ss.FSTrigger <- "all"
		}
	}()
	go func() {
		for {
			uuid := <-ss.FSTrigger
			if uuid == "all" {
				clients, err := ss.registrationRepository.GetAllClients()
				if err != nil {
					log.Println("quics: ", err)
					continue
				}

				for _, client := range clients {
					err = ss.FullScan(client.UUID)
					if err != nil {
						log.Println("quics: ", err)
						continue
					}
				}
			} else {
				err := ss.FullScan(uuid)
				if err != nil {
					log.Println("quics: ", err)
				}
			}
		}
	}()
	return nil
}

func (ss *SyncService) Rescan(request *types.RescanReq) (*types.RescanRes, error) {
	ss.FSTrigger <- request.UUID
	rescanRes := &types.RescanRes{
		UUID: request.UUID,
	}
	return rescanRes, nil
}

func (ss *SyncService) CallNeedContent(file *types.File) error {
	if file.ContentsExisted {
		return errors.New("quics: file contents is already existed")
	}

	transaction, err := ss.networkAdapter.OpenTransaction(types.NEEDCONTENT, file.LatestEditClient)
	if err != nil {
		return err
	}

	needContentReq := &types.NeedContentReq{
		UUID:                file.LatestEditClient,
		AfterPath:           file.AfterPath,
		LastUpdateTimestamp: file.LatestSyncTimestamp,
		LastUpdateHash:      file.LatestHash,
	}

	res, fileMetadata, fileContent, err := transaction.RequestNeedContent(needContentReq)
	if err != nil {
		return err
	}

	if res.UUID != file.LatestEditClient {
		return errors.New("quics: UUID is not equal")
	}
	if res.AfterPath != file.AfterPath {
		return errors.New("quics: AfterPath is not equal")
	}
	if res.LastUpdateTimestamp != file.LatestSyncTimestamp {
		return errors.New("quics: LastUpdateTimestamp is not equal")
	}
	if res.LastUpdateHash != file.LatestHash {
		return errors.New("quics: LastUpdateHash is not equal")
	}
	if fileMetadata == nil {
		return errors.New("quics: fileMetadata is nil")
	}
	if fileContent == nil {
		return errors.New("quics: fileContent is nil")
	}

	// save file to history dir
	err = ss.syncDirAdapter.SaveFileToHistoryDir(file.AfterPath, file.LatestSyncTimestamp, fileMetadata, fileContent)
	if err != nil {
		return err
	}

	// copy file to latest dir
	fileMetadata, fileContent, err = ss.syncDirAdapter.GetFileFromHistoryDir(file.AfterPath, file.LatestSyncTimestamp)
	if err != nil {
		return err
	}

	err = ss.syncDirAdapter.SaveFileToLatestDir(file.AfterPath, fileMetadata, fileContent)
	if err != nil {
		return err
	}

	// update file
	file.ContentsExisted = true
	err = ss.syncRepository.UpdateFile(file)
	if err != nil {
		return err
	}

	return nil
}

// GetFilesByRootDir returns files by root directory path
func (ss *SyncService) GetFilesByRootDir(rootDirPath string) []types.File {
	files, err := ss.syncRepository.GetAllFiles(rootDirPath)
	if err != nil {
		log.Println("quics: ", err)
		return nil
	}

	return files
}

// GetFiles returns all files in database
func (ss *SyncService) GetFiles() []types.File {

	files, err := ss.syncRepository.GetAllFiles("")
	if err != nil {
		log.Println("quics: ", err)
		return nil
	}
	return files
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

func (ss *SyncService) RollbackFileByHistory(request *types.RollBackReq) (*types.RollBackRes, error) {
	fileData, err := ss.syncRepository.GetFileByPath(request.AfterPath)
	if err != nil {
		return nil, err
	}

	historyData, err := ss.historyRepository.GetFileHistory(request.AfterPath, request.Version)
	if err != nil {
		return nil, err
	}

	newHistoryData := &types.FileHistory{
		Date:       time.Now().String(),
		UUID:       request.UUID,
		BeforePath: historyData.BeforePath,
		AfterPath:  historyData.AfterPath,
		Timestamp:  fileData.LatestSyncTimestamp + 1,
		Hash:       historyData.Hash,
		File:       historyData.File,
	}
	err = ss.historyRepository.SaveNewFileHistory(request.AfterPath, newHistoryData)
	if err != nil {
		return nil, err
	}

	newFileData := &types.File{
		BeforePath:          utils.GetQuicsSyncDirPath(),
		AfterPath:           fileData.AfterPath,
		RootDirKey:          fileData.RootDirKey,
		LatestHash:          newHistoryData.Hash,
		LatestSyncTimestamp: newHistoryData.Timestamp,
		LatestEditClient:    request.UUID,
		ContentsExisted:     true,
		NeedForceSync:       false,
		Metadata:            newHistoryData.File,
	}
	err = ss.syncRepository.SaveFileByPath(newFileData.AfterPath, newFileData)
	if err != nil {
		return nil, err
	}

	historyFileMetadata, historyFileInfo, err := ss.syncDirAdapter.GetFileFromHistoryDir(historyData.AfterPath, historyData.Timestamp)
	if err != nil {
		return nil, err
	}

	err = ss.syncDirAdapter.SaveFileToHistoryDir(newHistoryData.AfterPath, newHistoryData.Timestamp, historyFileMetadata, historyFileInfo)
	if err != nil {
		return nil, err
	}

	fileMetadata, fileInfo, err := ss.syncDirAdapter.GetFileFromHistoryDir(newHistoryData.AfterPath, newHistoryData.Timestamp)
	if err != nil {
		return nil, err
	}

	err = ss.syncDirAdapter.SaveFileToLatestDir(newFileData.AfterPath, fileMetadata, fileInfo)
	if err != nil {
		return nil, err
	}

	// call must sync
	rootDir, err := ss.syncRepository.GetRootDirByPath(newFileData.RootDirKey)
	if err != nil {
		return nil, err
	}

	UUIDs := rootDir.UUIDs

	err = ss.CallMustSync(newFileData.AfterPath, UUIDs)
	if err != nil {
		return nil, err
	}

	return &types.RollBackRes{
		UUID: request.UUID,
	}, nil
}

func (ss *SyncService) GetStagingNum(request *types.AskStagingNumReq) (*types.AskStagingNumRes, error) {
	// get file by afterPath
	file, err := ss.syncRepository.GetFileByPath(request.AfterPath)
	if err != nil {
		return nil, err
	}

	return &types.AskStagingNumRes{
		UUID:        request.UUID,
		ConflictNum: uint64(len(file.Conflict.StagingFiles)),
	}, nil
}

func (ss *SyncService) GetConflictFiles(request *types.AskStagingNumReq) ([]types.ConflictDownloadReq, error) {
	// get file by afterPath
	conflict, err := ss.syncRepository.GetConflict(request.AfterPath)
	if err != nil {
		return nil, err
	}

	response := []types.ConflictDownloadReq{}

	for uuid, stagingFile := range conflict.StagingFiles {

		conflictDownloadReq := types.ConflictDownloadReq{
			UUID:      request.UUID,
			Candidate: uuid,
			AfterPath: stagingFile.AfterPath,
		}

		response = append(response, conflictDownloadReq)
	}

	return response, nil
}

func (ss *SyncService) DownloadHistory(request *types.DownloadHistoryReq) (*types.DownloadHistoryRes, string, error) {
	history, err := ss.historyRepository.GetFileHistory(request.AfterPath, request.Version)
	if err != nil {
		return nil, "", err
	}

	file, err := ss.syncRepository.GetFileByPath(request.AfterPath)
	if err != nil {
		return nil, "", err
	}

	historyFileName := utils.ExtractFileNameFromHistoryFile(history.AfterPath)
	filePath := utils.GetQuicsHistoryPathByRootDir(file.RootDirKey) + "/" + historyFileName + "_" + fmt.Sprint(request.Version)

	return &types.DownloadHistoryRes{
		UUID: request.UUID,
	}, filePath, nil
}
