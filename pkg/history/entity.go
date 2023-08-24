package history

import (
	"github.com/quic-s/quics/pkg/client"
	"github.com/quic-s/quics/pkg/metadata"
)

type FileHistory struct {
	Id     int
	Date   string
	client client.Client
	File   metadata.FileMetadata // History must have file metadata at the point that client wanted in time
}
