package sharing

import "github.com/dgraph-io/badger/v3"

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
}

func NewSharingRepository(db *badger.DB) *Repository {
	return &Repository{DB: db}
}
