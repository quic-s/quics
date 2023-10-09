package sync

import (
	"github.com/quic-s/quics/pkg/network/qp"
	"io"

	"github.com/quic-s/quics-protocol/pkg/types/fileinfo"
	"github.com/quic-s/quics/pkg/types"
)

type Repository interface {
	IsExistFileByPath(afterPath string) (bool, error)
	SaveFileByPath(afterPath string, file *types.File) error
	GetFileByPath(afterPath string) (*types.File, error)
	UpdateFile(file *types.File) error
	GetAllFiles() []*types.File
}

type Service interface {
	GetFileMetadataForPleaseSync(pleaseFileMetaReq *types.PleaseFileMetaReq) (*types.PleaseFileMetaRes, error)
	UpdateFileWithoutContents(pleaseSyncReq *types.PleaseSyncReq) (*types.PleaseSyncRes, error)
	UpdateFileWithContents(pleaseTakeReq *types.PleaseTakeReq, fileInfo *fileinfo.FileInfo, fileContent io.Reader) (*types.PleaseTakeRes, error)
	CallMustSync(pleaseTakeRes *types.PleaseTakeRes) error

	GetFilesByRootDir(rootDirPath string) []*types.File
	GetFiles() []*types.File
	GetFileByPath(afterPath string) (*types.File, error)
	SyncRootDir(request *types.SyncRootDirReq) error
}

type NetworkAdapter interface {
	OpenMustSyncTransaction(uuid string) (*qp.Transaction, error)
}
