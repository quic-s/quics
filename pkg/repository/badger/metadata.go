package badger

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
)

const (
	PrefixMetadata string = "metadata_"
)

type MetadataRepository struct {
	db *badger.DB
}

// SaveFileMetadata saves new file metadata to badger
func (mr *MetadataRepository) SaveFileMetadata(path string, fileMetadata types.FileMetadata) error {
	key := []byte(PrefixMetadata + path)

	err := mr.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, fileMetadata.Encode())
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

// GetFileMetadataByPath find and return certain file metadata by path
func (mr *MetadataRepository) GetFileMetadataByPath(path string) *types.FileMetadata {
	key := []byte(PrefixMetadata + path)
	fileMetadata := &types.FileMetadata{}

	err := mr.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		fileMetadata = &types.FileMetadata{}
		if err := fileMetadata.Decode(val); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil
	}

	return fileMetadata
}
