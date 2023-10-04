package app

import (
	"fmt"
	"github.com/quic-s/quics/config"
	"github.com/quic-s/quics/pkg/core/download"
	"github.com/quic-s/quics/pkg/core/history"
	"github.com/quic-s/quics/pkg/core/metadata"
	"github.com/quic-s/quics/pkg/core/registration"
	"github.com/quic-s/quics/pkg/core/server"
	"github.com/quic-s/quics/pkg/core/sync"
	httphdl "github.com/quic-s/quics/pkg/handlers/http"
	http3hdl "github.com/quic-s/quics/pkg/handlers/http3"
	repo "github.com/quic-s/quics/pkg/repositories/badger"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
}

var SigCh chan os.Signal
var Password string

var DownloadHandler *httphdl.DownloadHandler
var HistoryHandler *http3hdl.HistoryHandler
var MetadataHandler *http3hdl.MetadataHandler
var RegistrationHandler *http3hdl.RegistrationHandler
var ServerHandler *http3hdl.ServerHandler
var SyncHandler *http3hdl.SyncHandler

// Initialize initialize program
func Initialize() (*App, error) {
	// define system call actions
	SigCh = make(chan os.Signal, 1)
	signal.Notify(SigCh, syscall.SIGINT, syscall.SIGTERM)

	// initialize server password
	Password = config.GetViperEnvVariables("PASSWORD")

	return &App{}, nil
}

func (app *App) InitAdapters() {
	repo.NewBadgerDB()

	downloadRepository := repo.NewDownloadRepository()
	downloadService := download.NewDownloadService(downloadRepository)
	DownloadHandler = httphdl.NewDownloadHandler(downloadService)

	historyRepository := repo.NewHistoryRepository()
	historyService := history.NewHistoryService(historyRepository)
	HistoryHandler = http3hdl.NewHistoryHandler(historyService)

	metadataRepository := repo.NewMetadataRepository()
	metadataService := metadata.NewMetadataService(metadataRepository)
	MetadataHandler = http3hdl.NewMetadataHandler(metadataService)

	registrationRepository := repo.NewRegistrationRepository()
	registrationService := registration.NewRegistrationService(registrationRepository)
	RegistrationHandler = http3hdl.NewRegistrationHandler(registrationService)

	serverRepository := repo.NewServerRepository()
	serverService := server.NewServerService(serverRepository)
	ServerHandler = http3hdl.NewServerHandler(serverService)

	syncRepository := repo.NewSyncRepository()
	syncService := sync.NewSyncService(syncRepository)
	SyncHandler = http3hdl.NewSyncHandler(syncService)
}

func (app *App) Close() {
	// if pressed ctrl + c, then stop server with closing database
	<-SigCh
	repo.CloseBadgerDB()
	fmt.Println("Database is closed successfully.")
}
