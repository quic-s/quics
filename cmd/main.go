package main

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/registration"
	"github.com/quic-s/quics/pkg/server"
	"github.com/quic-s/quics/pkg/sync"
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
var Conns []*qp.Connection
var RegistrationHandler *registration.Handler
var SyncHandler *sync.Handler
var ServerHandler *server.Handler

func init() {

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
	SyncHandler = sync.NewSyncHandler(DB)
	ServerHandler = server.NewServerHandler(DB)

	// initialize server password
	Password = config.GetViperEnvVariables("PASSWORD")
	err = ServerHandler.ServerService.UpdatePassword(Password)
	if err != nil {
		log.Fatalln("quis: Error while saving default password: ", err)
	}

	//define system call actions
	SigCh = make(chan os.Signal, 1)
	signal.Notify(SigCh, syscall.SIGINT, syscall.SIGTERM)

	// reset connections
	Conns = make([]*qp.Connection, 0)
}

func main() {

	// HTTP/3
	r := connectRestHandler()
	startHttp3Server(r)

	// protocol
	startQuicsProtocol()

	fmt.Println("************************************************************")
	fmt.Println("                           Start                            ")
	fmt.Println("************************************************************")

	// if pressed ctrl + c, then stop server with closing database
	StopServer()
}

func StopServer() {
	<-SigCh

	// close database
	err := DB.Close()
	if err != nil {
		log.Println("quis: Error while closing database when server is stopped.")
	}
	fmt.Println("quis: Database is closed successfully.")

	// close all conenctions
	for _, conn := range Conns {
		err := conn.Close()
		if err != nil {
			fmt.Println("quis: Closing connection failed with error: ", err)
		}
	}
}
