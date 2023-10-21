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
	GetAllHistories() ([]*types.FileHistory, error)
	GetHistoryByAfterPath(afterPath string) (*types.FileHistory, error)
}

type Service interface {
	StopServer() error
	ListenProtocol() error
	Ping(request *types.Ping) (*types.Ping, error)
	ShowClientLogs(all string, id string) error
	ShowDirLogs(all string, id string) error
	ShowFileLogs(all string, id string) error
	ShowHistoryLogs(all string, id string) error
	RemoveClient(all string, id string) error
	RemoveDir(all string, id string) error
	RemoveFile(all string, id string) error
	RollbackFile(path string, version string) error
	DownloadFile(path string, version string, target string) error
}
