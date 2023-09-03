package main

import (
	"crypto/tls"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	qp "github.com/quic-s/quics-protocol"
	pb "github.com/quic-s/quics-protocol/proto/v1"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/registeration"
	"github.com/quic-s/quics/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const restApiVersion = "v1"
const restUri = "/api/" + restApiVersion

func main() {
	// initialize badger database in .quics/badger directory
	opts := badger.DefaultOptions(config.GetDirPath() + "/badger")
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("Error while connecting to the database: %s", err)
	}
	defer db.Close()

	// ready to Cobra command
	if err := Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	setDefaultPassword(db)
	r := connectHandler(db)
	startHttp3Server(r)
	startQuicsProtocol()

	fmt.Println("************************************************************")
	fmt.Println("                           Start                            ")
	fmt.Println("************************************************************")

	// If press Ctrl + C, then stop server
	select {}
}

// setDefaultPassword sets default password of server from env for accessing client
func setDefaultPassword(db *badger.DB) {

}

// connectHandler creates mux router and connect handler to router
func connectHandler(db *badger.DB) *mux.Router {
	r := mux.NewRouter()

	clientHandler := registeration.NewRegistrationHandler(db)
	clientHandler.SetupRoutes(r.PathPrefix(restUri + "/clients").Subrouter())

	return r
}

func startHttp3Server(r *mux.Router) {
	quicConfig := quic.Config{}

	server := &http3.Server{
		Addr:       config.GetViperEnvVariables("REST_SERVER_ADDR"),
		QuicConfig: &quicConfig,
		Handler:    r,
	}

	// get directory path for certification
	quicsDir := config.GetDirPath()
	certFileDir := filepath.Join(quicsDir, config.GetViperEnvVariables("QUICS_CERT_NAME"))
	keyFileDir := filepath.Join(quicsDir, config.GetViperEnvVariables("QUICS_KEY_NAME"))

	// load the certificate and the key from the files
	_, err := tls.LoadX509KeyPair(certFileDir, keyFileDir)
	if err != nil {
		utils.SecurityFiles()
	}

	go func() {
		log.Fatal(server.ListenAndServeTLS(certFileDir, keyFileDir))
	}()

	fmt.Println("HTTP/3 server started successfully.")
}

func startQuicsProtocol() error {
	proto, err := qp.New()
	if err != nil {
		return err
	}

	err = proto.RecvMessage(func(conn quic.Connection, message *pb.Message) {
		log.Println(message.Message)
	})
	if err != nil {
		return err
	}

	go func() {
		portStr := config.GetViperEnvVariables("QUICS_PORT")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Error while getting port number: %s", err)
		}

		// protocol
		err = proto.Listen("0.0.0.0", port)
		if err != nil {
			log.Fatalf("Error while listening protocol: %s", err)
		}
	}()

	fmt.Println("QUIC-S protocol listened successfully.")

	return nil
}
