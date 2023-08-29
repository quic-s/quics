package history

import "github.com/dgraph-io/badger/v3"

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
}

func NewHistoryRepository(db *badger.DB) *Repository {
	return &Repository{DB: db}
}
