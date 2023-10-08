package main

import (
	"fmt"
	"log"

	"github.com/quic-s/quics/pkg/app"
)

const (
	RestApiVersion string = "v1"
	RestApiUri     string = "/api/" + RestApiVersion
)

func main() {

	// initialize application
	// TODO: check whether it is correct
	quics, err := app.New()
	if err != nil {
		log.Println("quics: ", err)
	}

	quics.Start()

	fmt.Println("************************************************************")
	fmt.Println("                           Start                            ")
	fmt.Println("************************************************************")

	quics.Close()
}
