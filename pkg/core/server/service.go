package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/quic-s/quics/pkg/config"
	"github.com/quic-s/quics/pkg/core/history"
	"github.com/quic-s/quics/pkg/core/registration"
	"github.com/quic-s/quics/pkg/core/sharing"
	"github.com/quic-s/quics/pkg/core/sync"
	"github.com/quic-s/quics/pkg/network/qp"
	"github.com/quic-s/quics/pkg/network/qp/connection"
	"github.com/quic-s/quics/pkg/repository/badger"
	"github.com/quic-s/quics/pkg/types"
	"github.com/quic-s/quics/pkg/utils"
)

type ServerService struct {
	port     int
	password string
	repo     *badger.Badger
	Proto    *qp.Protocol

	syncService sync.Service

	serverRepository Repository
}

func NewService(repo *badger.Badger, serverRepository Repository, syncDirAdapter sync.SyncDirAdapter) (Service, error) {
	password := ""

	server, err := repo.NewServerRepository().GetPassword()
	if err != nil {
		password = config.GetViperEnvVariables("PASSWORD")
	} else {
		password = server.Password
	}

	// get env variables (server password, port)
	port, err := strconv.Atoi(config.GetViperEnvVariables("QUICS_PORT"))
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	pool := connection.NewnPool()

	registrationRepository := repo.NewRegistrationRepository()
	historyRepository := repo.NewHistoryRepository()
	syncRepository := repo.NewSyncRepository()
	sharingRepository := repo.NewSharingRepository()

	registrationNetworkAdapter := qp.NewRegistrationAdapter(pool)
	syncNetworkAdapter := qp.NewSyncAdapter(pool)

	registrationService := registration.NewService(password, registrationRepository, registrationNetworkAdapter)
	historyService := history.NewService(historyRepository)
	syncService := sync.NewService(registrationRepository, historyRepository, syncRepository, syncNetworkAdapter, syncDirAdapter)
	sharingService := sharing.NewService(historyRepository, syncRepository, sharingRepository, syncDirAdapter)

	registrationHandler := qp.NewRegistrationHandler(registrationService)
	syncHandler := qp.NewSyncHandler(syncService)
	historyHandler := qp.NewHistoryHandler(historyService, sharingService)
	sharingHandler := qp.NewSharingHandler(sharingService)

	proto, err := qp.New("0.0.0.0", port, pool)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	proto.RecvTransactionHandleFunc(types.REGISTERCLIENT, registrationHandler.RegisterClient)
	proto.RecvTransactionHandleFunc(types.REGISTERROOTDIR, syncHandler.RegisterRootDir)
	proto.RecvTransactionHandleFunc(types.SYNCROOTDIR, syncHandler.SyncRootDir)
	proto.RecvTransactionHandleFunc(types.GETROOTDIRS, syncHandler.GetRemoteDirs)
	proto.RecvTransactionHandleFunc(types.PLEASESYNC, syncHandler.PleaseSync)
	proto.RecvTransactionHandleFunc(types.CONFLICTLIST, syncHandler.AskConflictList)
	proto.RecvTransactionHandleFunc(types.CONFLICTDOWNLOAD, syncHandler.ConflictDownload)
	proto.RecvTransactionHandleFunc(types.CHOOSEONE, syncHandler.ChooseOne)
	proto.RecvTransactionHandleFunc(types.RESCAN, syncHandler.Rescan)
	proto.RecvTransactionHandleFunc(types.HISTORYSHOW, historyHandler.ShowHistory)
	proto.RecvTransactionHandleFunc(types.ROLLBACK, syncHandler.RollbackFileByHistory)
	proto.RecvTransactionHandleFunc(types.HISTORYDOWNLOAD, syncHandler.DownloadHistory)
	proto.RecvTransactionHandleFunc(types.STARTSHARING, sharingHandler.StartSharing)
	proto.RecvTransactionHandleFunc(types.STOPSHARING, sharingHandler.StopSharing)

	return &ServerService{
		port:     port,
		password: password,
		repo:     repo,
		Proto:    proto,

		syncService:      syncService,
		serverRepository: serverRepository,
	}, nil
}

