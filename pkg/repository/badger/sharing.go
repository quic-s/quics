package badger

import "github.com/dgraph-io/badger/v3"

const (
	PrefixSharing string = "sharing_"
)

type SharingRepository struct {
	db *badger.DB
}

func (sr *SharingRepository) SaveNewDownloadLink(link string) error {
	return nil
}
