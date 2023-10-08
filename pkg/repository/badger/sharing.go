package badger

import "github.com/dgraph-io/badger/v3"

type SharingRepository struct {
	db *badger.DB
}
