package history

type MyHistoryService struct {
	historyRepository Repository
}

func NewHistoryService(historyRepository Repository) *MyHistoryService {
	return &MyHistoryService{
		historyRepository: historyRepository,
	}
}
