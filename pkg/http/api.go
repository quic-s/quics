package http

import (
	"fmt"
	"log"
	"net/http"
)

func StartRestServer() {
	go func() {

	}()

	handler := setupHandler()
	config := quis.Config{}
	server := http3.Server{
		Handler:    handler,
		QuicConfig: &qconf,
		Addr:       ":6121",
	}
}

// setupHandler
// setup handler for RESTful APIs
func setupHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%#v\n", r)
		_, err := w.Write([]byte("Start rest server..."))
		if err != nil {
			log.Fatalf("Error while showing log at starting server point: %s", err)
		}
	})

	return mux
}
