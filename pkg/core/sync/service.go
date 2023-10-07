package sync

import (
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

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

func (ss *SyncService) GetFileMetadata(pleaseFileMetaReq *types.PleaseFileMetaReq) (*types.PleaseFileMetaRes, error) {
	file, err := ss.syncRepository.GetFileByPath(pleaseFileMetaReq.AfterPath)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	pleaseFileMetaRes := &types.PleaseFileMetaRes{
		UUID:                pleaseFileMetaReq.UUID,
		AfterPath:           pleaseFileMetaReq.AfterPath,
		LatestHash:          file.LatestHash,
		LatestSyncTimestamp: file.LatestSyncTimestamp,
		ModifiedDate:        uint64(file.Metadata.ModTime.UnixNano()),
	}

	return pleaseFileMetaRes, nil
}

// SaveFileFromPleaseSync saves file from please sync request
func (ss *SyncService) SaveFileFromPleaseSync(path string, file types.File) error {
	err := ss.syncRepository.SaveFileByPath(path, file)
	if err != nil {
		return err
	}

	return nil
}

// GetFilesByRootDir returns files by root directory path
func (ss *SyncService) GetFilesByRootDir(rootDirPath string) []types.File {
	filesByRootDir := make([]types.File, 0)
	files := ss.syncRepository.GetAllFiles()

	for _, file := range files {
		if file.RootDir.AfterPath == rootDirPath {
			filesByRootDir = append(files, file)
		}
	}

	return filesByRootDir
}

// GetFiles returns all files in database
func (ss *SyncService) GetFiles() []types.File {
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

// SyncFileToLatestDir creates/updates sync file to latest directory
func (ss *SyncService) SyncFileToLatestDir(afterPath string, fileInfo *fileinfo.FileInfo, fileContent io.Reader) error {
	rootDirName, _ := getRootDirNameAndFileName(afterPath)
	filePath := utils.GetQuicsLatestPathByRootDir(rootDirName)

	err := fileInfo.WriteFileWithInfo(filePath, fileContent)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

// SyncFileToHistoryDir creates/updates sync file to history directory
func (ss *SyncService) SyncFileToHistoryDir(afterPath string, timestamp uint64, fileInfo *fileinfo.FileInfo, fileContent io.Reader) error {
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
