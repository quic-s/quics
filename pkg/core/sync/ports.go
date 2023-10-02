package sync

import "github.com/quic-s/quics/pkg/types"

type Repository interface {
	GetFileByPath(path string) *types.File
	SaveFileByPath(path string, file types.File) error
	GetAllFiles() []types.File
}

type Service interface {
	IsExistFile(path string) int
	SaveFileFromPleaseSync(path string, file types.File) error
	GetFilesByRootDir(rootDirPath string) []types.File
	GetFiles() []types.File
	GetFileByPath(path string) *types.File
}
