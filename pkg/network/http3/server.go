package http3

import (
	"net/http"

	"github.com/quic-s/quics/pkg/core/server"
)

type ServerHandler struct {
	ServerService server.Service
}

func NewServerHandler(serverService server.Service) *ServerHandler {
	return &ServerHandler{
		ServerService: serverService,
	}
}

func (sh *ServerHandler) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/server/stop", sh.StopRestServer)
	mux.HandleFunc("/api/v1/server/listen", sh.ListenProtocol)
	mux.HandleFunc("/api/v1/server/logs/clients", sh.ShowClientLogs)
	mux.HandleFunc("/api/v1/server/logs/directories", sh.ShowDirLogs)
	mux.HandleFunc("/api/v1/server/logs/files", sh.ShowFileLogs)
	mux.HandleFunc("/api/v1/server/disconnections/clients", sh.DisconnectClient)
	mux.HandleFunc("/api/v1/server/disconnections/directories", sh.DisconnectDir)
	mux.HandleFunc("/api/v1/server/disconnections/files", sh.DisconnectFile)

	return mux
}

func (sh *ServerHandler) StopRestServer(w http.ResponseWriter, r *http.Request) {
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
	switch r.Method {
	case "POST":
		err := sh.ServerService.ListenProtocol()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) ShowClientLogs(w http.ResponseWriter, r *http.Request) {
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

func (sh *ServerHandler) DisconnectClient(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.DisconnectClient(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) DisconnectDir(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.DisconnectDir(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) DisconnectFile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		all := r.URL.Query().Get("all")
		id := r.URL.Query().Get("id")

		err := sh.ServerService.DisconnectFile(all, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
