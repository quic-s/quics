package registration

import (
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
func (registrationRepository *Repository) SaveClient(uuid string, client Client) {
	err := registrationRepository.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(uuid), client.Encode())
		return err
	})
	if err != nil {
		log.Panicf("Error while creating client: %s", err)
	}
}

// GetClientByUuid gets client by client uuid
func (registrationRepository *Repository) GetClientByUuid(uuid string) *Client {
	var client *Client

	err := registrationRepository.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(uuid))
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		client = &Client{}
		if err := client.Decode(val); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil
	}

	return client
}

func (registrationRepository *Repository) SaveRootDir(path string, rootDir RootDirectory) {
	err := registrationRepository.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(path), rootDir.Encode())
		return err
	})
	if err != nil {
		log.Panicf("Error while creating client: %s", err)
	}
}
