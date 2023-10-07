package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/quic-s/quics/config"
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
	SigCh := make(chan os.Signal, 1)
	signal.Notify(SigCh, syscall.SIGINT, syscall.SIGTERM)

	// initialize server password
	Password := config.GetViperEnvVariables("PASSWORD")

	repo, err := badger.NewBadgerRepository()
	if err != nil {
		log.Println("quics: Error while connecting to the database: ", err)
	}

	pool := connection.NewnPool()

	proto, err := qp.New("0.0.0.0", 6122, pool)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	registrationRepository := repo.NewRegistrationRepository()
	registrationNetworkAdapter := qp.NewRegistrationAdapter(pool)
	registrationService := registration.NewService(Password, registrationRepository, registrationNetworkAdapter)
	registrationHandler := qp.NewRegistrationHandler(registrationService)

	proto.RecvTransactionHandleFunc(types.REGISTERCLIENT, registrationHandler.RegisterClient)
	proto.RecvTransactionHandleFunc(types.REGISTERROOTDIR, registrationHandler.RegisterRootDir)
	proto.RecvTransactionHandleFunc(types.GETROOTDIRS, registrationHandler.GetRemoteDirs)

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
		SigCh: SigCh,
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
