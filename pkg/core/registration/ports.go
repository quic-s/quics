package registration

import (
	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/types"
)

type Repository interface {
	SaveClient(uuid string, client *types.Client) error
	GetClientByUUID(uuid string) (*types.Client, error)
	SaveRootDir(path string, rootDir *types.RootDirectory) error
	GetRootDirByPath(path string) (*types.RootDirectory, error)
	GetAllRootDir() ([]*types.RootDirectory, error)
	GetSequence(key []byte, increment uint64) (uint64, error)
}

type Service interface {
	RegisterClient(request *types.ClientRegisterReq, conn *qp.Connection) (*types.ClientRegisterRes, error)
	RegisterRootDir(request *types.RootDirRegisterReq) (*types.RootDirRegisterRes, error)
	GetRootDirList() (*types.AskRootDirRes, error)
	GetRootDirByPath(path string) (*types.RootDirectory, error)
}

type NetworkAdapter interface {
	UpdateClientConnection(uuid string, conn *qp.Connection) error
}
