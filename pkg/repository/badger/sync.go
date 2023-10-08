package badger

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
)

const (
	PrefixFile string = "file_"
)

type SyncRepository struct {
	db *badger.DB
}

// GetFileByPath gets file by file path
func (sr *SyncRepository) GetFileByPath(path string) (*types.File, error) {
	key := []byte(PrefixFile + path)
	var file *types.File

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
func (sr *SyncRepository) SaveFileByPath(path string, file types.File) error {
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
func (sr *SyncRepository) GetAllFiles() []types.File {
	var files []types.File

	err := sr.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(PrefixFile)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			file := &types.File{}
			if err := file.Decode(val); err != nil {
				return err
			}

			files = append(files, *file)
		}

		return nil
	})
	if err != nil {
		return nil
	}

	return files
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
