package badger

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
)

const (
	PrefixHistory string = "history_"
)

type HistoryRepository struct {
}

func NewHistoryRepository() *HistoryRepository {
	return &HistoryRepository{}
}

// SaveNewFileHistory creates the history with file metadata
func (repository *HistoryRepository) SaveNewFileHistory(path string, fileHistory types.FileHistory) error {
	key := []byte(PrefixHistory + path)

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, fileHistory.Encode())
		return err
	})
	if err != nil {
		return err
	}

	return nil
}
