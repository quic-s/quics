package server

import "github.com/quic-s/quics/pkg/types"

type Repository interface {
	GetAllClients() ([]*types.Client, error)
	GetAllRootDirectories() ([]*types.RootDirectory, error)
	GetAllFiles() ([]*types.File, error)
	GetClientByUUID(uuid string) (*types.Client, error)
	GetRootDirectoryByPath(afterPath string) (*types.RootDirectory, error)
	GetFileByAfterPath(afterPath string) (*types.File, error)
	DeleteAllClients() error
	DeleteAllRootDirectories() error
	DeleteAllFiles() error
	DeleteClientByUUID(uuid string) error
	DeleteRootDirectoryByAfterPath(afterPath string) error
	DeleteFileByAfterPath(afterPath string) error
}

type Service interface {
	StopServer() error
	ListenProtocol() error
	Ping(request *types.Ping) (*types.Ping, error)
	ShowClientLogs(all string, id string) error
	ShowDirLogs(all string, id string) error
	ShowFileLogs(all string, id string) error
	DisconnectClient(all string, id string) error
	DisconnectDir(all string, id string) error
	DisconnectFile(all string, id string) error
}
