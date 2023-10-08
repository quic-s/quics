package history

import "github.com/quic-s/quics/pkg/types"

type Repository interface {
	SaveNewFileHistory(afterPath string, fileHistory *types.FileHistory) error
	GetFileHistory(afterPath string) (*types.FileHistory, error)
}

type Service interface {
}
