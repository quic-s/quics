package registration

import (
	"errors"
	"log"
	"strings"

	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/types"
)

type Service struct {
	registrationRepository *Repository
}

func NewRegistrationService(registrationRepository *Repository) *Service {
	return &Service{registrationRepository: registrationRepository}
}

// CreateNewClient creates new client entity
func (registrationService *Service) CreateNewClient(request types.ClientRegisterReq, password string, ip string) error {

	if request.ClientPassword != password {
		return errors.New("password is not correct")
	}

	// create new id using badger sequence
	seq, err := registrationService.registrationRepository.DB.GetSequence([]byte("client"), 1)
	if err != nil {
		log.Println("Error while creating new id: ", err)
		return err
	}
	defer seq.Release()
	newId, err := seq.Next()

	// initialize client information
	//var newUUID = utils.CreateUUID()
	var client = types.Client{
		Id:   newId,
		Ip:   ip,
		UUID: request.UUID,
	}

	// Save client to badger database
	registrationService.registrationRepository.SaveClient(request.UUID, client)

	return nil
}

// RegisterRootDir registers initial root directory to client database
func (registrationService *Service) RegisterRootDir(request types.RootDirRegisterReq) error {
	// get client entity by uuid in request data
	client := registrationService.registrationRepository.GetClientByUUID(request.UUID)

	// create root directory entity
	path := config.GetSyncDirPath() + request.AfterPath
	var rootDir = types.RootDirectory{
		Path:     path,
		Owner:    client.UUID,
		Password: request.RootDirPassword,
	}
	rootDirs := append(client.Root, rootDir)
	client.Root = rootDirs

	// save updated client entity
	registrationService.registrationRepository.SaveClient(client.UUID, *client)

	// save requested root directory
	registrationService.registrationRepository.SaveRootDir(path, rootDir)

	return nil
}

func (registrationService *Service) SyncRootDir(request types.SyncRootDirReq) error {
	client := registrationService.registrationRepository.GetClientByUUID(request.UUID)

	path := config.GetSyncDirPath() + request.AfterPath
	rootDir := registrationService.registrationRepository.GetRootDirByPath(path)

	// password check
	if rootDir.Password != request.RootDirPassword {
		return errors.New("password is not correct")
	}

	rootDirs := append(client.Root, *rootDir)
	client.Root = rootDirs

	// save updated client entity with new root directory
	registrationService.registrationRepository.SaveClient(client.UUID, *client)

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
