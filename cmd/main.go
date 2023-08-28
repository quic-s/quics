package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/client"
)

const Version = "v1"
const uri = "/api/" + Version

func main() {

	// FIXME: this is for test. Delete after.
	fmt.Println("database: ", config.RuntimeConf.Database.Path)

	// ready to Cobra command
	if err := Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// initialize badger database
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatalf("Error while connecting to the database: %s", err)
	}
	defer db.Close()

	r := connectHandler(db)
	startHttp3Server(r)
}

// connectHandler creates mux router and connect handler to router
func connectHandler(db *badger.DB) *mux.Router {
	r := mux.NewRouter()

	clientHandler := client.NewClientHandler(db)
	clientHandler.SetupRoutes(r.PathPrefix(uri + "/clients").Subrouter())

	return r
}

func startHttp3Server(r *mux.Router) {
	server := &http3.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("Starting HTTP/3 server...")

	// TODO: need to link cretification files
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
