package main

import (
	qp "github.com/quic-s/quics-protocol"
	"log"
	"time"
)

func ClientMessage(msgtype string, message []byte) {
	// initialize client
	quicClient, err := qp.New()
	if err != nil {
		log.Panicln(err)
	}
	log.Println("client Created")
	// start client
	err = quicClient.Dial("172.16.33.124" + ":" + "6122")
	if err != nil {
		log.Panicln(err)
	}
	log.Println("client Connected")
	// send message to server
	quicClient.SendMessage(msgtype, message)

	// delay for waiting message sent to server
	time.Sleep(3 * time.Second)
	quicClient.Close()
}

func main() {
	ClientMessage("get", []byte("test"))
}
