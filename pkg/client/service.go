package client

import (
	"encoding/binary"
	"log"

	"github.com/quic-s/quics/pkg/utils/common"
)

type Service struct {
	clientRepository *Repository
}

func NewClientService(clientRepository *Repository) *Service {
	return &Service{clientRepository: clientRepository}
}

// CreateNewClient creates new client entity
func (clientService *Service) CreateNewClient(ip *string) (string, error) {
	// create new id using badger sequence
	seq, err := clientService.clientRepository.DB.GetSequence([]byte("client"), 1)
	if err != nil {
		log.Panicf("Error while creating new id: %s", err)
	}
	defer seq.Release()
	newId, err := seq.Next()

	// FIXME: if not necessary, remove 2 line code below
	newIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(newIdBytes, newId)

	// initialize client information
	var newUuid = common.CreateUuid()
	var client = Client{
		Id:   newId,
		Ip:   *ip,
		Uuid: newUuid,
	}

	// Save client to badger database
	clientService.clientRepository.SaveClient(newUuid, client)

	return newUuid, nil
}
