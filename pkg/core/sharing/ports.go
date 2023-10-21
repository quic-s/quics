package sharing

import "github.com/quic-s/quics/pkg/types"

type Repository interface {
	SaveLink(sharing *types.Sharing) error
	GetLink(link string) (*types.Sharing, error)
	DeleteLink(link string) error
}

type Service interface {
	CreateLink(request *types.ShareReq) (*types.ShareRes, error)
	DeleteLink(request *types.StopShareReq) (*types.StopShareRes, error)
}
