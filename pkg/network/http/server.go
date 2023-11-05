package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/quic-s/quics/pkg/config"
	"github.com/quic-s/quics/pkg/core/server"
	"github.com/quic-s/quics/pkg/types"
	"github.com/quic-s/quics/pkg/utils"
)

type ServerHandler struct {
	ServerService server.Service
}

func NewServerHandler(serverService server.Service) *ServerHandler {
	return &ServerHandler{
		ServerService: serverService,
	}
}

func (sh *ServerHandler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/server/stop", sh.StopRestServer)
	mux.HandleFunc("/api/v1/server/listen", sh.ListenProtocol)
	mux.HandleFunc("/api/v1/server/password/set", sh.SetPassword)
	mux.HandleFunc("/api/v1/server/password/reset", sh.ResetPassword)
	mux.HandleFunc("/api/v1/server/logs/clients", sh.ShowClientLogs)
	mux.HandleFunc("/api/v1/server/logs/directories", sh.ShowDirLogs)
	mux.HandleFunc("/api/v1/server/logs/files", sh.ShowFileLogs)
	mux.HandleFunc("/api/v1/server/logs/histories", sh.ShowHistoryLogs)
	mux.HandleFunc("/api/v1/server/remove/clients", sh.RemoveClient)
	mux.HandleFunc("/api/v1/server/remove/directories", sh.RemoveDir)
	mux.HandleFunc("/api/v1/server/remove/files", sh.RemoveFile)
	mux.HandleFunc("/api/v1/server/download/files", sh.DownloadFile)
}

func (sh *ServerHandler) StopRestServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "POST":
		err := sh.ServerService.StopServer()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ListenProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "POST":
		err := sh.ServerService.ListenProtocol()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) SetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "POST":
		body := &types.Server{}

		buf := make([]byte, r.ContentLength)
		n, err := r.Body.Read(buf)
		if err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if int64(n) != r.ContentLength {
			http.Error(w, "failed to read request body", http.StatusInternalServerError)
			return
		}

		err = utils.UnmarshalRequestBody(buf, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = sh.ServerService.SetPassword(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "POST":
		err := sh.ServerService.ResetPassword()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ShowClientLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "GET":
		uuid := r.URL.Query().Get("uuid")

		clients, err := sh.ServerService.ShowClient(uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		response, err := json.Marshal(clients)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		n, err := w.Write(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if n != len(response) {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ShowDirLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "GET":
		afterPath := r.URL.Query().Get("afterpath")

		dirs, err := sh.ServerService.ShowDir(afterPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		response, err := json.Marshal(dirs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		n, err := w.Write(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if n != len(response) {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ShowFileLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "GET":
		afterPath := r.URL.Query().Get("afterpath")

		files, err := sh.ServerService.ShowFile(afterPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		response, err := json.Marshal(files)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		n, err := w.Write(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if n != len(response) {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ShowHistoryLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "GET":
		afterPath := r.URL.Query().Get("afterpath")

		histories, err := sh.ServerService.ShowHistory(afterPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		response, err := json.Marshal(histories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		n, err := w.Write(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if n != len(response) {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) RemoveClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "POST":
		afterPath := r.URL.Query().Get("afterpath")

		err := sh.ServerService.RemoveClient(afterPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) RemoveDir(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "POST":
		afterPath := r.URL.Query().Get("afterpath")

		err := sh.ServerService.RemoveDir(afterPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) RemoveFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "POST":
		afterPath := r.URL.Query().Get("afterpath")

		err := sh.ServerService.RemoveFile(afterPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "GET":
		afterPath := r.URL.Query().Get("afterpath")
		timestamp, err := strconv.Atoi(r.URL.Query().Get("timestamp"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fileInfo, fileContent, err := sh.ServerService.DownloadFile(afterPath, uint64(timestamp))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, fileName := filepath.Split(afterPath)
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
