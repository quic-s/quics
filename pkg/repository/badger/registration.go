package badger

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
)

const (
	PrefixClient  string = "client_"
	PrefixRootDir string = "root_dir_"
)

type RegistrationRepository struct {
	db *badger.DB
}

// SaveClient saves new client to badger and this system
func (rr *RegistrationRepository) SaveClient(uuid string, client *types.Client) error {
	key := []byte(PrefixClient + uuid)

	err := rr.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, client.Encode())
		return err
	})
	if err != nil {
		log.Println("quics: (SaveClient) ", err)
		return err
	}
	return nil
}

// GetClientByUUID gets client by client uuid
func (rr *RegistrationRepository) GetClientByUUID(uuid string) (*types.Client, error) {
	key := []byte(PrefixClient + uuid)

	client := &types.Client{}
	err := rr.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
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
		return nil, err
	}

	return client, nil
}

// GetAllClients gets all clients
func (rr *RegistrationRepository) GetAllClients() ([]types.Client, error) {
	clients := []types.Client{}

	err := rr.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(PrefixClient)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			client := types.Client{}
			if err := client.Decode(val); err != nil {
				return err
			}

			clients = append(clients, client)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return clients, nil
}

// GetSequence returns badger sequence by key
func (rr *RegistrationRepository) GetSequence(key []byte, increment uint64) (uint64, error) {
	seq, err := rr.db.GetSequence(key, increment)
	if err != nil {
		log.Panicln("quics: (GetSequence) ", err)
	}
	defer seq.Release()

	return seq.Next()
}
