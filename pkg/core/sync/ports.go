package sync

import (
	"io"

	"github.com/quic-s/quics/pkg/types"
)

type Repository interface {
	IsExistFileByPath(afterPath string) (bool, error)
	SaveFileByPath(afterPath string, file *types.File) error
	GetFileByPath(afterPath string) (*types.File, error)
	UpdateFile(file *types.File) error
	UpdateConflict(afterpath string, conflict *types.Conflict) error
	GetConflict(afterpath string) (*types.Conflict, error)
	GetConflictList(rootDirs []string) ([]types.Conflict, error)
	DeleteConflict(afterpath string) error
	GetAllFiles() []*types.File

	ErrKeyNotFound() error
}

type Service interface {
	SyncRootDir(request *types.RootDirRegisterReq) (*types.RootDirRegisterRes, error)
	UpdateFileWithoutContents(pleaseSyncReq *types.PleaseSyncReq) (*types.PleaseSyncRes, error)
	UpdateFileWithContents(pleaseTakeReq *types.PleaseTakeReq, fileMetadata *types.FileMetadata, fileContent io.Reader) (*types.PleaseTakeRes, error)
	CallMustSync(filePath string, UUIDs []string) error

	GetConflictList(*types.AskConflictListReq) (*types.AskConflictListRes, error)
	ChooseOne(request *types.PleaseFileReq) (*types.PleaseFileRes, error)
	CallForceSync(filePath string, UUIDs []string) error

	GetFilesByRootDir(rootDirPath string) []*types.File
	GetFiles() []*types.File
	GetFileByPath(afterPath string) (*types.File, error)
}

type SyncDirAdapter interface {
	SaveFileToLatestDir(afterPath string, fileMetadata *types.FileMetadata, fileContent io.Reader) error
	GetFileFromLatestDir(afterPath string) (*types.FileMetadata, io.Reader, error)
	DeleteFileFromLatestDir(afterPath string) error
	SaveFileToConflictDir(uuid string, afterPath string, fileMetadata *types.FileMetadata, fileContent io.Reader) error
	GetFileFromConflictDir(afterPath string, uuid string) (*types.FileMetadata, io.Reader, error)
	DeleteFilesFromConflictDir(afterPath string) error
	SaveFileToHistoryDir(afterPath string, timestamp uint64, fileMetadata *types.FileMetadata, fileContent io.Reader) error
	GetFileFromHistoryDir(afterPath string, timestamp uint64) (*types.FileMetadata, io.Reader, error)
}

type NetworkAdapter interface {
	OpenTransaction(transactionName string, uuid string) (Transaction, error)
}

type Transaction interface {
	RequestMustSync(*types.MustSyncReq) (*types.MustSyncRes, error)
	RequestGiveYou(giveYouReq *types.GiveYouReq, historyFilePath string) (*types.GiveYouRes, error)
	Close() error
}
