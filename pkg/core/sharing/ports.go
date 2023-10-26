package sharing

import (
	"io"

	"github.com/quic-s/quics/pkg/types"
)

type Repository interface {
	SaveLink(sharing *types.Sharing) error
	GetLink(link string) (*types.Sharing, error)
	DeleteLink(link string) error
	UpdateLink(*types.Sharing) error
}

type Service interface {
	CreateLink(request *types.ShareReq) (*types.ShareRes, error)
	DeleteLink(request *types.StopShareReq) (*types.StopShareRes, error)
	DownloadFile(uuid string, afterPath string) (*types.FileMetadata, io.Reader, error)
}

type SyncDirAdapter interface {
	GetFileFromHistoryDir(afterPath string, timestamp uint64) (*types.FileMetadata, io.Reader, error)
}
