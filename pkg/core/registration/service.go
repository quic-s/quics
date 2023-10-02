package registration

import (
	"errors"
	"github.com/quic-s/quics/pkg/types"
	"github.com/quic-s/quics/pkg/utils"
	"log"
)

type MyRegistrationService struct {
	registrationRepository Repository
}

// NewRegistrationService creates new registration service
func NewRegistrationService(registrationRepository Repository) *MyRegistrationService {
	return &MyRegistrationService{
		registrationRepository: registrationRepository,
	}
}

// CreateNewClient creates new client entity
func (service *MyRegistrationService) CreateNewClient(request types.ClientRegisterReq, password string, ip string) error {

	if request.ClientPassword != password {
		return errors.New("quics: (CreateNewClient) password is not correct")
	}

	// create new id using badger sequence
	newId, err := service.registrationRepository.GetSequence([]byte("client"), 1)
	if err != nil {
		log.Println("quics: (CreateNewClient) error while getting sequence")
		return err
	}

	// initialize client information
	var client = types.Client{
		Id:   newId,
		Ip:   ip,
		UUID: request.UUID,
	}

	// Save client to badger database
	service.registrationRepository.SaveClient(request.UUID, client)

	return nil
}

// RegisterRootDir registers initial root directory to client database
func (service *MyRegistrationService) RegisterRootDir(request types.RegisterRootDirReq) error {
	// get client entity by uuid in request data
	client := service.registrationRepository.GetClientByUUID(request.UUID)

	// create root directory entity
	path := utils.GetQuicsSyncDirPath() + request.AfterPath
	var rootDir = types.RootDirectory{
		BeforePath: utils.GetQuicsSyncDirPath(),
		AfterPath:  request.AfterPath,
		Owner:      client.UUID,
		Password:   request.RootDirPassword,
	}
	rootDirs := append(client.Root, rootDir)
	client.Root = rootDirs

	// save updated client entity
	service.registrationRepository.SaveClient(client.UUID, *client)

	// save requested root directory
	service.registrationRepository.SaveRootDir(path, rootDir)

	return nil
}

// SyncRootDir syncs root directory to other client from owner client
func (service *MyRegistrationService) SyncRootDir(request types.SyncRootDirReq) error {
	client := service.registrationRepository.GetClientByUUID(request.UUID)

	path := utils.GetQuicsSyncDirPath() + request.AfterPath
	rootDir := service.registrationRepository.GetRootDirByPath(path)

	// password check
	if rootDir.Password != request.RootDirPassword {
		return errors.New("quics: (SyncRootDir) password is not correct")
	}

	rootDirs := append(client.Root, *rootDir)
	client.Root = rootDirs

	// save updated client entity with new root directory
	service.registrationRepository.SaveClient(client.UUID, *client)

	return nil
}

// GetRootDirList gets root directory list of client
func (service *MyRegistrationService) GetRootDirList() []types.RootDirectory {
	rootDirs := service.registrationRepository.GetAllRootDir()
	return rootDirs
}

// GetRootDirByPath gets root directory by path
func (service *MyRegistrationService) GetRootDirByPath(path string) types.RootDirectory {
	rootDir := service.registrationRepository.GetRootDirByPath(path)
	return *rootDir
}
