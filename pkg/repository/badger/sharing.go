package badger

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
)

const (
	PrefixSharing string = "sharing_"
)

type SharingRepository struct {
	db *badger.DB
}

func (sr *SharingRepository) SaveLink(sharing *types.Sharing) error {
	key := []byte(PrefixSharing + sharing.Link)

	err := sr.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, sharing.Encode())
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (sr *SharingRepository) GetLink(link string) (*types.Sharing, error) {
	key := []byte(PrefixSharing + link)

	sharing := &types.Sharing{}

	err := sr.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		if err := sharing.Decode(val); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return sharing, nil
}

func (sr *SharingRepository) DeleteLink(link string) error {
	key := []byte(PrefixSharing + link)

	err := sr.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (sr *SharingRepository) UpdateLink(sharing *types.Sharing) error {
	key := []byte(PrefixSharing + sharing.Link)

	err := sr.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, sharing.Encode())
		return err
	})
	if err != nil {
		return err
	}

	return nil
}
