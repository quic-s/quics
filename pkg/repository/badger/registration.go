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
	var client *types.Client

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

func (rr *RegistrationRepository) SaveRootDir(afterPath string, rootDir *types.RootDirectory) error {
	key := []byte(PrefixRootDir + afterPath)

	err := rr.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, rootDir.Encode())
		return err
	})
	if err != nil {
		log.Println("quics: (SaveClient) ", err)
		return err
	}
	return nil
}

func (rr *RegistrationRepository) GetRootDirByPath(afterPath string) (*types.RootDirectory, error) {
	key := []byte(PrefixRootDir + afterPath)

	var rootDir *types.RootDirectory

	err := rr.db.View(func(txn *badger.Txn) error {
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
		return nil, err
	}

	return rootDir, nil
}

func (rr *RegistrationRepository) GetAllRootDir() ([]*types.RootDirectory, error) {
	rootDirs := []*types.RootDirectory{}
	err := rr.db.View(func(txn *badger.Txn) error {
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

			rootDirs = append(rootDirs, rootDir)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return rootDirs, nil
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
