package sync

import "github.com/quic-s/quics/pkg/client"

type Service struct {
	syncRepository *Repository
}

func NewRegistrationService(clientRepository *client.Repository) *Service {
	return &Service{clientRepository: clientRepository}
}
