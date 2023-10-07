package sharing

type SharingService struct {
	sharingService Repository
}

func NewService(sharingService Repository) *SharingService {
	return &SharingService{
		sharingService: sharingService,
	}
}
