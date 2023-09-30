package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics-protocol/pkg/utils/fileinfo"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/registration"
	"github.com/quic-s/quics/pkg/types"
)

const (
	RegisterRootdir     string = "LOCALROOT"
	RegisterRootdirList string = "SHOWREMOTELIST"
	RegisterSyncRootdir string = "REMOTEROOT"

	FileSyncRescan       string = "RESCAN"
	FilesSyncPleasesync  string = "PLEASESYNC"
	FilesSyncPleasefile  string = "PLEASEFILE"
	FilesSyncMustsync    string = "MUSTSYNC"
	FilesSyncTwooptions  string = "TWOOPTIONS"
	FilesSyncGiveyoufile string = "GIVEYOUFILE"
)

var Interval = 180 * time.Second
var requestRescan = make(chan bool)

func connectProtocolHandler(proto *qp.QP) {

	// [REGISTER] CLIENT: register root directory (local to remote)
	err := proto.RecvMessageWithResponseHandleFunc(RegisterRootdir, func(conn *qp.Connection, msgType string, data []byte) []byte {
		// decode request data

		var request types.RootDirRegisterReq
		if err := request.Decode(data); err != nil {
			log.Println("quics: Error while decoding request data")
			return []byte("FAIL")
		}

		err := RegistrationHandler.RegistrationService.RegisterRootDir(request)
		if err != nil {
			log.Println("quics: (RootDirRegisterReq) Error while creating root directory: ", err)
			return []byte("FAIL")
		}

		return []byte("OK")
	})
	if err != nil {
		log.Fatalln("quics: Error while receiving message from client. ", err)
	}

	// [REGISTER] CLIENT: sync root directory (remote to local)
	err = proto.RecvMessageWithResponseHandleFunc(RegisterSyncRootdir, func(conn *qp.Connection, msgType string, data []byte) []byte {
		var request types.SyncRootDirReq
		if err := request.Decode(data); err != nil {
			log.Println("quics: (SyncRootDirReq) Error while decoding request data")
			return []byte("FAIL")
		}

		// get root directory path of requested data
		err := RegistrationHandler.RegistrationService.SyncRootDir(request)
		if err != nil {
			log.Println("quics: (SyncRootDirReq) Error while creating root directory: ", err)
			return []byte("FAIL")
		}

		return []byte("OK")
	})
	if err != nil {
		log.Printf("quics: (SyncRootDirReq) Error while receiving message from client: %s\n", err)
	}

	// [REGISTER] get root directory list
	err = proto.RecvMessageWithResponseHandleFunc(RegisterRootdirList, func(conn *qp.Connection, msgType string, data []byte) []byte {
		// get all root direcotry list
		rootDirs := RegistrationHandler.RegistrationService.GetRootDirList()
		log.Println("quics: rootDirs: ", rootDirs)

		var rootDirPaths []byte
		for _, rootDir := range rootDirs {
			rootDirPath := rootDir.Path[len(config.GetSyncDirPath()):]
			rootDirPaths = append(rootDirPaths, []byte(rootDirPath)...)
		}

		log.Println("quics: rootDirPaths: ", rootDirPaths)

		return rootDirPaths
	})
	if err != nil {
		log.Printf("quics: Error while receiving message from client: %s\n", err)
	}

	// [SYNC] listen PleaseSync message
	err = proto.RecvFileMessageHandleFunc(FilesSyncPleasesync, func(conn *qp.Connection, fileMsgType string, msgData []byte, fileInfo *fileinfo.FileInfo, fileReader io.Reader) {
		// parse request data
		var request types.PleaseSyncReq
		if err := request.Decode(msgData); err != nil {
			log.Println("quics: (PleaseSync) Error while decoding request data")
			return
		}

		requestPaths := strings.Split(request.AfterPath, "/")

		path := config.GetSyncDirPath() + "/" + requestPaths[1] + "/latest/" + requestPaths[2]

		// check if file is exist
		isExistFile := SyncHandler.SyncService.IsExistFile(path)
		if isExistFile == 1 {

			log.Println("quics: file is exist")

			// if exsit, then verify if conflict is occurred
			file, isConflict := SyncHandler.SyncService.CheckIsOccurredConflict(path, request)

			if isConflict == 1 {
				// if conflict is occurred
				// make conflict directory and save conflict file

				log.Println("quics: conflict is occurred")

				rootDirPath := registration.ExtractRelateiveRootDirPath(file.RootDir.Path)
				conflictPath := config.GetSyncRootDirPath(rootDirPath) + "/conflict/"
				filePath := filepath.Join(conflictPath) + fileInfo.Name + "_" + strconv.FormatUint(request.LastUpdateTimestamp, 10)

				dirPaths := strings.Split(filePath, "/")
				dirPaths = dirPaths[:len(dirPaths)-1]
				dirPath := strings.Join(dirPaths, "/")

				if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
					log.Fatalf("Error creating directory: %s with %s", err, dirPath)
				}

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

				// to resolve conflict, send messgae for two options to client
				for {
					twoOptions := types.TwoOptions{
						ServerSideHash:          file.LatestHash,
						ServerSideSyncTimestamp: file.LatestSyncTimestamp,
						ClientSideHash:          request.LastUpdateHash,
						ClientSideTimestamp:     request.LastUpdateTimestamp,
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
				paths := strings.Split(request.AfterPath, "/")
				rootDirPath := config.GetSyncDirPath() + "/" + paths[1]
				rootDir := RegistrationHandler.RegistrationService.GetRootDirByPath(rootDirPath)
				historyDirPath := rootDirPath + "/history/"
				historyFilePath := historyDirPath + paths[2] + "_" + strconv.FormatUint(request.LastUpdateTimestamp, 10)

				// not conflict
				updatedFile := types.File{
					Path:                rootDirPath + "/latest/" + paths[2],
					RootDir:             rootDir,
					LatestHash:          request.LastUpdateHash,
					LatestSyncTimestamp: request.LastUpdateTimestamp,
				}

				inputFile, err := os.Create(updatedFile.Path)
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

				historyFile, err := os.Create(historyFilePath)
				if err != nil {
					log.Fatalln("quis: Error while creating file: ", err)
					return
				}
				defer historyFile.Close()

				n, err = io.Copy(historyFile, inputFile)
				if err != nil {
					fmt.Println("quics: Error while copying file: ", err)
					return
				}
				if n != fileInfo.Size {
					fmt.Println("quics: read only ", n, "bytes")
					return
				}

				err = SyncHandler.SyncService.SaveFileFromPleaseSync(path, updatedFile)
				if err != nil {
					log.Println("quics: Error while saving file: ", err)
				}

				// broadcast
				mustSyncReq := types.MustSyncReq{
					LatestHash:          updatedFile.LatestHash,
					LatestSyncTimestamp: updatedFile.LatestSyncTimestamp,
					BeforePath:          config.GetSyncDirPath(),
					AfterPath:           request.AfterPath,
				}

				encodedMustSyncReq, err := mustSyncReq.Encode()

				err = conn.SendMessage(FilesSyncMustsync, encodedMustSyncReq)
				if err != nil {
					log.Println("quics: Error while sending message for MustSync: ", err)
				}
			}
		} else {
			// not exist
			paths := strings.Split(request.AfterPath, "/")
			rootDirPath := config.GetSyncDirPath() + "/" + paths[1]
			rootDir := RegistrationHandler.RegistrationService.GetRootDirByPath(rootDirPath)
			historyDirPath := rootDirPath + "/history/"
			historyFilePath := historyDirPath + strconv.FormatUint(request.LastUpdateTimestamp, 10) + "_" + paths[2]

			updatedFile := types.File{
				Path:                rootDirPath + "/latest/" + paths[2],
				RootDir:             rootDir,
				LatestHash:          request.LastUpdateHash,
				LatestSyncTimestamp: request.LastUpdateTimestamp,
			}

			dirPaths := strings.Split(updatedFile.Path, "/")
			dirPaths = dirPaths[:len(dirPaths)-1]
			dirPath := strings.Join(dirPaths, "/")

			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				log.Fatalf("Error creating directory: %s with %s", err, dirPath)
			}

			if err := os.MkdirAll(historyDirPath, os.ModePerm); err != nil {
				log.Fatalf("Error creating directory: %s with %s", err, historyDirPath)
			}

			inputFile, err := os.Create(updatedFile.Path)
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

			historyFile, err := os.Create(historyFilePath)
			if err != nil {
				log.Fatalln("quis: Error while creating file: ", err)
				return
			}
			defer inputFile.Close()

			n, err = io.Copy(historyFile, fileReader)
			if err != nil {
				fmt.Println("quics: Error while copying file: ", err)
				return
			}
			if n != fileInfo.Size {
				fmt.Println("quics: read only ", n, "bytes")
				return
			}

			err = SyncHandler.SyncService.SaveFileFromPleaseSync(path, updatedFile)
			if err != nil {
				log.Println("quics: Error while saving file: ", err)
			}

			// broadcast
			mustSyncReq := types.MustSyncReq{
				LatestHash:          updatedFile.LatestHash,
				LatestSyncTimestamp: updatedFile.LatestSyncTimestamp,
				BeforePath:          config.GetSyncDirPath(),
				AfterPath:           request.AfterPath,
			}

			encodedMustSyncReq, err := mustSyncReq.Encode()

			err = conn.SendMessage(FilesSyncMustsync, encodedMustSyncReq)
			if err != nil {
				log.Println("quics: Error while sending message for MustSync: ", err)
			}
		}
	})
	if err != nil {
		log.Printf("quics: Error while receiving message from client: %s\n", err)
	}

	// [SYNC] RESCAN
	err = proto.RecvMessageHandleFunc(FileSyncRescan, func(conn *qp.Connection, msgType string, data []byte) {
		requestRescan <- true
	})
	if err != nil {
		log.Printf("quics: Error while receiving message from client: %s\n", err)
	}

	// [SYNC] PleaseFile
	err = proto.RecvMessageHandleFunc(FilesSyncPleasefile, func(conn *qp.Connection, msgType string, data []byte) {
		// decode request data
		var request types.PleaseFileReq
		if err := request.Decode(data); err != nil {
			log.Println("quics: Error while decoding request data")
			return
		}

		paths := strings.Split(request.AfterPath, "/")
		rootDirPath := config.GetSyncDirPath() + "/" + paths[1]
		rootDirFilePath := rootDirPath + "/latest/" + paths[2]
		historyDirPath := rootDirPath + "/history/"
		historyFilePath := historyDirPath + paths[2] + "_" + strconv.FormatUint(request.SyncTimestamp, 10)

		file := SyncHandler.SyncService.GetFileByPath(rootDirFilePath)

		giveYouFile := types.GiveYouFile{
			LatestHash:          file.LatestHash,
			LatestSyncTimestamp: file.LatestSyncTimestamp,
			BeforePath:          config.GetSyncDirPath(),
			AfterPath:           request.AfterPath,
		}
		encodedGiveYouFile, err := giveYouFile.Encode()
		if err != nil {
			log.Println("quics: Error while encoding GiveYouFile: ", err)
		}

		err = conn.SendFileMessage(FilesSyncGiveyoufile, encodedGiveYouFile, historyFilePath)
		if err != nil {
			log.Println("quics: Error while sending file to client: ", err)
			return
		}
	})
	if err != nil {
		log.Println("quics: Error while receiving message from client: ", err)
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

		var request types.ClientRegisterReq
		if err := request.Decode(data); err != nil {
			log.Println("quics: (ClientRegisterReq while decoding request data: ", err)
			err := conn.Close()
			if err != nil {
				log.Println("quics: Error while closing the connection with client: ", err)
			}
			return
		}

		password := ServerHandler.ServerService.GetPassword()

		err := RegistrationHandler.RegistrationService.CreateNewClient(request, password, conn.Conn.RemoteAddr().String())
		if err != nil {
			log.Println("quics: (ClientRegisterReq while creating new client: ", err)
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

	go func() {
		ticker := time.NewTicker(Interval)

		for {
			select {
			case <-ticker.C:
				CallMustSync()
				ticker.Stop()
				ticker = time.NewTicker(Interval)
			case <-requestRescan:
				CallMustSync()
				ticker.Stop()
				ticker = time.NewTicker(Interval)
			}
		}
	}()
}

func CallMustSync() {

	files := SyncHandler.SyncService.GetFiles()

	for _, file := range files {

		var result string
		if strings.HasPrefix(file.Path, config.GetSyncDirPath()) {
			result = strings.TrimPrefix(file.Path, config.GetSyncDirPath())
		}

		target := "/latest/"
		afterPath := strings.Replace(result, target, "/", -1)

		mustSyncReq := types.MustSyncReq{
			LatestHash:          file.LatestHash,
			LatestSyncTimestamp: file.LatestSyncTimestamp,
			BeforePath:          config.GetSyncDirPath(),
			AfterPath:           afterPath,
		}
		encodedMustSyncReq, err := mustSyncReq.Encode()
		if err != nil {
			log.Println("quics: Error while encoding MustSyncReq: ", err)
		}

		for _, conn := range Conns {
			err := conn.SendMessage(FilesSyncMustsync, encodedMustSyncReq)
			if err != nil {
				log.Println("quics: Error while sending message for MustSync: ", err)
			}
		}
	}
}
