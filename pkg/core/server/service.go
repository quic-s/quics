package server

import (
	"fmt"
	"log"

	"github.com/quic-s/quics/pkg/app"
)

type ServerService struct {
	quics            *app.App
	serverRepository Repository
}

func NewService(app *app.App, serverRepository Repository) *ServerService {
	return &ServerService{
		quics:            app,
		serverRepository: serverRepository,
	}
}

// ListenProtocol is executed when server starts
func (ss *ServerService) ListenProtocol() error {
	fmt.Println("************************************************************")
	fmt.Println("                     Listen Protocol                        ")
	fmt.Println("************************************************************")

	go func() {
		err := ss.quics.Start()
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		err = ss.quics.Close()
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		return
	}()

	return nil
}

// StopServer stop quic-s server
func (ss *ServerService) StopServer() error {
	fmt.Println("************************************************************")
	fmt.Println("                           Stop                             ")
	fmt.Println("************************************************************")

	return nil
}
