package metadata

type MyMetadataService struct {
	metadataRepository Repository
}

func NewMetadataService(metadataRepository Repository) *MyMetadataService {
	return &MyMetadataService{
		metadataRepository: metadataRepository,
	}
}
