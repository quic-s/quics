package client

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
)

type Handler struct {
	db            *badger.DB
	clientService *Service
}

func NewClientHandler(db *badger.DB) *Handler {
	clientRepository := NewClientRepository(db)
	clientService := NewClientService(clientRepository)
	return &Handler{db: db, clientService: clientService}
}

func (clientHandler *Handler) SetupRoutes(r *mux.Router) {
	// 2.1 create (connect) new client
	r.HandleFunc("/", clientHandler.createClientHandler).Methods(http.MethodPost)

	// 2.2 get the status of a client
	//r.HandleFunc("/{clientId}", clientHandler.getClientHandler).Methods(http.MethodGet)

	// 2.3 get all clients
	//r.HandleFunc("/", clientHandler.getAllClientsHandler).Methods(http.MethodGet)

	// 2.4 disconnect client
	//r.HandleFunc("/{clientId}", clientHandler.disconnectClient).Methods(http.MethodPost)
}

// createClientHandler implements the API of 2.1 create (connect new client)
func (clientHandler *Handler) createClientHandler(w http.ResponseWriter, r *http.Request) {
	var request CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Faield to create client because of wrong request data", http.StatusBadRequest)
		log.Fatalf("Error while decoding request data: %s", err)
		return
	}

	uuid, err := clientHandler.clientService.CreateNewClient(&request.Ip)
	if err != nil {
		http.Error(w, "Faield to create client", http.StatusInternalServerError)
		log.Fatalf("Error while returning response data: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uuid)
}
