package badger

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/utils"
)

type Badger struct {
	db *badger.DB
}

func NewBadgerRepository() (*Badger, error) {
	// initialize badger database in .quics/badger directory
	opts := badger.DefaultOptions(utils.GetQuicsDirPath() + "/badger")
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		log.Println("quics: Error while connecting to the database: ", err)
		return nil, err
	}

	return &Badger{
		db: db,
	}, nil
}

func (b *Badger) Close() error {
	err := b.db.Close()
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (b *Badger) NewHistoryRepository() *HistoryRepository {
	return &HistoryRepository{
		db: b.db,
	}
}

func (b *Badger) NewMetadataRepository() *MetadataRepository {
	return &MetadataRepository{
		db: b.db,
	}
}

func (b *Badger) NewRegistrationRepository() *RegistrationRepository {
	return &RegistrationRepository{
		db: b.db,
	}
}

func (b *Badger) NewServerRepository() *ServerRepository {
	return &ServerRepository{
		db: b.db,
	}
}

func (b *Badger) NewSharingRepository() *SharingRepository {
	return &SharingRepository{
		db: b.db,
	}
}

func (b *Badger) NewSyncRepository() *SyncRepository {
	return &SyncRepository{
		db: b.db,
	}
}
