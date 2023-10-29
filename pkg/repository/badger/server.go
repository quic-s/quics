package badger

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
)

const (
	PrefixServerPassword = "password_"
)

type ServerRepository struct {
	db *badger.DB
}

func (sr *ServerRepository) UpdatePassword(server *types.Server) error {
	key := []byte(PrefixServerPassword)

	err := sr.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set(key, server.Encode()); err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (sr *ServerRepository) DeletePassword() error {
	key := []byte(PrefixServerPassword)

	err := sr.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete(key); err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (sr *ServerRepository) GetPassword() (*types.Server, error) {
	key := []byte(PrefixServerPassword)
	server := &types.Server{}

	err := sr.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		if err := server.Decode(val); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (sr *ServerRepository) GetAllClients() ([]*types.Client, error) {
	clients := []*types.Client{}

	err := sr.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte(PrefixClient)); it.ValidForPrefix([]byte(PrefixClient)); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			client := &types.Client{}
			if err := client.Decode(val); err != nil {
				log.Println("quics err: ", err)
				return err
			}

			clients = append(clients, client)
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	return clients, nil
}

func (sr *ServerRepository) GetAllRootDirectories() ([]*types.RootDirectory, error) {
	rootDirs := []*types.RootDirectory{}

	err := sr.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte(PrefixRootDir)); it.ValidForPrefix([]byte(PrefixRootDir)); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			rootDir := &types.RootDirectory{}
			if err := rootDir.Decode(val); err != nil {
				log.Println("quics err: ", err)
				return err
			}

			rootDirs = append(rootDirs, rootDir)
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	return rootDirs, nil
}

func (sr *ServerRepository) GetAllFiles() ([]*types.File, error) {
	files := []*types.File{}

	err := sr.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte(PrefixFile)); it.ValidForPrefix([]byte(PrefixFile)); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			file := &types.File{}
			if err := file.Decode(val); err != nil {
				log.Println("quics err: ", err)
				return err
			}

			files = append(files, file)
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	return files, nil
}

func (sr *ServerRepository) GetClientByUUID(uuid string) (*types.Client, error) {
	key := []byte(PrefixClient + uuid)
	client := &types.Client{}

	err := sr.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			log.Println("quics err: ", err)
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			log.Println("quics err: ", err)
			return err
		}

		if err := client.Decode(val); err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	return client, nil
}

func (sr *ServerRepository) GetRootDirectoryByPath(afterPath string) (*types.RootDirectory, error) {
	key := []byte(PrefixRootDir + afterPath)
	rootDir := &types.RootDirectory{}

	err := sr.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			log.Println("quics err: ", err)
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			log.Println("quics err: ", err)
			return err
		}

		if err := rootDir.Decode(val); err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	return rootDir, nil
}

func (sr *ServerRepository) GetFileByAfterPath(afterPath string) (*types.File, error) {
	key := []byte(PrefixFile + afterPath)
	file := &types.File{}

	err := sr.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			log.Println("quics err: ", err)
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			log.Println("quics err: ", err)
			return err
		}

		if err := file.Decode(val); err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	return file, nil
}

func (sr *ServerRepository) DeleteAllClients() error {
	key := []byte(PrefixClient)

	err := sr.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(key); it.ValidForPrefix(key); it.Next() {
			item := it.Item()
			if err := txn.Delete(item.Key()); err != nil {
				log.Println("quics err: ", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (sr *ServerRepository) DeleteAllRootDirectories() error {
	key := []byte(PrefixRootDir)

	err := sr.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(key); it.ValidForPrefix(key); it.Next() {
			item := it.Item()
			if err := txn.Delete(item.Key()); err != nil {
				log.Println("quics err: ", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (sr *ServerRepository) DeleteAllFiles() error {
	key := []byte(PrefixFile)

	err := sr.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(key); it.ValidForPrefix(key); it.Next() {
			item := it.Item()
			if err := txn.Delete(item.Key()); err != nil {
				log.Println("quics err: ", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (sr *ServerRepository) DeleteClientByUUID(uuid string) error {
	key := []byte(PrefixClient + uuid)

	err := sr.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete(key); err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (sr *ServerRepository) DeleteRootDirectoryByAfterPath(afterPath string) error {
	key := []byte(PrefixRootDir + afterPath)

	err := sr.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete(key); err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (sr *ServerRepository) DeleteFileByAfterPath(afterPath string) error {
	key := []byte(PrefixFile + afterPath)

	err := sr.db.Update(func(txn *badger.Txn) error {
		if err := txn.Delete(key); err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (sr *ServerRepository) GetAllHistories() ([]*types.FileHistory, error) {
	histories := []*types.FileHistory{}

	err := sr.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte(PrefixHistory)); it.ValidForPrefix([]byte(PrefixHistory)); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				log.Println("quics err: ", err)
				return err
			}

			history := &types.FileHistory{}
			if err := history.Decode(val); err != nil {
				log.Println("quics err: ", err)
				return err
			}

			histories = append(histories, history)
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	return histories, nil
}

func (sr *ServerRepository) GetHistoryByAfterPath(afterPath string) (*types.FileHistory, error) {
	key := []byte(PrefixHistory + afterPath)
	history := &types.FileHistory{}

	err := sr.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			log.Println("quics err: ", err)
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			log.Println("quics err: ", err)
			return err
		}

		if err := history.Decode(val); err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	return history, nil
}
