package badger

import (
	"strconv"

	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
)

const (
	PrefixHistory string = "history_"
)

type HistoryRepository struct {
	db *badger.DB
}

// SaveNewFileHistory creates the history with file metadata
func (hr *HistoryRepository) SaveNewFileHistory(afterPath string, fileHistory *types.FileHistory) error {
	key := []byte(PrefixHistory + afterPath + "_" + strconv.FormatUint(fileHistory.Timestamp, 10))

	err := hr.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, fileHistory.Encode())
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

// GetFileHistory returns the history of the file
func (hr *HistoryRepository) GetFileHistory(afterPath string, timestamp uint64) (*types.FileHistory, error) {
	key := []byte(PrefixHistory + afterPath + "_" + strconv.FormatUint(timestamp, 10))
	fileHistory := &types.FileHistory{}

	err := hr.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = fileHistory.Decode(val)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return fileHistory, nil
}
