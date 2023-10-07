package badger

import "github.com/dgraph-io/badger/v3"

type ServerRepository struct {
	db *badger.DB
}
