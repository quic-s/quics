package sync

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/quic-s/quics-protocol/pkg/types/fileinfo"
	"github.com/quic-s/quics/pkg/core/history"
	"github.com/quic-s/quics/pkg/core/registration"
	"github.com/quic-s/quics/pkg/types"
	"github.com/quic-s/quics/pkg/utils"
)

type SyncService struct {
	registrationRepository registration.Repository
	historyRepository      history.Repository
	syncRepository         Repository
}

func NewService(registrationRepository registration.Repository, historyRepository history.Repository, syncRepository Repository) *SyncService {
	return &SyncService{
		registrationRepository: registrationRepository,
		historyRepository:      historyRepository,
		syncRepository:         syncRepository,
	}
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
			RootDir:             types.RootDirectory{},
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
	file, err := ss.syncRepository.GetFileByPath(pleaseSyncReq.AfterPath)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	// TODO: event handling is needed

	file.LatestHash = pleaseSyncReq.LastUpdateHash
	file.LatestSyncTimestamp = pleaseSyncReq.LastUpdateTimestamp
	err = ss.syncRepository.UpdateFile(file)
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
}

// UpdateFileWithContents updates file (ContentExisted = true)
func (ss *SyncService) UpdateFileWithContents(pleaseTakeReq *types.PleaseTakeReq, fileInfo *fileinfo.FileInfo, fileContent io.Reader) (*types.PleaseTakeRes, error) {
	file, err := ss.syncRepository.GetFileByPath(pleaseTakeReq.AfterPath)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	// save latest file to {rootDir}
	err = saveFileToLatestDir(file.AfterPath, fileInfo, fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	// save history file to {rootDir}.history
	err = saveFileToHistoryDir(file.AfterPath, file.LatestSyncTimestamp, fileInfo, fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	// create file history entity
	fileHistory := &types.FileHistory{
		Date:       time.Now().Format("yyyy-MM-dd"),
		UUID:       pleaseTakeReq.UUID,
		BeforePath: file.BeforePath,
		AfterPath:  file.AfterPath,
		Hash:       file.LatestHash,
		File:       file.Metadata,
	}
	err = ss.historyRepository.SaveNewFileHistory(fileHistory.AfterPath, fileHistory)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	return nil, nil
}

// GetFileMetadataForMustSync returns file metadata by path
func (ss *SyncService) GetFileMetadataForMustSync(afterPath string) (*types.MustSyncReq, error) {
	file, err := ss.syncRepository.GetFileByPath(afterPath)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	mustSyncReq := &types.MustSyncReq{
		LatestHash:          file.LatestHash,
		LatestSyncTimestamp: file.LatestSyncTimestamp,
		BeforePath:          file.BeforePath,
		AfterPath:           file.AfterPath,
	}

	return mustSyncReq, nil
}

// GetFilesByRootDir returns files by root directory path
func (ss *SyncService) GetFilesByRootDir(rootDirPath string) []*types.File {
	filesByRootDir := make([]*types.File, 0)
	files := ss.syncRepository.GetAllFiles()

	for _, file := range files {
		if file.RootDir.AfterPath == rootDirPath {
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

// SyncRootDir syncs root directory to other client from owner client
func (ss *SyncService) SyncRootDir(request *types.SyncRootDirReq) error {
	client, err := ss.registrationRepository.GetClientByUUID(request.UUID)
	if err != nil {
		log.Println("quics: ", err)
	}

	path := utils.GetQuicsSyncDirPath() + request.AfterPath
	rootDir, err := ss.registrationRepository.GetRootDirByPath(path)
	if err != nil {
		log.Println("quics: ", err)
	}

	// password check
	if rootDir.Password != request.RootDirPassword {
		return errors.New("quics: (SyncRootDir) password is not correct")
	}

	rootDirs := append(client.Root, *rootDir)
	client.Root = rootDirs

	// save updated client entity with new root directory
	ss.registrationRepository.SaveClient(client.UUID, client)

	return nil
}

// CallMustSync calls must sync transaction
func (ss *SyncService) CallMustSync(pleaseTakeRes *types.PleaseTakeRes) {
	// extract root directory of this file

	// initailize connections

	// call adapter with goroutine
}

// ********************************************************************************
//                                  Private Logic
// ********************************************************************************

// FIXME: should think file with directories
func getRootDirNameAndFileName(afterPath string) (string, string) {
	requestPaths := strings.Split(afterPath, "/")
	rootDirName := requestPaths[1]
	fileName := requestPaths[2]
	return rootDirName, fileName
}

// SyncFileToLatestDir creates/updates sync file to latest directory
func saveFileToLatestDir(afterPath string, fileInfo *fileinfo.FileInfo, fileContent io.Reader) error {
	rootDirName, _ := getRootDirNameAndFileName(afterPath)
	filePath := utils.GetQuicsRootDirPath(rootDirName)

	err := fileInfo.WriteFileWithInfo(filePath, fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

// SyncFileToHistoryDir creates/updates sync file to history directory
func saveFileToHistoryDir(afterPath string, timestamp uint64, fileInfo *fileinfo.FileInfo, fileContent io.Reader) error {
	// create history directory
	rootDirName, fileName := getRootDirNameAndFileName(afterPath)
	historyDirPath := utils.GetQuicsHistoryPathByRootDir(rootDirName)
	historyFilePath := historyDirPath + strconv.FormatUint(timestamp, 10) + "_" + fileName

	err := fileInfo.WriteFileWithInfo(historyFilePath, fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}
