package server

import (
	"fmt"
)

// StartServer executes when server starts
func StartServer() {
	fmt.Println("Start server...")
	fmt.Println("Server started successfully.")
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
