package registration

import (
	"github.com/quic-s/quics/pkg/types"
)

type Repository interface {
	SaveClient(uuid string, client types.Client)
	GetClientByUUID(uuid string) *types.Client
	SaveRootDir(path string, rootDir types.RootDirectory)
	GetRootDirByPath(path string) *types.RootDirectory
	GetAllRootDir() []types.RootDirectory
	GetSequence(key []byte, increment uint64) (uint64, error)
}

type Service interface {
	CreateNewClient(request types.ClientRegisterReq, password string, ip string) error
	RegisterRootDir(request types.RegisterRootDirReq) error
	SyncRootDir(request types.SyncRootDirReq) error
	GetRootDirList() []types.RootDirectory
	GetRootDirByPath(path string) types.RootDirectory
}
