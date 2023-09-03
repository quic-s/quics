package server

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/registeration"
)

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
	SetDefaultPassword(newId []byte, client registeration.Client)
	UpdateDefaultPassword(id string) (*registeration.Client, error)
}
