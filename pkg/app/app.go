package app

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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

type App struct {
	repo     *badger.Badger
	Proto    *qp.Protocol
	Password string

	registrationService registration.Service
	syncService         sync.Service
}

// New initialize program
func New(repo *badger.Badger) (*App, error) {
	// get env variables (server password, port)
	password := config.GetViperEnvVariables("PASSWORD")
	port, err := strconv.Atoi(config.GetViperEnvVariables("QUICS_PORT"))
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
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

	historyRepository := repo.NewHistoryRepository()
	syncNetworkAdapter := qp.NewSyncAdapter(pool)
	syncDirAdapter := fs.NewSyncDir(utils.GetQuicsSyncDirPath())

	syncRepository := repo.NewSyncRepository()
	syncService := sync.NewService(registrationRepository, historyRepository, syncRepository, syncNetworkAdapter, syncDirAdapter)
	syncHandler := qp.NewSyncHandler(syncService)

	proto.RecvTransactionHandleFunc(types.REGISTERROOTDIR, syncHandler.RegisterRootDir)
	proto.RecvTransactionHandleFunc(types.SYNCROOTDIR, syncHandler.SyncRootDir)
	proto.RecvTransactionHandleFunc(types.GETROOTDIRS, syncHandler.GetRemoteDirs)
	proto.RecvTransactionHandleFunc(types.PLEASESYNC, syncHandler.PleaseSync)
	proto.RecvTransactionHandleFunc(types.CONFLICTLIST, syncHandler.AskConflictList)
	proto.RecvTransactionHandleFunc(types.CHOOSEONE, syncHandler.ChooseOne)
	proto.RecvTransactionHandleFunc(types.RESCAN, syncHandler.Rescan)

	return &App{
		repo:                repo,
		Proto:               proto,
		registrationService: registrationService,
		syncService:         syncService,
	}, nil
}

func (a *App) Start() {
	// start quics protocol server
	a.syncService.BackgroundFullScan(300)
	err := a.Proto.Start()
	if err != nil {
		log.Println("quics: ", err)
		return
	}
}

func (a *App) Close() error {
	// define system call actions
	interruptCh := make(chan os.Signal)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	// if pressed ctrl + c, then stop server with closing database
	go func() {
		<-interruptCh

		err := a.repo.Close()
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		os.Exit(0)
	}()

	return nil
}

func (a *App) Stop() error {

	go func() {
		err := a.repo.Close()
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		log.Println("quics: Closed")
		os.Exit(0)
	}()

	return nil
}
