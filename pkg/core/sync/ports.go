package sync

import (
	"io"

	"github.com/quic-s/quics/pkg/types"
)

type Repository interface {
	SaveRootDir(afterPath string, rootDir *types.RootDirectory) error
	GetRootDirByPath(afterPath string) (*types.RootDirectory, error)
	GetAllRootDir() ([]types.RootDirectory, error)

	IsExistFileByPath(afterPath string) (bool, error)
	SaveFileByPath(afterPath string, file *types.File) error
	GetFileByPath(afterPath string) (*types.File, error)
	UpdateFile(file *types.File) error
	GetAllFiles(prefix string) ([]types.File, error)

	UpdateConflict(afterpath string, conflict *types.Conflict) error
	GetConflict(afterpath string) (*types.Conflict, error)
	GetConflictList(rootDirs []string) ([]types.Conflict, error)
	DeleteConflict(afterpath string) error

	ErrKeyNotFound() error
}

type Service interface {
	RegisterRootDir(request *types.RootDirRegisterReq) (*types.RootDirRegisterRes, error)
	SyncRootDir(request *types.RootDirRegisterReq) (*types.RootDirRegisterRes, error)
	GetRootDirList() (*types.AskRootDirRes, error)
	GetRootDirByPath(afterPath string) (*types.RootDirectory, error)

	UpdateFileWithoutContents(pleaseSyncReq *types.PleaseSyncReq) (*types.PleaseSyncRes, error)
	UpdateFileWithContents(pleaseTakeReq *types.PleaseTakeReq, fileMetadata *types.FileMetadata, fileContent io.Reader) (*types.PleaseTakeRes, error)
	CallMustSync(filePath string, UUIDs []string) error

	GetConflictList(*types.AskConflictListReq) (*types.AskConflictListRes, error)
	ChooseOne(request *types.PleaseFileReq) (*types.PleaseFileRes, error)
	CallForceSync(filePath string, UUIDs []string) error

	FullScan(uuid string) error
	BackgroundFullScan(interval uint64) error
	Rescan(*types.RescanReq) (*types.RescanRes, error)

	GetFilesByRootDir(rootDirPath string) []types.File
	GetFiles() []types.File
	GetFileByPath(afterPath string) (*types.File, error)

	RollbackFileByHistory(request *types.RollBackReq) (*types.RollBackRes, error)

	DownloadHistory(request *types.DownloadHistoryReq) (*types.DownloadHistoryRes, string, error)

	GetStagingNum(request *types.AskStagingNumReq) (*types.AskStagingNumRes, error)
	GetConflictFiles(request *types.AskStagingNumReq) ([]types.ConflictDownloadReq, error)
}

type SyncDirAdapter interface {
	SaveFileToLatestDir(afterPath string, fileMetadata *types.FileMetadata, fileContent io.Reader) error
	GetFileFromLatestDir(afterPath string) (*types.FileMetadata, io.Reader, error)
	DeleteFileFromLatestDir(afterPath string) error
	SaveFileToConflictDir(uuid string, afterPath string, fileMetadata *types.FileMetadata, fileContent io.Reader) error
	GetFileFromConflictDir(afterPath string, uuid string) (*types.FileMetadata, io.Reader, error)
	GetFileInfoFromConflictDir(afterPath string, uuid string) (*types.FileMetadata, error)
	DeleteFilesFromConflictDir(afterPath string) error
	SaveFileToHistoryDir(afterPath string, timestamp uint64, fileMetadata *types.FileMetadata, fileContent io.Reader) error
	GetFileFromHistoryDir(afterPath string, timestamp uint64) (*types.FileMetadata, io.Reader, error)
	GetFileInfoFromHistoryDir(afterPath string, timestamp uint64) (*types.FileMetadata, error)
}

type NetworkAdapter interface {
	OpenTransaction(transactionName string, uuid string) (Transaction, error)
}

type Transaction interface {
	RequestMustSync(*types.MustSyncReq) (*types.MustSyncRes, error)
	RequestGiveYou(giveYouReq *types.GiveYouReq, historyFilePath string) (*types.GiveYouRes, error)
	RequestForceSync(mustSyncReq *types.MustSyncReq, historyFilePath string) (*types.MustSyncRes, error)
	RequestAskAllMeta(askAllMetaReq *types.AskAllMetaReq) (*types.AskAllMetaRes, error)
	RequestNeedSync(needSyncReq *types.NeedSyncReq) (*types.NeedSyncRes, error)
	RequestNeedContent(needContentReq *types.NeedContentReq) (*types.NeedContentRes, *types.FileMetadata, io.Reader, error)
	Close() error
}
