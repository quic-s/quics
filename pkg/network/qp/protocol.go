package qp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/pkg/network/qp/connection"
	"github.com/quic-s/quics/pkg/types"
)

type Protocol struct {
	udpaddr            *net.UDPAddr
	tlsConf            *tls.Config
	initialTransaction func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error
	Proto              *qp.QP
	Pool               *connection.Pool
}

func New(ip string, port int, pool *connection.Pool) (*Protocol, error) {
	// initialize protocol server
	proto, err := qp.New(qp.LOG_LEVEL_INFO)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	// initialize udp server address
	UDPAddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: port,
	}

	// initialize certificate for connection
	cert, err := qp.GetCertificate("", "")
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	// initialize tls config for connection with quics protocol
	tlsConfig := &tls.Config{
		Certificates: cert,
		NextProtos:   []string{"quic-s"},
	}

	return &Protocol{
		udpaddr: UDPAddr,
		tlsConf: tlsConfig,
		Proto:   proto,
		Pool:    pool,
	}, nil
}

// startQuicsProtocol starts quics protocol server
func (p *Protocol) Start() error {
	errChan := make(chan error)
	go func() {
		// listen quics protocol with client
		err := p.Proto.ListenWithTransaction(p.udpaddr, p.tlsConf, p.initialTransaction)
		if err != nil {
			log.Println("quics: ", err)
			errChan <- err
		}
	}()

	err := <-errChan
	if err != nil {
		log.Println("quics: ", err)
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
		log.Println("quics: ", err)
		return err
	}
	return nil
}

// // connectProtocolHandler connects handler to quics protocol
// func setHandler(proto *qp.QP) {
// 	err := proto.RecvTransactionHandleFunc(types.REGISTERCLIENT, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
// 	})
// 	if err != nil {
// 		log.Println("quics: ", err)
// 	}

// 	err = proto.RecvTransactionHandleFunc(types.REGISTERROOTDIR, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
// 	})
// 	if err != nil {
// 		log.Println("quics: ", err)
// 	}

// 	err = proto.RecvTransactionHandleFunc(types.SYNCROOTDIR, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
// 	})
// 	if err != nil {
// 		log.Println("quics: ", err)
// 	}

// 	err = proto.RecvTransactionHandleFunc(types.GETROOTDIRS, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
// 	})
// 	if err != nil {
// 		log.Println("quics: ", err)
// 	}

// 	err = proto.RecvTransactionHandleFunc(types.PLEASESYNC, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
// 	})
// 	if err != nil {
// 		log.Println("quics: ", err)
// 	}
// }

// func openMustSyncTransaction(conn *qp.Connection) {
// 	var err error

// 	// must sync transaction
// 	// 1. (server) Open transaction
// 	// 2. (server) MustSyncReq with file metadata to all registered clients without where the file come from
// 	// 3. (client) MustSyncRes if file update is available
// 	// 3-1. (server) If all request data are exist, then go to step 4
// 	// 3-2. (server) If not, then this transaction should be closed
// 	// 4. (server) GiveYouReq for giving file contents
// 	// 5. (client) GiveYouRes
// 	err = conn.OpenTransaction(types.MUSTSYNC, func(stream *qp.Stream, transactionName string, transactionId []byte) error {
// 		// TODO: implement must sync transaction process

// 		return nil
// 	})
// 	if err != nil {
// 		log.Println("quics: ", err)
// 	}
// }
