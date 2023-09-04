package main

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/registration"
	"log"
	"os"
)

const (
	RestApiVersion   = "v1"
	RestServerApiUri = "/api/" + RestApiVersion
)

var DB *badger.DB
var RegistrationHandler *registration.Handler

func init() {
	// initialize badger database in .quics/badger directory
	opts := badger.DefaultOptions(config.GetDirPath() + "/badger")
	opts.Logger = nil

	DB, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("Error while connecting to the database: %s", err)
	}

	RegistrationHandler = registration.NewRegistrationHandler(DB)
}

func main() {

	// TODO: when server stopped, call these command below
	//defer db.Close()

	// ready to Cobra command
	if err := Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	setDefaultPassword()

	// HTTP/3
	r := connectRestHandler()
	startHttp3Server(r)

	// protocol
	startQuicsProtocol()

	fmt.Println("************************************************************")
	fmt.Println("                           Start                            ")
	fmt.Println("************************************************************")

	// If press Ctrl + C, then stop server
	select {}
}

// setDefaultPassword sets default password of server from env for accessing client
func setDefaultPassword() {

}
