package history

import (
	"github.com/quic-s/quics/pkg/metadata"
)

type FileHistory struct {
	Id   uint64
	Date string
	Uuid string
	Path string                // path of stored the file with history
	File metadata.FileMetadata // must have file metadata at the point that client wanted in time
}
