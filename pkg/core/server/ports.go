package server

import (
	"io"

	"github.com/quic-s/quics/pkg/types"
)

type Repository interface {
	UpdatePassword(server *types.Server) error
	DeletePassword() error
	GetPassword() (*types.Server, error)
	GetAllClients() ([]types.Client, error)
	GetAllRootDirectories() ([]types.RootDirectory, error)
	GetAllFiles() ([]types.File, error)
	GetClientByUUID(uuid string) (*types.Client, error)
	GetRootDirectoryByPath(afterPath string) (*types.RootDirectory, error)
	GetFileByAfterPath(afterPath string) (*types.File, error)
	DeleteAllClients() error
	DeleteAllRootDirectories() error
	DeleteAllFiles() error
	DeleteClientByUUID(uuid string) error
	DeleteRootDirectoryByAfterPath(afterPath string) error
	DeleteFileByAfterPath(afterPath string) error
	GetAllHistories() ([]types.FileHistory, error)
	GetHistoryByAfterPath(afterPath string) (*types.FileHistory, error)
}

type Service interface {
	StopServer() error
	ListenProtocol() error
	SetPassword(request *types.Server) error
	ResetPassword() error
	Ping(request *types.Ping) (*types.Ping, error)
	ShowClient(uuid string) ([]types.Client, error)
	ShowDir(afterPath string) ([]types.RootDirectory, error)
	ShowFile(afterPath string) ([]types.File, error)
	ShowHistory(afterPath string) ([]types.FileHistory, error)
	RemoveClient(uuid string) error
	RemoveDir(afterPath string) error
	RemoveFile(afterPath string) error
	DownloadFile(afterPath string, timestamp uint64) (*types.FileMetadata, io.Reader, error)
}

type SyncDirAdapter interface {
	GetFileFromHistoryDir(afterPath string, timestamp uint64) (*types.FileMetadata, io.Reader, error)
}
