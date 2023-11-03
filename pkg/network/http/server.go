package http

import (
	"net/http"

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

		err := utils.UnmarshalRequestBody(r, body)
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
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.ShowClientLogs(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ShowDirLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "GET":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.ShowDirLogs(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ShowFileLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "GET":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.ShowFileLogs(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ShowHistoryLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "GET":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.ShowHistoryLogs(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) RemoveClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", "h3=\":"+config.GetViperEnvVariables("REST_SERVER_H3_PORT")+"\"")
	switch r.Method {
	case "POST":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.RemoveClient(all, id)
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
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.RemoveDir(all, id)
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
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.RemoveFile(all, id)
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
		path := r.URL.Query().Get("path")
		version := r.URL.Query().Get("version")
		target := r.URL.Query().Get("target")

		err := sh.ServerService.DownloadFile(path, version, target)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
