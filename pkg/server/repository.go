package server

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/types"
	"log"
)

const (
	PASSWORD string = "PASSWORD"
)

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
	SetPassword(password string, client types.Client)
	UpdateDefaultPassword(id string) (*types.Client, error)
}

func NewServerRepository(db *badger.DB) *Repository {
	return &Repository{DB: db}
}

// SetPassword set server password to database from env file
func (serverRepository *Repository) SetPassword(password string) {
	err := serverRepository.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(PASSWORD), []byte(password))
		return err
	})
	if err != nil {
		log.Panicln("Error while setting server: ", err)
	}
}

// GetPassword get server password from database
func (serverRepository *Repository) GetPassword() string {
	var password []byte

	err := serverRepository.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(PASSWORD))
		if err != nil {
			log.Println("quis: Error while getting password: ", err)
		}
		password, err = item.ValueCopy(password)
		return err
	})
	if err != nil {
		log.Fatalln("quics: Error while executing GetPassword() with database.")
	}

	return string(password)
}
