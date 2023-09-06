package sync

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
)

const (
	PrefixFile string = "file_"
)

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
	GetFilesByRootDirPath(rootDirPath string) []types.File
}

func NewSyncRepository(db *badger.DB) *Repository {
	return &Repository{DB: db}
}

func (syncRepository *Repository) GetFileByPath(path string) *types.File {
	key := []byte(PrefixFile + path)
	var file *types.File

	err := syncRepository.DB.View(func(txn *badger.Txn) error {
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
		return nil
	}

	return file
}
