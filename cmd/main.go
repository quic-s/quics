package cmd

import (
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-s/quics/pkg/utils/env"
	"github.com/quic-s/quics/pkg/utils/security"
	"log"
	"os"
	"path/filepath"

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
	quicConfig := quic.Config{}

	server := &http3.Server{
		Addr:       ":8080",
		QuicConfig: &quicConfig,
		Handler:    r,
	}

	// get directory path for certification
	quicsDir := env.GetDirPath()
	certFileDir := filepath.Join(quicsDir, env.GetViperEnvVariables("QUICS_CERT_NAME"))
	keyFileDir := filepath.Join(quicsDir, env.GetViperEnvVariables("QUICS_KEY_NAME"))

	// load the certificate and the key from the files
	_, err := tls.LoadX509KeyPair(certFileDir, keyFileDir)
	if err != nil {
		security.SecurityFiles()
	}

	fmt.Println("Starting HTTP/3 server...")

	// TODO: need to link cretification files
	log.Fatal(server.ListenAndServeTLS(certFileDir, keyFileDir))
}
