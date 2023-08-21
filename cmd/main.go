package main

import (
	"fmt"
	"github.com/quic-s/quics/config"
	"log"
	"os"
)

func main() {

	// FIXME: this is for test. Delete after.
	fmt.Println("database: ", config.RuntimeConf.Database.Path)

	if err := Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
