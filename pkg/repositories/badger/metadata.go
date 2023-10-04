package badger

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
)

const (
	PrefixMetadata string = "metadata_"
)

type MetadataRepository struct {
}

func NewMetadataRepository() *MetadataRepository {
	return &MetadataRepository{}
}

// SaveFileMetadata saves new file metadata to badger
func (repository *MetadataRepository) SaveFileMetadata(path string, fileMetadata types.FileMetadata) error {
	key := []byte(PrefixMetadata + path)

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, fileMetadata.Encode())
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

// GetFileMetadataByPath find and return certain file metadata by path
func (repository *MetadataRepository) GetFileMetadataByPath(path string) *types.FileMetadata {
	key := []byte(PrefixMetadata + path)
	var fileMetadata *types.FileMetadata

	err := db.View(func(txn *badger.Txn) error {
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
