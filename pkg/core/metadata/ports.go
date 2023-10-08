package metadata

import "github.com/quic-s/quics/pkg/types"

type Repository interface {
	SaveFileMetadata(fileMetadata *types.FileMetadata) error
	GetFileMetadataByPath(path string) (*types.FileMetadata, error)
}

type Service interface {
}
