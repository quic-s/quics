package registration

import (
	"errors"
	"log"

	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/types"
)

type RegistrationService struct {
	password               string
	registrationRepository Repository
	networkAdapter         NetworkAdapter
}

// NewRegistrationService creates new registration service
func NewService(password string, registrationRepository Repository, networkAdapter NetworkAdapter) Service {
	return &RegistrationService{
		password:               password,
		registrationRepository: registrationRepository,
		networkAdapter:         networkAdapter,
	}
}

// CreateNewClient creates new client entity
func (rs *RegistrationService) RegisterClient(request *types.ClientRegisterReq, conn *qp.Connection) (*types.ClientRegisterRes, error) {
	log.Println("quics: RegisterClient: ", request)
	if request.ClientPassword != rs.password {
		return nil, errors.New("[RegistrationService.RegitserClient] password is not correct")
	}
	client, err := rs.registrationRepository.GetClientByUUID(request.UUID)
	if err != nil && err != rs.registrationRepository.ErrKeyNotFound() {
		err = errors.New("[RegistrationService.RegitserClient] get client by uuid: " + err.Error())
		return nil, err
	}

	// if client is already existed, just update connection
	if client != nil && request.UUID == client.UUID {
		err = rs.networkAdapter.UpdateClientConnection(request.UUID, conn)
		if err != nil {
			err = errors.New("[RegistrationService.RegitserClient] update client connection: " + err.Error())
			return nil, err
		}
		return &types.ClientRegisterRes{
			UUID: request.UUID,
		}, nil
	}

	// create new id using badger sequence
	newId, err := rs.registrationRepository.GetSequence([]byte("client"), 1)
	if err != nil {
		err = errors.New("[RegistrationService.RegitserClient] get sequence: " + err.Error())
		return nil, err
	}

	// initialize client information
	client = &types.Client{
		Id:   newId,
		UUID: request.UUID,
	}

	// Save client to badger database
	err = rs.registrationRepository.SaveClient(request.UUID, client)
	if err != nil {
		err = errors.New("[RegistrationService.RegitserClient] save client to repository: " + err.Error())
		return nil, err
	}

	err = rs.networkAdapter.UpdateClientConnection(request.UUID, conn)
	if err != nil {
		err = errors.New("[RegistrationService.RegitserClient] update client connection: " + err.Error())
		return nil, err
	}

	return &types.ClientRegisterRes{
		UUID: request.UUID,
	}, nil
}

// CreateNewClient creates new client entity
func (rs *RegistrationService) DisconnectClient(request *types.DisconnectClientReq, conn *qp.Connection) (*types.DisconnectClientRes, error) {
	log.Println("quics: DisconnectClient: ", request)
	// Save client to badger database
	err := rs.registrationRepository.DeleteClient(request.UUID)
	if err != nil {
		err = errors.New("[RegistrationService.DisconnectClient] delete client from repository: " + err.Error())
		return nil, err
	}

	err = rs.networkAdapter.DeleteConnection(request.UUID)
	if err != nil {
		err = errors.New("[RegistrationService.DisconnectClient] delete client connection: " + err.Error())
		return nil, err
	}

	return &types.DisconnectClientRes{
		UUID: request.UUID,
	}, nil
}
