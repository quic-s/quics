package history

import (
	"github.com/quic-s/quics/pkg/metadata"
)

type FileHistory struct {
	Id   uint64                `json:"id"`
	Date string                `json:"date"`
	Uuid string                `json:"uuid"`
	File metadata.FileMetadata `json:"file"` // must have file metadata at the point that client wanted in time
}
