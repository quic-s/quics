package client

import (
	"encoding/json"
	"log"

	"github.com/dgraph-io/badger/v3"
)

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
	SaveClient(newId []byte, client Client)
	GetClientById(id string) (*Client, error)
}

func NewClientRepository(db *badger.DB) *Repository {
	return &Repository{DB: db}
}

// SaveClient saves new client to badger and this system
func (clientRepository *Repository) SaveClient(newId []byte, client Client) {
	clientJson, err := json.Marshal(client)
	if err != nil {
		log.Panicf("Error while marshaling request data: %s", err)
	}

	err = clientRepository.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(newId, clientJson)
		return err
	})
	if err != nil {
		log.Panicf("Error while creating client: %s", err)
	}
}

// GetClientById gets client by client uuid
func (clientRepository *Repository) GetClientById(id string) (*Client, error) {
	var client *Client
	err := clientRepository.DB.View(func(txn *badger.Txn) error {
		return nil
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}