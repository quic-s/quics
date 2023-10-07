package sync

import (
	"io"

	"github.com/quic-s/quics-protocol/pkg/types/fileinfo"
	"github.com/quic-s/quics/pkg/types"
)

type Repository interface {
	GetFileByPath(path string) (*types.File, error)
	SaveFileByPath(path string, file types.File) error
	GetAllFiles() []types.File
}

type Service interface {
	GetFileMetadata(path string) (*types.FileMetadata, error)
	SyncRootDir(request *types.SyncRootDirReq) error
	SyncFileToLatestDir(afterPath string, fileInfo *fileinfo.FileInfo, fileContent io.Reader) error
	SyncFileToHistoryDir(afterPath string, timestamp uint64, fileInfo *fileinfo.FileInfo, fileContent io.Reader) error
	SaveFileFromPleaseSync(path string, file types.File) error
	GetFilesByRootDir(rootDirPath string) []types.File
	GetFiles() []types.File
	GetFileByPath(path string) *types.File
}

type NetworkAdapter interface {
}
