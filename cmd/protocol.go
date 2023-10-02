package main

import (
	"crypto/tls"
	"fmt"
	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/types"
	"github.com/quic-s/quics/pkg/utils"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// startQuicsProtocol starts quics protocol server
func startQuicsProtocol() {

	// initialize protocol server
	proto, err := qp.New(qp.LOG_LEVEL_INFO)
	if err != nil {
		log.Println("quics: ", err)
		return
	}

	// initialize server port
	portStr := config.GetViperEnvVariables("QUICS_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Println("quics: ", err)
		return
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
		return
	}

	// initialize tls config for connection with quics protocol
	tlsConfig := &tls.Config{
		Certificates: cert,
		NextProtos:   []string{"quic-s"},
	}

	// connect handler to quics protocol
	connectProtocolHandler(proto)

	go func() {
		// listen quics protocol with client
		err = proto.Listen(UDPAddr, tlsConfig, func(conn *qp.Connection) {
			fmt.Println("Successfully connected with ", conn.Conn.RemoteAddr().String())
		})
		if err != nil {
			log.Println("quics: ", err)
			return
		}
	}()

	fmt.Println("QUIC-S protocol listening successfully.")
}

// connectProtocolHandler connects handler to quics protocol
func connectProtocolHandler(proto *qp.QP) {
	var err error

	// register client
	// 1. (client) Open transaction
	// 2. (client) Send request data for registering client
	// 3. (server) Receive request data
	// 4. (server) Create new client to database
	// TODO: 5. (server) Send response data for registering client
	err = proto.RecvTransactionHandleFunc(types.REGISTERCLIENT, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
		log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

		data, err := stream.RecvBMessage()
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		var request types.ClientRegisterReq
		if err := request.Decode(data); err != nil {
			log.Println("quics: ", err)
			err := conn.Close()
			if err != nil {
				log.Println("quics: ", err)
			}
			return
		}

		// create new client to database
		err = RegistrationHandler.RegistrationService.CreateNewClient(request, Password, conn.Conn.RemoteAddr().String())
		if err != nil {
			log.Println("quics: ", err)
			err := conn.Close()
			if err != nil {
				log.Println("quics: ", err)
			}
			return
		}

		// TODO: is it necessary to send response data?
	})
	if err != nil {
		log.Println("quics: ", err)
	}

	// register root directory
	// 1. (client) Open transaction
	// 2. (client) Send request data for registering root directory
	// 3. (server) Receive request data
	// 4. (server) Register root directory of client to database
	// TODO: 5. (server) Send response data for registering root directory
	err = proto.RecvTransactionHandleFunc(types.REGISTERROOTDIR, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
		log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

		data, err := stream.RecvBMessage()
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		var request types.RegisterRootDirReq
		if err = request.Decode(data); err != nil {
			log.Println("quics: ", err)
			return
		}

		// Register root directory of client to database
		err = RegistrationHandler.RegistrationService.RegisterRootDir(request)
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		// TODO: is it necessary to send response data?
	})
	if err != nil {
		log.Println("quics: ", err)
	}

	// sync root directory
	// 1. (client) Open transaction
	// 2. (client) Send request data for syncing root directory
	// 3. (server) Receive request data
	// 4. (server) Sync root directory of client to database
	// TODO: 5. (server) Send response data for syncing root directory
	err = proto.RecvTransactionHandleFunc(types.SYNCROOTDIR, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
		log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

		data, err := stream.RecvBMessage()
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		var request types.SyncRootDirReq
		if err := request.Decode(data); err != nil {
			log.Println("quics: ", err)
			return
		}

		// get root directory path of requested data
		err = RegistrationHandler.RegistrationService.SyncRootDir(request)
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		// TODO: is it necessary to send response data?
	})
	if err != nil {
		log.Println("quics: ", err)
	}

	// get root directory list
	// 1. (client) Open transaction
	// 2. (client) Send request for getting root directory list
	// 3. (server) Receive request data
	// 4. (server) Get root directory list from database
	// 5. (server) Send response data for getting root directory list
	err = proto.RecvTransactionHandleFunc(types.GETROOTDIRS, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
		log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

		rootDirs := RegistrationHandler.RegistrationService.GetRootDirList()

		var rootDirPaths []byte
		for _, rootDir := range rootDirs {
			rootDirPath := rootDir.AfterPath
			rootDirPaths = append(rootDirPaths, []byte(rootDirPath)...)
		}

		err = stream.SendBMessage(rootDirPaths)
		if err != nil {
			log.Println("quics: ", err)
			return
		}
	})
	if err != nil {
		log.Println("quics: ", err)
	}

	// please sync transaction
	// 1. (client) Open transaction
	// 2. (client) PleaseFileMetaReq for getting a file metadata
	// 3. (server) Find and return certain file metadata
	// 4. (server) PleaseFileMetaRes for returning a file metadata
	// 5. (client) PleaseSyncReq if file update is available
	// 6. (server) Update the history with file metadata and set flag 'ContentsExisted' = false
	// 7. (server) PleaseSyncRes
	// 8. (client) PleaseTakeReq for sync a file
	// 9. (server) Get file contents and set flag 'ContentsExisted' = true
	// 10. (server) PleaseTakeRes
	// 11. (server) Go to the MustSync transaction
	err = proto.RecvTransactionHandleFunc(types.PLEASESYNC, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) {
		log.Println("quics: message received ", conn.Conn.RemoteAddr().String())

		data, err := stream.RecvBMessage()
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		var pleaseFileMetaReq types.PleaseFileMetaReq
		if err := pleaseFileMetaReq.Decode(data); err != nil {
			log.Println("quics: ", err)
			return
		}

		// TODO: find and return certain file metadata

		var pleaseSyncReq types.PleaseSyncReq
		if err := pleaseSyncReq.Decode(data); err != nil {
			log.Println("quics: ", err)
			return
		}

		// FIXME: change the condition from whether the file is exist to whether the request data is empty or not full
		// TODO: should think file with directories

		requestPaths := strings.Split(pleaseSyncReq.AfterPath, "/")
		rootDirName := requestPaths[1]
		rootDirPath := utils.GetQuicsRootDirPath(rootDirName)
		rootDir := RegistrationHandler.RegistrationService.GetRootDirByPath(rootDirPath)

		historyDirPath := utils.GetQuicsHistoryPathByRootDir(rootDirName)
		historyFilePath := historyDirPath + strconv.FormatUint(pleaseSyncReq.LastUpdateTimestamp, 10) + "_" + paths[2]

		fileSavedPath := utils.GetQuicsLatestPathByRootDir(rootDirName) // including latest directory

		updatedFile := types.File{
			BeforePath:          fileSavedPath,
			AfterPath:           pleaseSyncReq.AfterPath,
			RootDir:             rootDir,
			LatestHash:          pleaseSyncReq.LastUpdateHash,
			LatestSyncTimestamp: pleaseSyncReq.LastUpdateTimestamp,
		}

		dirPaths := strings.Split(updatedFile.BeforePath, "/")
		dirPaths = dirPaths[:len(dirPaths)-1]
		dirPath := strings.Join(dirPaths, "/")

		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			log.Println("quics: ", err)
			return
		}

		//inputFile, err := os.Create(updatedFile.Path)
		//if err != nil {
		//	log.Println("quis: ", err)
		//	return
		//}
		//defer inputFile.Close()
		//
		//n, err := io.Copy(inputFile, fileReader)
		//if err != nil {
		//	fmt.Println("quics: Error while copying file: ", err)
		//	return
		//}
		//if n != fileInfo.Size {
		//	fmt.Println("quics: read only ", n, "bytes")
		//	return
		//}
		//
		//if err := os.MkdirAll(historyDirPath, os.ModePerm); err != nil {
		//	log.Println("quics: ", err)
		//	return
		//}
		//
		//historyFile, err := os.Create(historyFilePath)
		//if err != nil {
		//	log.Fatalln("quis: Error while creating file: ", err)
		//	return
		//}
		//defer inputFile.Close()
		//
		//n, err = io.Copy(historyFile, fileReader)
		//if err != nil {
		//	fmt.Println("quics: Error while copying file: ", err)
		//	return
		//}
		//if n != fileInfo.Size {
		//	fmt.Println("quics: read only ", n, "bytes")
		//	return
		//}
		//
		//err = SyncHandler.SyncService.SaveFileFromPleaseSync(path, updatedFile)
		//if err != nil {
		//	log.Println("quics: Error while saving file: ", err)
		//}

		// open must sync transaction
		openMustSyncTransaction(conn)

	})
	if err != nil {
		log.Println("quics: ", err)
	}
}

func openMustSyncTransaction(conn *qp.Connection) {
	var err error

	// must sync transaction
	// 1. (server) Open transaction
	// 2. (server) MustSyncReq with file metadata to all registered clients without where the file come from
	// 3. (client) MustSyncRes if file update is available
	// 3-1. (server) If all request data are exist, then go to step 4
	// 3-2. (server) If not, then this transaction should be closed
	// 4. (server) GiveYouReq for giving file contents
	// 5. (client) GiveYouRes
	err = conn.OpenTransaction(types.MUSTSYNC, func(stream *qp.Stream, transactionName string, transactionId []byte) error {
		// TODO: implement must sync transaction process

		return nil
	})
	if err != nil {
		log.Println("quics: ", err)
	}
}
