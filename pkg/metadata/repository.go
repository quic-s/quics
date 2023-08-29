package metadata

import "github.com/dgraph-io/badger/v3"

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
}

func NewMetadataRepository(db *badger.DB) *Repository {
	return &Repository{DB: db}
}
