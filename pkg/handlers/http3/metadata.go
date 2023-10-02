package http3

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/metadata"
)

type MetadataHandler struct {
	metadataService metadata.Service
}

func NewMetadataHandler(metadataService metadata.Service) *MetadataHandler {
	return &MetadataHandler{
		metadataService: metadataService,
	}
}

func (metadataHandler *MetadataHandler) SetupRoutes(r *mux.Router) {
	
}
