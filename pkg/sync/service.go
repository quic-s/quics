package sync

type Service struct {
	syncRepository *Repository
}

func NewSyncService(syncRepository *Repository) *Service {
	return &Service{syncRepository: syncRepository}
}
