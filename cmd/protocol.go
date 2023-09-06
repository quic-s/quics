package main

import (
	"crypto/tls"
	"fmt"
	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/types"
	"log"
	"net"
	"strconv"
)

const (
	REGISTER_CLIENT            = "CLIENT"
	REGISTER_DISCONNECT_CLIENT = "NOTCLIENTANYMORE"

	REGISTER_ROOTDIR            = "LOCALROOT"
	REGISTER_SYNC_ROOTDIR       = "REMOTEROOT"
	REGISTER_DISCONNECT_ROOTDIR = "NOTROOTDIRANYMORE"

	FILES_SYNC_RESCAN     = "RESCAN"
	FILES_SYNC_PLEASESYNC = "PLEASESYNC"
	FILES_SYNC_PLEASEFILE = "PLEASEFILE"

	FILES_SHARING = "SHARING"
)

func connectProtocolHandler(proto *qp.QP) {

	// [REGISTER] CLIENT: register root directory (local to remote)
	err := proto.RecvMessageWithResponseHandleFunc(REGISTER_ROOTDIR, func(conn *qp.Connection, msgType string, data []byte) []byte {
		// decode request data
		var request types.RegisterRootDirRequest
		if err := request.Decode(data); err != nil {
			log.Println("quics: Error while decoding request data")
			return []byte("FAIL")
		}

		err := RegistrationHandler.RegistrationService.RegisterRootDir(request)
		if err != nil {
			log.Println("quics: (RegisterRootDirRequest) Error while creating root directory: ", err)
			return []byte("FAIL")
		}

		return []byte("OK")
	})
	if err != nil {
		log.Fatalln("quics: Error while receiving message from client. ", err)
	}

	// [REGISTER] CLIENT: sync root directory (remote to local)
	err = proto.RecvMessageWithResponseHandleFunc(REGISTER_SYNC_ROOTDIR, func(conn *qp.Connection, msgType string, data []byte) []byte {
		// decode request data
		var request types.SyncRootDirRequest
		if err := request.Decode(data); err != nil {
			log.Println("quics: (SyncRootDirRequest) Error while decoding request data")
			return []byte("FAIL")
		}

		//err := RegistrationHandler.RegistrationService.SyncRootDir(request)
		//if err != nil {
		//	log.Println("[QUICS] (SyncRootDirRequest) Error while creating root directory: ", err)
		//	return []byte("FAIL")
		//}

		return []byte("OK")
	})
	if err != nil {
		log.Printf("[QUICS-PROTOCOL] Error while receiving message from client: %s\n", err)
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

	cert, err := qp.GetCertificate("", "")
	if err != nil {
		log.Println("quics-protocol: ", err)
		return
	}

	tlsConfig := &tls.Config{
		Certificates: cert,
		NextProtos:   []string{"quic-s"},
	}

	// TODO: conn instance를 리스트화 시켜서 이후 브로드 캐스트 시 사용할 수 있도록 해야 한다.
	conn := func(conn *qp.Connection, msgType string, data []byte) {
		// decode request data
		var request types.RegisterClientRequest
		if err := request.Decode(data); err != nil {
			log.Println("[QUICS] (RegisterClientRequest) Error while decoding request data: ", err)
			err := conn.Close()
			if err != nil {
				log.Println("[QUICS-PROTOCOL] Error while closing the connection with client: ", err)
			}
			return
		}

		err := RegistrationHandler.RegistrationService.CreateNewClient(request.Uuid, request.ClientPassword, conn.Conn.RemoteAddr().String())
		if err != nil {
			log.Println("[QUICS] (RegisterClientRequest) Error while creating new client: ", err)
			err := conn.Close()
			if err != nil {
				log.Println("[QUICS-PROTOCOL] Error while closing the connection with client: ", err)
			}
			return
		}
	}

	// connect handler
	connectProtocolHandler(proto)

	go func() {
		// protocol
		err = proto.ListenWithMessage(UDPAddr, tlsConfig, conn)
		if err != nil {
			log.Fatalf("Error while listening protocol: %s", err)
		}
	}()

	fmt.Println("QUIC-S protocol listening successfully.")
}
