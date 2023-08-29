package history

type Service struct {
	historyRepository *Repository
}

func NewHistoryService(historyRepository *Repository) *Service {
	return &Service{historyRepository: historyRepository}
}
