package main

import (
	"fmt"
	"github.com/quic-s/quics/pkg/app"
	"log"
)

const (
	RestApiVersion string = "v1"
	RestApiUri     string = "/api/" + RestApiVersion
)

func main() {

	// initialize application
	// TODO: check whether it is correct
	quics, err := app.Initialize()
	if err != nil {
		log.Println("quics: ", err)
		return
	}

	// initialize adapters
	quics.InitAdapters()

	// start HTTP/3 server
	r := connectRestHandler()
	startHttp3Server(r)

	// start quics protocol server
	startQuicsProtocol()

	fmt.Println("************************************************************")
	fmt.Println("                           Start                            ")
	fmt.Println("************************************************************")

	quics.Close()
}
