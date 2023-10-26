package http

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

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
		uuid := r.URL.Query().Get("uuid")
		afterPath := r.URL.Query().Get("file")

		fileInfo, fileContent, err := sh.sharingService.DownloadFile(uuid, afterPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, fileName := filepath.Split(afterPath)
		w.Header().Set("Alt-Svc", "h3=\":6121\"")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Length", fmt.Sprint(fileInfo.Size))

		n, err := io.Copy(w, fileContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if n != fileInfo.Size {
			http.Error(w, "file is modified", http.StatusInternalServerError)
		}
	}
}
