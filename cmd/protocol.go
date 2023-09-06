package main

import (
	"crypto/tls"
	"fmt"
	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics-protocol/pkg/utils/fileinfo"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/registration"
	"github.com/quic-s/quics/pkg/types"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

const (
	RegisterRootdir     string = "LOCALROOT"
	RegisterRootdirList string = "SHOWREMOTELIST"
	RegisterSyncRootdir string = "REMOTEROOT"

	FilesSyncPleasesync string = "PLEASESYNC"
	FilesSyncPleasefile string = "PLEASEFILE"
	FilesSyncMustuync   string = "MUSTSYNC"
	FilesSyncTwooptions string = "TWOOPTIONS"
)

func connectProtocolHandler(proto *qp.QP) {

	// [REGISTER] CLIENT: register root directory (local to remote)
	err := proto.RecvMessageWithResponseHandleFunc(RegisterRootdir, func(conn *qp.Connection, msgType string, data []byte) []byte {
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
	err = proto.RecvMessageWithResponseHandleFunc(RegisterSyncRootdir, func(conn *qp.Connection, msgType string, data []byte) []byte {
		var request types.SyncRootDirRequest
		if err := request.Decode(data); err != nil {
			log.Println("quics: (SyncRootDirRequest) Error while decoding request data")
			return []byte("FAIL")
		}

		// get root directory path of requested data
		err := RegistrationHandler.RegistrationService.SyncRootDir(request)
		if err != nil {
			log.Println("quics: (SyncRootDirRequest) Error while creating root directory: ", err)
			return []byte("FAIL")
		}

		return []byte("OK")
	})
	if err != nil {
		log.Printf("quics: (SyncRootDirRequest) Error while receiving message from client: %s\n", err)
	}

	// [REGISTER] get root directory list
	err = proto.RecvMessageWithResponseHandleFunc(RegisterRootdirList, func(conn *qp.Connection, msgType string, data []byte) []byte {

		// get all root direcotry list
		rootDirs := RegistrationHandler.RegistrationService.GetRootDirList()

		var rootDirPaths []byte
		for _, rootDir := range rootDirs {
			rootDirPath := rootDir.Path[len(config.GetSyncDirPath()):]
			rootDirPaths = append(rootDirPaths, []byte(rootDirPath)...)
		}

		return rootDirPaths
	})
	if err != nil {
		log.Printf("quics: Error while receiving message from client: %s\n", err)
	}

	// [SYNC] listen PleaseSync message
	err = proto.RecvFileMessageHandleFunc(FilesSyncPleasesync, func(conn *qp.Connection, fileMsgType string, msgData []byte, fileInfo *fileinfo.FileInfo, fileReader io.Reader) {
		var request types.PleaseSync
		if err := request.Decode(msgData); err != nil {
			log.Println("quics: (PleaseSync) Error while decoding request data")
			return
		}

		// verify if conflict is occurred
		path := config.GetDirPath() + request.AfterPath
		file, isConflict := SyncHandler.SyncService.CheckIsOccurredConflict(path, request)

		if isConflict == 1 {
			// conflict
			rootDirPath := registration.ExtractRelateiveRootDirPath(file.RootDir.Path)
			conflictPath := config.GetSyncRootDirPath(rootDirPath)
			filePath := filepath.Join(conflictPath) + fileInfo.Name + "_" + strconv.FormatUint(request.LastUpdatedTimestamp, 10)
			inputFile, err := os.Create(filePath)
			if err != nil {
				log.Fatalln("quis: Error while creating file: ", err)
				return
			}
			defer inputFile.Close()

			n, err := io.Copy(inputFile, fileReader)
			if err != nil {
				fmt.Println("quics: Error while copying file: ", err)
				return
			}
			if n != fileInfo.Size {
				fmt.Println("quics: read only ", n, "bytes")
				return
			}
			log.Println("quics: conflict file saved")

			for {
				twoOptions := types.TwoOptions{
					ServerSideHash:          file.LatestHash,
					ServerSideSyncTimestamp: file.LatestSyncTimestamp,
					ClientSideHash:          request.LastUpdateHash,
					ClientSideTimestamp:     request.LastUpdatedTimestamp,
				}

				data, err := twoOptions.Encode()
				if err != nil {
					fmt.Println("quics: (PleaseSync) Error while encoding data: ", err)
				}

				response, err := conn.SendMessageWithResponse(FilesSyncTwooptions, data)
				if err != nil {
					fmt.Println("quics: (PleaseSync) Error while sending message to client: ", err)
					continue
				}

				if string(response) == "OK" {
					break
				} else {
					continue
				}
			}

		} else {
			//paths := strings.Split(request.AfterPath, "/")
			//rootDirPath := config.GetSyncDirPath() + paths[0]
			//rootDir := RegistrationHandler.RegistrationService.GetRootDirByPath(rootDirPath)
			//
			//// not conflict
			//updatedFile := types.File{
			//	Path:                config.GetSyncDirPath() + request.AfterPath,
			//	RootDir:             rootDir,
			//	LatestHash:          request.LastUpdateHash,
			//	LatestSyncTimestamp: request.LastUpdatedTimestamp,
			//}
			//
			//SyncHandler.SyncService.
		}
	})
	if err != nil {
		log.Printf("quics: Error while receiving message from client: %s\n", err)
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

	conn := func(conn *qp.Connection, msgType string, data []byte) {

		// connect
		Conns = append(Conns, conn)

		var request types.RegisterClientRequest
		if err := request.Decode(data); err != nil {
			log.Println("quics: (RegisterClientRequest) Error while decoding request data: ", err)
			err := conn.Close()
			if err != nil {
				log.Println("quics: Error while closing the connection with client: ", err)
			}
			return
		}

		password := ServerHandler.ServerService.GetPassword()

		err := RegistrationHandler.RegistrationService.CreateNewClient(request, password, conn.Conn.RemoteAddr().String())
		if err != nil {
			log.Println("quics: (RegisterClientRequest) Error while creating new client: ", err)
			err := conn.Close()
			if err != nil {
				log.Println("quics: Error while closing the connection with client: ", err)
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
