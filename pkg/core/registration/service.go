package registration

import (
	"errors"
	"log"

	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/types"
	"github.com/quic-s/quics/pkg/utils"
)

type RegistrationService struct {
	password               string
	registrationRepository Repository
	networkAdapter         NetworkAdapter
}

// NewRegistrationService creates new registration service
func NewService(password string, registrationRepository Repository, networkAdapter NetworkAdapter) *RegistrationService {
	return &RegistrationService{
		password:               password,
		registrationRepository: registrationRepository,
		networkAdapter:         networkAdapter,
	}
}

// CreateNewClient creates new client entity
func (rs *RegistrationService) RegisterClient(request *types.ClientRegisterReq, conn *qp.Connection) (*types.ClientRegisterRes, error) {
	if request.ClientPassword != rs.password {
		return nil, errors.New("quics: (CreateNewClient) password is not correct")
	}

	// create new id using badger sequence
	newId, err := rs.registrationRepository.GetSequence([]byte("client"), 1)
	if err != nil {
		log.Println("quics: (CreateNewClient) error while getting sequence")
		return nil, err
	}

	// initialize client information
	client := &types.Client{
		Id:   newId,
		UUID: request.UUID,
	}

	// Save client to badger database
	rs.registrationRepository.SaveClient(request.UUID, client)

	rs.networkAdapter.UpdateClientConnection(request.UUID, conn)

	return &types.ClientRegisterRes{
		UUID: request.UUID,
	}, nil
}

// RegisterRootDir registers initial root directory to client database
func (rs *RegistrationService) RegisterRootDir(request *types.RootDirRegisterReq) (*types.RootDirRegisterRes, error) {
	// get client entity by uuid in request data
	client, err := rs.registrationRepository.GetClientByUUID(request.UUID)
	if err != nil {
		return nil, err
	}

	// create root directory entity
	path := utils.GetQuicsSyncDirPath() + request.AfterPath
	rootDir := types.RootDirectory{
		BeforePath: utils.GetQuicsSyncDirPath(),
		AfterPath:  request.AfterPath,
		Owner:      client.UUID,
		Password:   request.RootDirPassword,
	}
	rootDirs := append(client.Root, rootDir)
	client.Root = rootDirs

	// save updated client entity
	err = rs.registrationRepository.SaveClient(client.UUID, client)
	if err != nil {
		return nil, err
	}

	// save requested root directory
	err = rs.registrationRepository.SaveRootDir(path, &rootDir)
	if err != nil {
		return nil, err
	}

	return &types.RootDirRegisterRes{
		UUID: request.UUID,
	}, nil
}

// GetRootDirList gets root directory list of client
func (rs *RegistrationService) GetRootDirList() ([]*types.RootDirectory, error) {
	rootDirs, err := rs.registrationRepository.GetAllRootDir()
	if err != nil {
		log.Println("quics: ", err)
	}

	return rootDirs, err
}

// GetRootDirByPath gets root directory by path
func (rs *RegistrationService) GetRootDirByPath(path string) (*types.RootDirectory, error) {
	rootDir, err := rs.registrationRepository.GetRootDirByPath(path)
	if err != nil {
		log.Println("quics: ", err)
	}

	return rootDir, err
}
