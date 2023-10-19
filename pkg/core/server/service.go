package server

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/quic-s/quics/pkg/config"
	"github.com/quic-s/quics/pkg/core/registration"
	"github.com/quic-s/quics/pkg/core/sync"
	"github.com/quic-s/quics/pkg/fs"
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

func NewService(repo *badger.Badger, serverRepository Repository) (Service, error) {
	// get env variables (server password, port)
	password := config.GetViperEnvVariables("PASSWORD")
	port, err := strconv.Atoi(config.GetViperEnvVariables("QUICS_PORT"))
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	pool := connection.NewnPool()

	registrationRepository := repo.NewRegistrationRepository()
	registrationNetworkAdapter := qp.NewRegistrationAdapter(pool)
	registrationService := registration.NewService(password, registrationRepository, registrationNetworkAdapter)
	registrationHandler := qp.NewRegistrationHandler(registrationService)

	historyRepository := repo.NewHistoryRepository()
	syncNetworkAdapter := qp.NewSyncAdapter(pool)
	syncDirAdapter := fs.NewSyncDir(utils.GetQuicsSyncDirPath())

	syncRepository := repo.NewSyncRepository()
	syncService := sync.NewService(registrationRepository, historyRepository, syncRepository, syncNetworkAdapter, syncDirAdapter)
	syncHandler := qp.NewSyncHandler(syncService)

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
	proto.RecvTransactionHandleFunc(types.CHOOSEONE, syncHandler.ChooseOne)
	proto.RecvTransactionHandleFunc(types.RESCAN, syncHandler.Rescan)

	return &ServerService{
		port:     port,
		password: password,
		repo:     repo,

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
	err := ss.Proto.Start()
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
			for _, root := range client.Root {
				fmt.Printf("*   UUID: %s   |   ID: %d   |   IP: %s   |   Root Directoreis: %s   *\n", client.UUID, client.Id, client.Ip, root)
			}
		}

		return nil
	}

	if id != "" {
		client, err := ss.serverRepository.GetClientByUUID(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		for _, root := range client.Root {
			fmt.Printf("*   UUID: %s   |   ID: %d   |   IP: %s   |   Root Directory: %s   *\n", client.UUID, client.Id, client.Ip, root.AfterPath)
		}

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
			for _, UUID := range dir.UUIDs {
				fmt.Printf("*   Root Directory: %s   |   Owner: %s   |   Password: %s   |   UUID: %s   *\n", dir.AfterPath, dir.Owner, dir.Password, UUID)
			}
		}

		return nil
	}

	if id != "" {
		dir, err := ss.serverRepository.GetRootDirectoryByPath(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		for _, UUID := range dir.UUIDs {
			fmt.Printf("*   Root Directory: %s   |   Owner: %s   |   Password: %s   |   UUID: %s   *\n", dir.AfterPath, dir.Owner, dir.Password, UUID)
		}

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
			fmt.Printf("*   File: %s   |   Root Directory: %s   |   LatestHash: %s   |   LatestSyncTimestamp: %d   |   ContentsExisted: %t   |   Metadata: %s   *\n", file.AfterPath, file.RootDirKey, file.LatestHash, file.LatestSyncTimestamp, file.ContentsExisted, file.Metadata.ModTime)
		}

		return nil
	}

	if id != "" {
		file, err := ss.serverRepository.GetFileByAfterPath(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		fmt.Printf("*   File: %s   |   Root Directory: %s   |   LatestHash: %s   |   LatestSyncTimestamp: %d   |   ContentsExisted: %t   |   Metadata: %s   *\n", file.AfterPath, file.RootDirKey, file.LatestHash, file.LatestSyncTimestamp, file.ContentsExisted, file.Metadata.ModTime)

		return nil
	}

	return nil
}

func (ss *ServerService) DisconnectClient(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                    Disconnect Client                       ")
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

func (ss *ServerService) DisconnectDir(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                  Disconnect Directory                      ")
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

func (ss *ServerService) DisconnectFile(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                     Disconnect File                        ")
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
