package badger

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/config"
	"log"
)

var db *badger.DB

// NewBadgerDB initializes badger database
func NewBadgerDB() {
	var err error

	// initialize badger database in .quics/badger directory
	opts := badger.DefaultOptions(config.GetQuicsDirPath() + "/badger")
	opts.Logger = nil
	db, err = badger.Open(opts)
	if err != nil {
		log.Println("quics: Error while connecting to the database: ", err)
	}
}

// CloseBadgerDB closes badger database
func CloseBadgerDB() {
	err := db.Close()
	if err != nil {
		log.Println("quics: Error while closing database when server is stopped.")
	}
}
