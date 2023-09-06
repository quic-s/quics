package server

import (
	"fmt"
)

type Service struct {
	serverRepository *Repository
}

func NewServerService(serverRepository *Repository) *Service {
	return &Service{serverRepository: serverRepository}
}

func (serverService *Service) UpdatePassword(password string) error {
	serverService.serverRepository.SetPassword(password)
	return nil
}

func (serverService *Service) GetPassword() string {
	password := serverService.serverRepository.GetPassword()
	return password
}

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
