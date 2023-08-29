package sharing

type Service struct {
	sharingRepository *Repository
}

func NewSharingService(sharingRepository *Repository) *Service {
	return &Service{sharingRepository: sharingRepository}
}
