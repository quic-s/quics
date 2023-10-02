package http

import (
	"github.com/gorilla/mux"
	"github.com/quic-s/quics/pkg/core/download"
)

type DownloadHandler struct {
	downloadService download.Service
}

func NewDownloadHandler(downloadService download.Service) *DownloadHandler {
	return &DownloadHandler{
		downloadService: downloadService,
	}
}

func (handler *DownloadHandler) SetupRoutes(r *mux.Router) {
	
}
