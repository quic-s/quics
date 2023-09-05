package registration

import (
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
	var client = Client{
		Id:   newId,
		Ip:   ip,
		Uuid: uuid,
	}

	// Save client to badger database
	registrationService.registrationRepository.SaveClient(uuid, client)

	return nil
}

// RegisterRootDir registers initial root directory to client database
//func (registrationService *Service) RegisterRootDir(request RegisterRootDirRequest) error {
//	// get client entity by uuid in request data
//	client, err := registrationService.registrationRepository.GetClientByUuid(request.Uuid)
//	if err != nil {
//		log.Printf("Error while registering root directory: %s\n", err)
//		return err
//	}
//
//	// create root directory entity
//	// TODO: need to check the time zone
//	var rootDir = RootDirectory{
//		Owner: client.Uuid,
//		Path:  request.BeforePath + request.AfterPath,
//	}
//	rootDirs := append(client.Root, rootDir)
//	client.Root = rootDirs
//
//	// save updated client entity
//	registrationService.registrationRepository.SaveClient(client.Uuid, *client)
//
//	return nil
//}

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
