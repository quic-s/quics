package registration

import (
	"github.com/quic-s/quics/pkg/types"
	"log"

	"github.com/dgraph-io/badger/v3"
)

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
	SaveClient(newId []byte, client types.Client)
	GetClientById(id string) (*types.Client, error)
}

func NewClientRepository(db *badger.DB) *Repository {
	return &Repository{DB: db}
}

// SaveClient saves new client to badger and this system
func (registrationRepository *Repository) SaveClient(uuid string, client types.Client) {
	err := registrationRepository.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(uuid), client.Encode())
		return err
	})
	if err != nil {
		log.Panicf("Error while creating client: %s", err)
	}
}

// GetClientByUuid gets client by client uuid
func (registrationRepository *Repository) GetClientByUuid(uuid string) *types.Client {
	var client *types.Client

	err := registrationRepository.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(uuid))
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		client = &types.Client{}
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

func (registrationRepository *Repository) SaveRootDir(path string, rootDir types.RootDirectory) {
	err := registrationRepository.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(path), rootDir.Encode())
		return err
	})
	if err != nil {
		log.Panicf("Error while creating client: %s", err)
	}
}
