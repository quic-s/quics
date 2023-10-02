package server

import (
	"fmt"
)

type MyServerService struct {
	serverRepository Repository
}

func NewServerService(serverRepository Repository) *MyServerService {
	return &MyServerService{
		serverRepository: serverRepository,
	}
}

// StartServer executes when server starts
func (service *MyServerService) StartServer() {
	fmt.Println("Start server...")
	fmt.Println("Server started successfully.")
}

// StopServer Stop quics server
func (service *MyServerService) StopServer() {
	fmt.Println("Stop server...")
	fmt.Println("Server stopped successfully.")
}

// RebootServer Reboot quics server
func (service *MyServerService) RebootServer() {
	fmt.Println("Stop server...")
	fmt.Println("Server rebooted successfully.")
}

// ShutdownServer Shutdown quics server
func (service *MyServerService) ShutdownServer() {
	fmt.Println("Shutdown server...")
	fmt.Println("Server shutdown successfully.")
}
