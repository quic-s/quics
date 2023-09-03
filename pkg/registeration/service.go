package registeration

import (
	"encoding/binary"
	"github.com/quic-s/quics/pkg/utils"
	"log"
)

type Service struct {
	clientRepository *Repository
}

func NewRegistrationService(clientRepository *Repository) *Service {
	return &Service{clientRepository: clientRepository}
}

// CreateNewClient creates new client entity
func (registrationService *Service) CreateNewClient(ip *string) (string, error) {
	// create new id using badger sequence
	seq, err := registrationService.clientRepository.DB.GetSequence([]byte("client"), 1)
	if err != nil {
		log.Panicf("Error while creating new id: %s", err)
	}
	defer seq.Release()
	newId, err := seq.Next()

	// FIXME: if not necessary, remove 2 line code below
	newIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(newIdBytes, newId)

	// initialize client information
	var newUuid = utils.CreateUuid()
	var client = Client{
		Id:   newId,
		Ip:   *ip,
		Uuid: newUuid,
	}

	// Save client to badger database
	registrationService.clientRepository.SaveClient(newUuid, client)

	return newUuid, nil
}

// RegisterRootDir registers initial root directory to client database
func (registrationService *Service) RegisterRootDir(request RegisterRootDirRequest) (string, error) {
	// get client entity by uuid in request data
	client, err := registrationService.clientRepository.GetClientByUuid(request.Uuid)
	if err != nil {
		log.Fatalf("Error while registering root directory: %s", err)
	}

	// create root directory entity
	// TODO: need to check the time zone
	var rootDir = RootDirectory{
		Owner: client.Uuid,
		Path:  request.BeforePath + request.AfterPath,
	}
	rootDirs := append(client.Root, rootDir)
	client.Root = rootDirs

	// save updated client entity
	registrationService.clientRepository.SaveClient(client.Uuid, *client)

	return "Success to registeration root directroy", nil
}
