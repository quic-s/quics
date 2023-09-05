package main

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/registration"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	RestApiVersion   string = "v1"
	RestServerApiUri string = "/api/" + RestApiVersion
)

var SigCh chan os.Signal
var Password string
var DB *badger.DB
var RegistrationHandler *registration.Handler

func init() {
	// initialize server password
	Password = config.GetViperEnvVariables("PASSWORD")

	// initialize badger database in .quics/badger directory
	opts := badger.DefaultOptions(config.GetDirPath() + "/badger")
	opts.Logger = nil

	var err error
	DB, err = badger.Open(opts)
	if err != nil {
		log.Println("Error while connecting to the database: ", err)
	}

	// initialize hanlder
	RegistrationHandler = registration.NewRegistrationHandler(DB)

	//define system call actions
	SigCh = make(chan os.Signal, 1)
	signal.Notify(SigCh, syscall.SIGINT, syscall.SIGTERM)
}

func main() {

	// ready to Cobra command
	if err := Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// HTTP/3
	r := connectRestHandler()
	startHttp3Server(r)

	// protocol
	startQuicsProtocol()

	fmt.Println("************************************************************")
	fmt.Println("                           Start                            ")
	fmt.Println("************************************************************")

	StopServer()
}

func StopServer() {
	<-SigCh
	err := DB.Close()
	if err != nil {
		log.Println("quis: Error while closing database when server is stopped.")
	}
	fmt.Println("quis: Database is closed successfully.")
}
