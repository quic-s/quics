package registration

import (
	"github.com/quic-s/quics/pkg/client"
	"github.com/quic-s/quics/pkg/sync"
	"log"
)

type Service struct {
	clientRepository *client.Repository
}

func NewRegistrationService(clientRepository *client.Repository) *Service {
	return &Service{clientRepository: clientRepository}
}

// RegisterRootDir registers initial root directory to client database
func (registrationService *Service) RegisterRootDir(request sync.RegisterRootDirRequest) (string, error) {
	// get client entity by uuid in request data
	client, err := registrationService.clientRepository.GetClientByUuid(request.Uuid)
	if err != nil {
		log.Fatalf("Error while registering root directory: %s", err)
	}

	// create root directory entity
	// TODO: need to check the time zone
	var rootDir = sync.RootDirectory{
		Path: request.Path,
		Date: request.Date,
	}
	client.Root = rootDir

	// save updated client entity
	registrationService.clientRepository.SaveClient(client.Uuid, *client)

	return "Success to register root directroy", nil
}