// StopServer stop quic-s server
func (ss *ServerService) StopServer() error {
	fmt.Println("************************************************************")
	fmt.Println("                           Stop                             ")
	fmt.Println("************************************************************")

	go func() {
		err := ss.repo.Close()
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		log.Println("quics: Closed")
		os.Exit(0)
	}()

	return nil
}

// ListenProtocol is executed when server starts
func (ss *ServerService) ListenProtocol() error {
	fmt.Println("************************************************************")
	fmt.Println("                     Listen Protocol                        ")
	fmt.Println("************************************************************")

	// start quics protocol server
	ss.syncService.BackgroundFullScan(300)
	errChan := make(chan error)
	go func() {
		go func() {
			time.Sleep(3 * time.Second)
			errChan <- nil
		}()
		err := ss.Proto.Start()
		if err != nil {
			log.Println("quics: ", err)
			errChan <- err
		}
	}()

	err := <-errChan
	if err != nil {
		return err
	}
	return nil
}

func (ss *ServerService) SetPassword(request *types.Server) error {
	fmt.Println("************************************************************")
	fmt.Println("                       Set Password                         ")
	fmt.Println("************************************************************")

	err := ss.serverRepository.UpdatePassword(request)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

func (ss *ServerService) ResetPassword() error {
	fmt.Println("************************************************************")
	fmt.Println("                      Reset Password                        ")
	fmt.Println("************************************************************")

	err := ss.serverRepository.DeletePassword()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

func (ss *ServerService) Ping(request *types.Ping) (*types.Ping, error) {
	client, err := ss.serverRepository.GetClientByUUID(request.UUID)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	if request.UUID != client.UUID {
		return nil, errors.New("quics: (Ping) UUID is not correct")
	}

	return &types.Ping{
		UUID: request.UUID,
	}, nil
}

func (ss *ServerService) ShowClientLogs(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                        Client Log                          ")
	fmt.Println("************************************************************")

	if all != "" {
		clients, err := ss.serverRepository.GetAllClients()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		for _, client := range clients {
			fmt.Sprintf("\n")
			fmt.Sprintf("\tUUID           : %s\n", client.UUID)
			fmt.Sprintf("\tID             : %d\n", client.Id)
			fmt.Sprintf("\tIP             : %s\n", client.Ip)

			for _, root := range client.Root {
				fmt.Sprintf("\tRoot Directory : %s\n", root.AfterPath)
			}
			fmt.Sprintf("\n")
		}

		return nil
	}

	if id != "" {
		client, err := ss.serverRepository.GetClientByUUID(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		fmt.Sprintf("\n")
		fmt.Sprintf("\tUUID           : %s\n", client.UUID)
		fmt.Sprintf("\tID             : %d\n", client.Id)
		fmt.Sprintf("\tIP             : %s\n", client.Ip)

		for _, root := range client.Root {
			fmt.Sprintf("\tRoot Directory : %s\n", root.AfterPath)
		}

		fmt.Sprintf("\n")

		return nil
	}

	return nil
}

func (ss *ServerService) ShowDirLogs(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                       Directory Log                        ")
	fmt.Println("************************************************************")

	if all != "" {
		dirs, err := ss.serverRepository.GetAllRootDirectories()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		for _, dir := range dirs {
			fmt.Sprintf("\n")
			fmt.Sprintf("\tRoot Directory : %s\n", dir.AfterPath)
			fmt.Sprintf("\tOwner          : %s\n", dir.Owner)
			fmt.Sprintf("\tPassword       : %s\n", dir.Password)

			for _, UUID := range dir.UUIDs {
				fmt.Sprintf("\tClient         : %s\n", UUID)
			}

			fmt.Sprintf("\n")
		}

		return nil
	}

	if id != "" {
		dir, err := ss.serverRepository.GetRootDirectoryByPath(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		fmt.Sprintf("\n")
		fmt.Sprintf("\tRoot Directory : %s\n", dir.AfterPath)
		fmt.Sprintf("\tOwner          : %s\n", dir.Owner)
		fmt.Sprintf("\tPassword       : %s\n", dir.Password)

		for _, UUID := range dir.UUIDs {
			fmt.Sprintf("\tClient         : %s\n", UUID)
		}

		fmt.Sprintf("\n")

		return nil
	}

	return nil
}

func (ss *ServerService) ShowFileLogs(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                         File Log                           ")
	fmt.Println("************************************************************")

	if all != "" {
		files, err := ss.serverRepository.GetAllFiles()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		for _, file := range files {
			fmt.Sprintf("\n")
			fmt.Sprintf("\tFile                  : %s\n", file.AfterPath)
			fmt.Sprintf("\tRoot Directory        : %d\n", file.RootDirKey)
			fmt.Sprintf("\tLatest Hash           : %s\n", file.LatestHash)
			fmt.Sprintf("\tLatest Sync Timestamp : %s\n", file.LatestSyncTimestamp)
			fmt.Sprintf("\tLatest Edit Client    : %s\n", file.LatestEditClient)
			fmt.Sprintf("\tContents Existed      : %s\n", file.ContentsExisted)
			fmt.Sprintf("\tMetadata              : %s\n", file.Metadata.ModTime)
			fmt.Sprintf("\n")
		}

		return nil
	}

	if id != "" {
		file, err := ss.serverRepository.GetFileByAfterPath(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		fmt.Sprintf("\n")
		fmt.Sprintf("\tFile                  : %s\n", file.AfterPath)
		fmt.Sprintf("\tRoot Directory        : %d\n", file.RootDirKey)
		fmt.Sprintf("\tLatest Hash           : %s\n", file.LatestHash)
		fmt.Sprintf("\tLatest Sync Timestamp : %s\n", file.LatestSyncTimestamp)
		fmt.Sprintf("\tLatest Edit Client    : %s\n", file.LatestEditClient)
		fmt.Sprintf("\tContents Existed      : %s\n", file.ContentsExisted)
		fmt.Sprintf("\tMetadata              : %s\n", file.Metadata.ModTime)
		fmt.Sprintf("\n")

		return nil
	}

	return nil
}

func (ss *ServerService) ShowHistoryLogs(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                       History Log                          ")
	fmt.Println("************************************************************")

	if all != "" {
		histories, err := ss.serverRepository.GetAllHistories()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		for _, history := range histories {
			fmt.Sprintf("\n")
			fmt.Sprintf("\tPath      : %s\n", history.BeforePath+history.AfterPath)
			fmt.Sprintf("\tDate      : %d\n", history.Date)
			fmt.Sprintf("\tClient    : %s\n", history.UUID)
			fmt.Sprintf("\tTimestamp : %s\n", history.Timestamp)
			fmt.Sprintf("\tHash      : %s\n", history.Hash)
			fmt.Sprintf("\n")
		}

		return nil
	}

	if id != "" {
		history, err := ss.serverRepository.GetHistoryByAfterPath(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		fmt.Sprintf("\n")
		fmt.Sprintf("\tPath      : %s\n", history.BeforePath+history.AfterPath)
		fmt.Sprintf("\tDate      : %d\n", history.Date)
		fmt.Sprintf("\tClient    : %s\n", history.UUID)
		fmt.Sprintf("\tTimestamp : %s\n", history.Timestamp)
		fmt.Sprintf("\tHash      : %s\n", history.Hash)
		fmt.Sprintf("\n")

		return nil
	}

	return nil
}

func (ss *ServerService) RemoveClient(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                       Remove Client                        ")
	fmt.Println("************************************************************")

	if all != "" {
		err := ss.serverRepository.DeleteAllClients()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	if id != "" {
		err := ss.serverRepository.DeleteClientByUUID(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	return nil
}

func (ss *ServerService) RemoveDir(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                      Remove Directory                      ")
	fmt.Println("************************************************************")

	if all != "" {
		err := ss.serverRepository.DeleteAllRootDirectories()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	if id != "" {
		err := ss.serverRepository.DeleteRootDirectoryByAfterPath(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	return nil
}

func (ss *ServerService) RemoveFile(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                        Remove File                         ")
	fmt.Println("************************************************************")

	if all != "" {
		err := ss.serverRepository.DeleteAllFiles()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	if id != "" {
		err := ss.serverRepository.DeleteFileByAfterPath(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	return nil
}

func (ss *ServerService) DownloadFile(path string, version string, target string) error {
	fmt.Println("************************************************************")
	fmt.Println("                      Download File                         ")
	fmt.Println("************************************************************")

	if strings.Contains(target, utils.GetQuicsDirPath()) {
		return errors.New("quics: target path should not be in .quics directory")
	}

	sourceFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// copy contents
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
