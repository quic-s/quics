package metadata

type Service struct {
	metadataRepository *Repository
}

func NewMetadataService(metadataRepository *Repository) *Service {
	return &Service{metadataRepository: metadataRepository}
}
