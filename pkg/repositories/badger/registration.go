package badger

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
	"log"
)

const (
	PrefixClient  string = "client_"
	PrefixRootDir string = "root_dir_"
)

type RegistrationRepository struct {
}

func NewRegistrationRepository() *RegistrationRepository {
	return &RegistrationRepository{}
}

// SaveClient saves new client to badger and this system
func (repository *RegistrationRepository) SaveClient(uuid string, client types.Client) {
	key := []byte(PrefixClient + uuid)

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, client.Encode())
		return err
	})
	if err != nil {
		log.Println("quics: (SaveClient) ", err)
	}
}

// GetClientByUUID gets client by client uuid
func (repository *RegistrationRepository) GetClientByUUID(uuid string) *types.Client {
	key := []byte(PrefixClient + uuid)
	var client *types.Client

	err := db.View(func(txn *badger.Txn) error {
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

func (repository *RegistrationRepository) SaveRootDir(path string, rootDir types.RootDirectory) {
	key := []byte(PrefixRootDir + path)

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, rootDir.Encode())
		return err
	})
	if err != nil {
		log.Println("quics: (SaveClient) ", err)
	}
}

func (repository *RegistrationRepository) GetRootDirByPath(path string) *types.RootDirectory {
	key := []byte(PrefixRootDir + path)

	var rootDir *types.RootDirectory

	err := db.View(func(txn *badger.Txn) error {
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

func (repository *RegistrationRepository) GetAllRootDir() []types.RootDirectory {
	var rootDirs []types.RootDirectory

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte(PrefixRootDir)); it.ValidForPrefix([]byte(PrefixRootDir)); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			rootDir := &types.RootDirectory{}
			if err := rootDir.Decode(val); err != nil {
				return err
			}

			rootDirs = append(rootDirs, *rootDir)
		}

		return nil
	})
	if err != nil {
		return nil
	}

	return rootDirs
}

// GetSequence returns badger sequence by key
func (repository *RegistrationRepository) GetSequence(key []byte, increment uint64) (uint64, error) {
	seq, err := db.GetSequence(key, increment)
	if err != nil {
		log.Panicln("quics: (GetSequence) ", err)
	}
	defer seq.Release()

	return seq.Next()
}
