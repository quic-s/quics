package registration

import (
	"log"

	"github.com/quic-s/quics/pkg/types"

	"github.com/dgraph-io/badger/v3"
)

const (
	PrefixClient  string = "client_"
	PrefixRootdir string = "root_dir_"
)

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
	SaveClient(newId []byte, client types.Client)
	GetClientById(id string) (*types.Client, error)
}

func NewRegistrationRepository(db *badger.DB) *Repository {
	return &Repository{DB: db}
}

// SaveClient saves new client to badger and this system
func (registrationRepository *Repository) SaveClient(uuid string, client types.Client) {
	key := []byte(PrefixClient + uuid)

	err := registrationRepository.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, client.Encode())
		return err
	})
	if err != nil {
		log.Panicf("Error while creating client: %s", err)
	}
}

// GetClientByUUID gets client by client uuid
func (registrationRepository *Repository) GetClientByUUID(uuid string) *types.Client {
	key := []byte(PrefixClient + uuid)
	var client *types.Client

	err := registrationRepository.DB.View(func(txn *badger.Txn) error {
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
		return nil
	}

	return client
}

func (registrationRepository *Repository) SaveRootDir(path string, rootDir types.RootDirectory) {
	key := []byte(PrefixRootdir + path)

	err := registrationRepository.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, rootDir.Encode())
		return err
	})
	if err != nil {
		log.Panicf("Error while creating client: %s", err)
	}
}

func (registrationRepository *Repository) GetRootDirByPath(path string) *types.RootDirectory {
	key := []byte(PrefixRootdir + path)

	var rootDir *types.RootDirectory

	err := registrationRepository.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		rootDir = &types.RootDirectory{}
		if err := rootDir.Decode(val); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil
	}

	return rootDir
}

func (registrationRepository *Repository) GetAllRootDir() []*types.RootDirectory {
	var rootDirs []*types.RootDirectory

	err := registrationRepository.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte(PrefixRootdir)); it.ValidForPrefix([]byte(PrefixRootdir)); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			rootDir := &types.RootDirectory{}
			if err := rootDir.Decode(val); err != nil {
				return err
			}

			rootDirs = append(rootDirs, rootDir)
		}

		return nil
	})
	if err != nil {
		return nil
	}

	return rootDirs
}
