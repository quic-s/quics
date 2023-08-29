package server

import (
	"fmt"
	"github.com/quic-go/quic-go"
	qp "github.com/quic-s/quics-protocol"
	pb "github.com/quic-s/quics-protocol/proto/v1" // defines message contents
	"log"
)

// StartServer Start quics server
func StartServer() error {
	proto, err := qp.New()
	if err != nil {
		return err
	}

	proto.RecvMessage(func(conn quic.Connection, message *pb.Message) {
		// TODO
		log.Println(message.Message)
	})

	go func() {

		log.Println("Start to listening protocol...")

		err := proto.Listen("0.0.0.0", 6122)
		if err != nil {
			log.Fatalf("error with: %s", err)
		}

		//err := http.ListenAndServe(string(rune(config.RuntimeConf.Server.Port)), nil)
		//if err != nil {
		//	log.Println(err)
		//}
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
