package server

import (
	"fmt"
	"github.com/quic-s/quics/config"
	"log"
	"net/http"
)

// StartServer Start quics server
func StartServer() error {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Start server...")
		fmt.Println(w)
	})

	go func() {
		err := http.ListenAndServe(string(rune(config.RuntimeConf.Server.Port)), nil)
		if err != nil {
			log.Println(err)
		}
	}()

	fmt.Println("Server started successfully.")
	fmt.Println("Press Ctrl + C to stop the server.")

	// If press Ctrl + C, then stop server
	select {}
}

// StopServer Stop quics server
func StopServer() {
	fmt.Println("Stop server...")
	fmt.Println("Server stopped successfully.")
}

// RebootServer Reboot quics server
func RebootServer() {
	fmt.Println("Stop server...")
	fmt.Println("Server rebooted successfully.")
}

// ShutdownServer Shutdown quics server
func ShutdownServer() {
	fmt.Println("Shutdown server...")
	fmt.Println("Server shutdown successfully.")
}
