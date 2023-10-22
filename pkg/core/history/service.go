package history

import (
	"errors"

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

	// request.CntFromHead만큼 뒤에서 세서 개수를 보낸다.
	length := len(histories)
	if request.CntFromHead > uint64(length) {
		return nil, errors.New("index out of bound")
	}

	histories = histories[length-1-int(request.CntFromHead):]

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
