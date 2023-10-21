package sharing

type SharingService struct {
	sharingRepository Repository
}

const (
	PrefixLink = "http://"
)

func NewService(sharingRepository Repository) *SharingService {
	return &SharingService{
		sharingRepository: sharingRepository,
	}
}

func (ss *SharingService) CreateLink(UUID string, afterPath string, count uint64) (string, error) {
	return "", nil
}
