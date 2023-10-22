package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/quic-s/quics/pkg/core/sharing"
)

type SharingHandler struct {
	sharingService sharing.Service
}

func NewSharingHandler(sharingService sharing.Service) *SharingHandler {
	return &SharingHandler{
		sharingService: sharingService,
	}
}

func (sh *SharingHandler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/download/files", sh.DownloadFile)
}

func (sh *SharingHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uuid := r.URL.Query().Get("id")
		afterPath := r.URL.Query().Get("file")

		file, fileInfo, err := sh.sharingService.DownloadFile(uuid, afterPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename=")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprint(fileInfo.Size()))

		_, err = io.Copy(w, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
