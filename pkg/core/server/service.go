package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
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
)

type ServerService struct {
	port     int
	password string
	repo     *badger.Badger
	Proto    *qp.Protocol

	syncService sync.Service

	syncDirAdapter   SyncDirAdapter
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
		log.Println("quics err: ", err)
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
		log.Println("quics err: ", err)
		return nil, err
	}

	proto.RecvTransactionHandleFunc(types.REGISTERCLIENT, registrationHandler.RegisterClient)
	proto.RecvTransactionHandleFunc(types.REGISTERROOTDIR, syncHandler.RegisterRootDir)
	proto.RecvTransactionHandleFunc(types.DISCONNECTROOTDIR, syncHandler.DisconnectRootDir)
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
		syncDirAdapter:   syncDirAdapter,
		serverRepository: serverRepository,
	}, nil
}

// StopServer stop quic-s server
func (ss *ServerService) StopServer() error {
	fmt.Println("************************************************************")
	fmt.Println("                           Stop                             ")
	fmt.Println("************************************************************")

	err := ss.repo.Close()
	if err != nil {
		return err
	}

	err = ss.Proto.Close()
	if err != nil {
		return err
	}
	log.Println("quics: Closed")

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
			log.Println("quics err: ", err)
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
	log.Println("quics: set password")

	err := config.WriteViperEnvVariables("PASSWORD", request.Password)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (ss *ServerService) ResetPassword() error {
	log.Println("quics: reset password")

	err := config.WriteViperEnvVariables("PASSWORD", config.DefaultPassword)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (ss *ServerService) Ping(request *types.Ping) (*types.Ping, error) {
	client, err := ss.serverRepository.GetClientByUUID(request.UUID)
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	if request.UUID != client.UUID {
		return nil, errors.New("quics: (Ping) UUID is not correct")
	}

	return &types.Ping{
		UUID: request.UUID,
	}, nil
}

func (ss *ServerService) ShowClient(uuid string) ([]types.Client, error) {
	log.Println("quics: show client logs (uudi: ", uuid, ")")

	if uuid == "" {
		clients, err := ss.serverRepository.GetAllClients()
		if err != nil {
			log.Println("quics err: ", err)
			return nil, err
		}
		return clients, nil
	}
	client, err := ss.serverRepository.GetClientByUUID(uuid)
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	return []types.Client{*client}, nil
}

func (ss *ServerService) ShowDir(afterPath string) ([]types.RootDirectory, error) {
	log.Println("quics: show dir logs (afterPath: ", afterPath, ")")

	if afterPath == "" {
		dirs, err := ss.serverRepository.GetAllRootDirectories()
		if err != nil {
			log.Println("quics err: ", err)
			return nil, err
		}

		return dirs, nil
	}

	dir, err := ss.serverRepository.GetRootDirectoryByPath(afterPath)
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}
	return []types.RootDirectory{*dir}, nil
}

func (ss *ServerService) ShowFile(afterPath string) ([]types.File, error) {
	log.Println("quics: show file logs (afterPath: ", afterPath, ")")

	if afterPath == "" {
		files, err := ss.serverRepository.GetAllFiles()
		if err != nil {
			log.Println("quics err: ", err)
			return nil, err
		}

		return files, nil
	}

	file, err := ss.serverRepository.GetFileByAfterPath(afterPath)
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}
	return []types.File{*file}, nil
}

func (ss *ServerService) ShowHistory(afterPath string) ([]types.FileHistory, error) {
	log.Println("quics: show history logs (afterPath: ", afterPath, ")")

	if afterPath == "" {
		histories, err := ss.serverRepository.GetAllHistories()
		if err != nil {
			log.Println("quics err: ", err)
			return nil, err
		}

		return histories, nil
	}

	history, err := ss.serverRepository.GetHistoryByAfterPath(afterPath)
	if err != nil {
		log.Println("quics err: ", err)
		return nil, err
	}

	fmt.Printf("*   Path: %s   |   Date: %s   |   UUID: %s   |   Timestamp: %d   |   Hash: %s   |*\n", history.BeforePath+history.AfterPath, history.Date, history.UUID, history.Timestamp, history.Hash)

	return []types.FileHistory{*history}, nil
}

func (ss *ServerService) RemoveClient(uuid string) error {
	log.Println("quics: remove client (uuid: ", uuid, ")")

	if uuid == "" {
		err := ss.serverRepository.DeleteAllClients()
		if err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	}
	err := ss.serverRepository.DeleteClientByUUID(uuid)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (ss *ServerService) RemoveDir(afterPath string) error {
	log.Println("quics: remove dir (afterPath: ", afterPath, ")")

	if afterPath == "" {
		err := ss.serverRepository.DeleteAllRootDirectories()
		if err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	}

	err := ss.serverRepository.DeleteRootDirectoryByAfterPath(afterPath)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (ss *ServerService) RemoveFile(afterPath string) error {
	log.Println("quics: remove file (afterPath: ", afterPath, ")")

	if afterPath == "" {
		err := ss.serverRepository.DeleteAllFiles()
		if err != nil {
			log.Println("quics err: ", err)
			return err
		}

		return nil
	}
	err := ss.serverRepository.DeleteFileByAfterPath(afterPath)
	if err != nil {
		log.Println("quics err: ", err)
		return err
	}

	return nil
}

func (ss *ServerService) DownloadFile(afterPath string, timestamp uint64) (*types.FileMetadata, io.Reader, error) {
	log.Println("quics: download file (afterPath: ", afterPath, ")")

	return ss.syncDirAdapter.GetFileFromHistoryDir(afterPath, timestamp)
}
