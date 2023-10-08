package server

import (
	"fmt"
)

type ServerService struct {
	serverRepository Repository
}

func NewService(serverRepository Repository) *ServerService {
	return &ServerService{
		serverRepository: serverRepository,
	}
}

// StartServer executes when server starts
func (ss *ServerService) StartServer() {
	fmt.Println("Start server...")
	fmt.Println("Server started successfully.")
}

// StopServer Stop quics server
func (ss *ServerService) StopServer() {
	fmt.Println("Stop server...")
	fmt.Println("Server stopped successfully.")
}

// RebootServer Reboot quics server
func (ss *ServerService) RebootServer() {
	fmt.Println("Stop server...")
	fmt.Println("Server rebooted successfully.")
}

// ShutdownServer Shutdown quics server
func (ss *ServerService) ShutdownServer() {
	fmt.Println("Shutdown server...")
	fmt.Println("Server shutdown successfully.")
}
