package history

import "github.com/quic-s/quics/pkg/types"

type Repository interface {
	SaveNewFileHistory(afterPath string, fileHistory *types.FileHistory) error
	GetFileHistory(afterPath string, timestamp uint64) (*types.FileHistory, error)
	GetFileHistoriesForClient(afterPath string, cntFromHead uint64) ([]types.FileHistory, error)
}

type Service interface {
	ShowHistory(request *types.ShowHistoryReq) (*types.ShowHistoryRes, error)
	DownloadHistory(request *types.DownloadHistoryReq) (*types.DownloadHistoryRes, string, error)
}
