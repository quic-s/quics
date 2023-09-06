package registration

import (
	"errors"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/types"
	"log"
	"strings"
)

type Service struct {
	registrationRepository *Repository
}

func NewRegistrationService(registrationRepository *Repository) *Service {
	return &Service{registrationRepository: registrationRepository}
}

// CreateNewClient creates new client entity
func (registrationService *Service) CreateNewClient(request types.RegisterClientRequest, password string, ip string) error {

	log.Println("request password: ", request.ClientPassword)
	log.Println("server password: ", password)

	// TODO: check if the password is correct
	//if request.ClientPassword != password {
	//	return errors.New("password is not correct")
	//}

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
		Uuid: request.Uuid,
	}

	// Save client to badger database
	registrationService.registrationRepository.SaveClient(request.Uuid, client)

	return nil
}

// RegisterRootDir registers initial root directory to client database
func (registrationService *Service) RegisterRootDir(request types.RegisterRootDirRequest) error {
	// get client entity by uuid in request data
	client := registrationService.registrationRepository.GetClientByUuid(request.Uuid)

	// create root directory entity
	path := config.GetSyncDirPath() + request.AfterPath
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

func (registrationService *Service) SyncRootDir(request types.SyncRootDirRequest) error {
	client := registrationService.registrationRepository.GetClientByUuid(request.Uuid)

	path := config.GetSyncDirPath() + request.AfterPath
	rootDir := registrationService.registrationRepository.GetRootDirByPath(path)

	// password check
	if rootDir.Password != request.RootDirPassword {
		return errors.New("password is not correct")
	}

	rootDirs := append(client.Root, *rootDir)
	client.Root = rootDirs

	// save updated client entity with new root directory
	registrationService.registrationRepository.SaveClient(client.Uuid, *client)

	return nil
}

func (registrationService *Service) GetRootDirList() []*types.RootDirectory {
	rootDirs := registrationService.registrationRepository.GetAllRootDir()
	return rootDirs
}

func (registrationService *Service) GetRootDirByPath(path string) types.RootDirectory {
	rootDir := registrationService.registrationRepository.GetRootDirByPath(path)
	return *rootDir
}

func ExtractRelateiveRootDirPath(fullPath string) string {
	dirPath := config.GetDirPath()

	// if fullPath starts with dirPath, extract the rest of the string
	if strings.HasPrefix(fullPath, dirPath) {
		relativePath := strings.TrimPrefix(fullPath, dirPath)
		return relativePath
	}

	// if fullPath does not start with dirPath, return fullPath as it is.
	return fullPath
}
