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
func (clientRepository *Repository) SaveClient(uuid string, client Client) {
	err := clientRepository.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(uuid), client.Encode())
		return err
	})
	if err != nil {
		log.Panicf("Error while creating client: %s", err)
	}
}

// GetClientByUuid gets client by client uuid
//func (clientRepository *Repository) GetClientByUuid(uuid string) *Client {
//	err := clientRepository.DB.View(func(txn *badger.Txn) error {
//		item, err := txn.Get([]byte(uuid))
//		if err != nil {
//			return err
//		}
//
//		err = item.Value(func(val []byte) error {
//			return Client.Decode(val)
//		})
//		if err != nil {
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		return nil
//	}
//
//	return Client.Decode(item)
//}
