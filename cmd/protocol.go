package main

import (
	"crypto/tls"
	"fmt"
	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics-protocol/pkg/utils/fileinfo"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/registration"
	"io"
	"log"
	"net"
	"strconv"
)

const (
	REGISTER_CLIENT  = "CLIENT"
	REGISTER_ROOTDIR = "ROOTDIR"

	FILES             = "FILE"
	FILIES_DOWNLOAD   = "DOWNLOAD"
	FILES_CREATE      = "CREATE"
	FILES_DELETE      = "DELETE"
	FILES_RENAME      = "RENAME"
	FILES_SYNC_RESCAN = "RESCAN"
	FILES_HISTORY     = "HISTORY"
	FILES_SHARING     = "SHARING"
)

func connectProtocolHandler(proto *qp.QP) {
	// [REGISTER] CLIENT: register client from client
	err := proto.RecvMessageHandleFunc(REGISTER_CLIENT, func(conn *qp.Connection, msgType string, data []byte) {
		// decode request data
		var request registration.RegisterClientRequest
		if err := request.Decode(data); err != nil {
			fmt.Println("[QUICS] (RegisterClientRequest) Error while decoding request data")
			return
		}

		clientUuid, err := RegistrationHandler.RegistrationService.CreateNewClient(request.Uuid, request.ClientPassword, conn.Conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("[QUICS] (RegisterClientRequest) Error while creating new client: ", err)
		}

		response := registration.RegisterClientResponse{
			Uuid: clientUuid,
		}

		// encode response data of request
		encodedResponse, err := registration.RegisterClientResponse.Encode(response)
		if err != nil {
			fmt.Println("[QUICS] (RegisterClientResponse) Error while encoding response data")
		}

		err = conn.SendMessage(REGISTER_CLIENT, encodedResponse)
		if err != nil {
			fmt.Println("[QUICS-PROTOCOL] (ReigsterClientResponse) Error while sending message to client")
		}
	})
	if err != nil {
		log.Printf("[QUICS-PROTOCOL] Error while receiving message from client: %s\n", err)
	}

	// register root directory from client (RegisterRootDirRequest)
	err = proto.RecvMessageHandleFunc(REGISTER_ROOTDIR, func(conn *qp.Connection, msgType string, data []byte) {

	})
	if err != nil {
		log.Printf("[QUICS-PROTOCOL] Error while receiving message from client: %s\n", err)
	}

	// please file message from client
	err = proto.RecvMessageHandleFunc(FILES, func(conn *qp.Connection, msgType string, data []byte) {

	})
	if err != nil {
		log.Printf("[QUICS-PROTOCOL] Error while receiving message from client: %s\n", err)
	}

	// please sync (file) from client
	err = proto.RecvFileMessageHandleFunc(FILES, func(conn *qp.Connection, fileMsgType string, msgData []byte, fileInfo *fileinfo.FileInfo, fileReader io.Reader) {

	})
	if err != nil {
		log.Printf("[QUICS-PROTOCOL] Error while receiving file and message from client: %s\n", err)
	}

}

func startQuicsProtocol() {
	proto, err := qp.New(qp.LOG_LEVEL_INFO)
	if err != nil {
		log.Fatalf("Error while creating connection protocol: %s", err)
	}

	portStr := config.GetViperEnvVariables("QUICS_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error while getting port number: %s", err)
	}

	UDPAddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: port,
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quics-protocol"},
	}

	// TODO: conn instance를 리스트화 시켜서 이후 브로드 캐스트 시 사용할 수 있도록 해야 한다.
	conn := func(conn *qp.Connection) {
		log.Println("[QUICS-PROTOCOL] Created new connection: ", conn.Conn.RemoteAddr().String())
	}

	// connect handler
	connectProtocolHandler(proto)

	go func() {
		// protocol
		err = proto.Listen(UDPAddr, tlsConfig, conn)
		if err != nil {
			log.Fatalf("Error while listening protocol: %s", err)
		}
	}()

	fmt.Println("QUIC-S protocol listening successfully.")
}
