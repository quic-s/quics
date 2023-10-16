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
	err = rs.registrationRepository.SaveClient(request.UUID, client)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	err = rs.networkAdapter.UpdateClientConnection(request.UUID, conn)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	return &types.ClientRegisterRes{
		UUID: request.UUID,
	}, nil
}
