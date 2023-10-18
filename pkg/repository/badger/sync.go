package badger

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
)

const (
	PrefixFile     string = "file_"
	PrefixConflict string = "conflict_"
)

type SyncRepository struct {
	db *badger.DB
}

func (sr *SyncRepository) SaveRootDir(afterPath string, rootDir *types.RootDirectory) error {
	key := []byte(PrefixRootDir + afterPath)

	err := sr.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, rootDir.Encode())
		return err
	})
	if err != nil {
		log.Println("quics: (SaveClient) ", err)
		return err
	}
	return nil
}

func (sr *SyncRepository) GetRootDirByPath(afterPath string) (*types.RootDirectory, error) {
	key := []byte(PrefixRootDir + afterPath)

	rootDir := &types.RootDirectory{}
	err := sr.db.View(func(txn *badger.Txn) error {
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

func (sr *SyncRepository) GetAllRootDir() ([]types.RootDirectory, error) {
	rootDirs := []types.RootDirectory{}
	err := sr.db.View(func(txn *badger.Txn) error {
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

			rootDir := types.RootDirectory{}
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

// IsExistFileByPath checks if file exists by file path
func (sr *SyncRepository) IsExistFileByPath(afterPath string) (bool, error) {
	key := []byte(PrefixFile + afterPath)

	err := sr.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		return err
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetFileByPath gets file by file path
func (sr *SyncRepository) GetFileByPath(path string) (*types.File, error) {
	key := []byte(PrefixFile + path)
	file := &types.File{}

	err := sr.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		file = &types.File{}
		if err := file.Decode(val); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return file, nil
}

// SaveFileByPath saves new file to badger
func (sr *SyncRepository) SaveFileByPath(path string, file *types.File) error {
	key := []byte(PrefixFile + path)

	err := sr.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, file.Encode())
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

// GetAllFiles gets all files
func (sr *SyncRepository) GetAllFiles(prefix string) ([]types.File, error) {
	key := []byte(PrefixFile + prefix)
	files := []types.File{}

	err := sr.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(key); it.ValidForPrefix(key); it.Next() {
			item := it.Item()

			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			file := types.File{}
			if err := file.Decode(val); err != nil {
				return err
			}

			files = append(files, file)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// UpdateContentsExisted updates contents existed flag (if exist then true, or not then false)
func (sr *SyncRepository) UpdateContentsExisted(path string, contentsExisted bool) error {
	key := []byte(PrefixFile + path)

	err := sr.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		file := &types.File{}
		if err := file.Decode(val); err != nil {
			return err
		}

		file.ContentsExisted = contentsExisted

		if err := txn.Set(key, file.Encode()); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (sr *SyncRepository) UpdateFile(file *types.File) error {
	key := []byte(PrefixFile + file.AfterPath)

	err := sr.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set(key, file.Encode()); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (sr *SyncRepository) UpdateConflict(afterpath string, conflict *types.Conflict) error {
	key := []byte(PrefixConflict + afterpath)

	err := sr.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, conflict.Encode())
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (sr *SyncRepository) GetConflict(afterpath string) (*types.Conflict, error) {
	key := []byte(PrefixConflict + afterpath)
	conflict := &types.Conflict{}

	err := sr.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		if err := conflict.Decode(val); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return conflict, nil
}

func (sr *SyncRepository) GetConflictList(rootDirs []string) ([]types.Conflict, error) {
	conflictMetadataList := []types.Conflict{}
	for _, rootDir := range rootDirs {
		key := []byte(PrefixConflict + rootDir)

		sr.db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchSize = 10
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Seek(key); it.ValidForPrefix(key); it.Next() {
				item := it.Item()

				val, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}

				conflictMetadata := types.Conflict{}
				if err := conflictMetadata.Decode(val); err != nil {
					return err
				}

				conflictMetadataList = append(conflictMetadataList, conflictMetadata)
			}

			return nil
		})
	}
	return conflictMetadataList, nil
}

func (sr *SyncRepository) DeleteConflict(afterpath string) error {
	key := []byte(PrefixConflict + afterpath)

	err := sr.db.DropPrefix(key)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SyncRepository) ErrKeyNotFound() error {
	return badger.ErrKeyNotFound
}
