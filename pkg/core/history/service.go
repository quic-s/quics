package history

type HistoryService struct {
	historyRepository Repository
}

func NewHistoryService(historyRepository Repository) *HistoryService {
	return &HistoryService{
		historyRepository: historyRepository,
	}
}
