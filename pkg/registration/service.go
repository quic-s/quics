package registration

import (
	"github.com/quic-s/quics/pkg/types"
	"log"
)

type Service struct {
	registrationRepository *Repository
}

func NewRegistrationService(registrationRepository *Repository) *Service {
	return &Service{registrationRepository: registrationRepository}
}

// CreateNewClient creates new client entity
func (registrationService *Service) CreateNewClient(uuid string, password string, ip string) error {

	// TODO: check if the password is correct

	// create new id using badger sequence
	seq, err := registrationService.registrationRepository.DB.GetSequence([]byte("client"), 1)
	if err != nil {
		log.Println("Error while creating new id: ", err)
		return err
	}
	defer seq.Release()
	newId, err := seq.Next()

	// initialize client information
	//var newUuid = utils.CreateUuid()
	var client = types.Client{
		Id:   newId,
		Ip:   ip,
		Uuid: uuid,
	}

	// Save client to badger database
	registrationService.registrationRepository.SaveClient(uuid, client)

	return nil
}

// RegisterRootDir registers initial root directory to client database
func (registrationService *Service) RegisterRootDir(request types.RegisterRootDirRequest) error {
	// get client entity by uuid in request data
	client := registrationService.registrationRepository.GetClientByUuid(request.Uuid)

	// create root directory entity
	path := request.BeforePath + request.AfterPath
	var rootDir = types.RootDirectory{
		Path:     path,
		Owner:    client.Uuid,
		Password: request.RootDirPassword,
	}
	rootDirs := append(client.Root, rootDir)
	client.Root = rootDirs

	// save updated client entity
	registrationService.registrationRepository.SaveClient(client.Uuid, *client)

	// save requested root directory
	registrationService.registrationRepository.SaveRootDir(path, rootDir)

	return nil
}

//func (registrationService *Service) SyncRootDir(request SyncRootDirRequest) error {
//	// get client entity by uuid in request data
//	client, err := registrationService.registrationRepository.GetClientByUuid(request.Uuid)
//	if err != nil {
//		log.Printf("Error while registering root directory: %s\n", err)
//		return err
//	}
//
//	// save updated client entity
//	registrationService.registrationRepository.SaveClient(client.Uuid, *client)
//
//	return nil
//}
