package registration

import (
	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/types"
)

type Repository interface {
	SaveClient(uuid string, client *types.Client) error
	GetClientByUUID(uuid string) (*types.Client, error)
	GetAllClients() ([]types.Client, error)
	DeleteClient(uuid string) error
	GetSequence(key []byte, increment uint64) (uint64, error)
	ErrKeyNotFound() error
}

type Service interface {
	RegisterClient(request *types.ClientRegisterReq, conn *qp.Connection) (*types.ClientRegisterRes, error)
}

type NetworkAdapter interface {
	UpdateClientConnection(uuid string, conn *qp.Connection) error
	DeleteConnection(uuid string) error
}
