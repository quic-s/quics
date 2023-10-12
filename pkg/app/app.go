package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/quic-s/quics/pkg/config"
	"github.com/quic-s/quics/pkg/core/registration"
	"github.com/quic-s/quics/pkg/network/qp"
	"github.com/quic-s/quics/pkg/network/qp/connection"
	"github.com/quic-s/quics/pkg/repository/badger"
	"github.com/quic-s/quics/pkg/types"
)

type App struct {
	repo     *badger.Badger
	Proto    *qp.Protocol
	SigCh    chan os.Signal
	Password string
}

// Initialize initialize program
func New() (*App, error) {
	// define system call actions
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// get env variables
	password := config.GetViperEnvVariables("PASSWORD")
	// ip := config.GetViperEnvVariables("IP")
	port, err := strconv.Atoi(config.GetViperEnvVariables("QUICS_PORT"))
	if err != nil {
		log.Println("quics: quics port is not integer: ", err)
		return nil, err
	}

	repo, err := badger.NewBadgerRepository()
	if err != nil {
		log.Println("quics: Error while connecting to the database: ", err)
	}

	pool := connection.NewnPool()

	proto, err := qp.New("0.0.0.0", port, pool)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	registrationRepository := repo.NewRegistrationRepository()
	registrationNetworkAdapter := qp.NewRegistrationAdapter(pool)
	registrationService := registration.NewService(password, registrationRepository, registrationNetworkAdapter)
	registrationHandler := qp.NewRegistrationHandler(registrationService)

	proto.RecvTransactionHandleFunc(types.REGISTERCLIENT, registrationHandler.RegisterClient)
	proto.RecvTransactionHandleFunc(types.REGISTERROOTDIR, registrationHandler.RegisterRootDir)
	proto.RecvTransactionHandleFunc(types.GETROOTDIRS, registrationHandler.GetRemoteDirs)

	// historyRepository := repo.NewHistoryRepository()

	// syncRepository := repo.NewSyncRepository()
	// syncService := sync.NewService(registrationRepository, historyRepository, syncRepository)
	// syncHandler := qp.NewSyncHandler(syncService)

	// proto.RecvTransactionHandleFunc(types.SYNCROOTDIR, syncHandler.SyncRootDir)

	// historyRepository := repo.NewHistoryRepository()
	// historyService := history.NewHistoryService(historyRepository)
	// HistoryHandler := http3hdl.NewHistoryHandler(historyService)

	// metadataRepository := repo.NewMetadataRepository()
	// metadataService := metadata.NewMetadataService(metadataRepository)
	// MetadataHandler := http3hdl.NewMetadataHandler(metadataService)

	// registrationRepository := repo.NewRegistrationRepository()
	// registrationService := registration.NewService(registrationRepository)
	// RegistrationHandler := http3hdl.NewRegistrationHandler(registrationService)

	// serverRepository := repo.NewServerRepository()
	// serverService := server.NewService(serverRepository)
	// ServerHandler := http3hdl.NewServerHandler(serverService)

	// sharingRepository := repo.NewSharingRepository()
	// sharingService := sharing.NewService(sharingRepository)
	// sharingHandler := httphdl.NewSharingHandler(sharingService)

	// syncRepository := repo.NewSyncRepository()
	// syncService := sync.NewService(syncRepository)
	// SyncHandler := http3hdl.NewSyncHandler(syncService)

	return &App{
		repo:  repo,
		Proto: proto,
		SigCh: sigCh,
		// HistoryHandler:      HistoryHandler,
		// MetadataHandler:     MetadataHandler,
		// RegistrationHandler: RegistrationHandler,
		// ServerHandler:       ServerHandler,
		// SharingHandler:      SharingHandler,
		// SyncHandler:         SyncHandler,
	}, nil
}

func (a *App) Start() {
	// start quics protocol server
	err := a.Proto.Start()
	if err != nil {
		log.Println("quics: ", err)
	}
}

func (a *App) Close() {
	// if pressed ctrl + c, then stop server with closing database
	<-a.SigCh
	a.repo.Close()
	fmt.Println("Database is closed successfully.")
}
