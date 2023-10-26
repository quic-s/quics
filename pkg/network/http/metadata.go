package http

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/metadata"
)

type MetadataHandler struct {
	MetadataService metadata.Service
}

func NewMetadataHandler(metadataService metadata.Service) *MetadataHandler {
	return &MetadataHandler{
		MetadataService: metadataService,
	}
}

func (metadataHandler *MetadataHandler) SetupRoutes(r *mux.Router) {

}
