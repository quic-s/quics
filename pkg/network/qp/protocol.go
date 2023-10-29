package qp

import (
	"crypto/tls"
	"fmt"
	"log"

	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/network/qp/connection"
	"github.com/quic-s/quics/pkg/types"
)

type Protocol struct {
	udpaddr            string
	tlsConf            *tls.Config
	initialTransaction func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error
	Proto              *qp.QP
	Pool               *connection.Pool
}

func New(ip string, port int, pool *connection.Pool) (*Protocol, error) {
	// initialize protocol server
	proto, err := qp.New(qp.LOG_LEVEL_ERROR)
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	// initialize certificate for connection
	cert, err := qp.GetCertificate("", "")
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	// initialize tls config for connection with quics protocol
	tlsConfig := &tls.Config{
		Certificates: cert,
		NextProtos:   []string{"quic-s"},
	}

	err = proto.RecvTransactionHandleFunc(types.PING, ping)
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	return &Protocol{
		udpaddr: ":6122",
		tlsConf: tlsConfig,
		Proto:   proto,
		Pool:    pool,
	}, nil
}

// Start starts quics protocol server
func (p *Protocol) Start() error {
	errChan := make(chan error)
	go func() {
		// listen quics protocol with client
		err := p.Proto.ListenWithTransaction(p.udpaddr, p.tlsConf, p.initialTransaction)
		if err != nil {
			log.Println("quics err: ", err)
			errChan <- err
		}
	}()

	err := <-errChan
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}
	fmt.Println("QUIC-S protocol listening successfully.")
	return nil
}

func (p *Protocol) RecvTransactionHandleFunc(transactionName string, handleFunc func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error) error {
	if transactionName == types.REGISTERCLIENT {
		p.initialTransaction = handleFunc
		return nil
	}
	err := p.Proto.RecvTransactionHandleFunc(transactionName, handleFunc)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}
	return nil
}

func ping(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
	log.Println("quics: Ping received ", conn.Conn.RemoteAddr().String())

	data, err := stream.RecvBMessage()
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	request := &types.Ping{}
	if err = request.Decode(data); err != nil {
		log.Println("quics err: ", err)
		return err
	}

	response, err := request.Encode()
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	err = stream.SendBMessage(response)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}
	return nil
}
