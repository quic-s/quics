package history

import (
	"github.com/quic-s/quics/pkg/types"
)

type HistoryService struct {
	historyRepository Repository
}

func NewService(historyRepository Repository) *HistoryService {
	return &HistoryService{
		historyRepository: historyRepository,
	}
}

func (hs *HistoryService) ShowHistory(request *types.ShowHistoryReq) (*types.ShowHistoryRes, error) {
	histories, err := hs.historyRepository.GetFileHistoriesForClient(request.AfterPath, request.CntFromHead)
	if err != nil {
		return nil, err
	}

	return &types.ShowHistoryRes{
		History: histories,
	}, nil
}

func (hs *HistoryService) DownloadHistory(request *types.DownloadHistoryReq) (*types.DownloadHistoryRes, string, error) {
	history, err := hs.historyRepository.GetFileHistory(request.AfterPath, request.Version)
	if err != nil {
		return nil, "", err
	}

	filePath := history.BeforePath + "/" + history.AfterPath

	return &types.DownloadHistoryRes{
		UUID: request.UUID,
	}, filePath, nil
}
