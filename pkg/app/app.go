package app

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-s/quics/pkg/config"
	"github.com/quic-s/quics/pkg/core/server"
	"github.com/quic-s/quics/pkg/core/sharing"
	quicshttp "github.com/quic-s/quics/pkg/network/http"
	quicshttp3 "github.com/quic-s/quics/pkg/network/http3"
	"github.com/quic-s/quics/pkg/repository/badger"
	"github.com/quic-s/quics/pkg/utils"
)

type App struct {
	certFileDir string
	keyFileDir  string
	restServer  *http3.Server
}

// New initialize program
func New(ip string, port string) (*App, error) {
	repo, err := badger.NewBadgerRepository()
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	serverRepository := repo.NewServerRepository()
	historyRepository := repo.NewHistoryRepository()
	syncRepository := repo.NewSyncRepository()
	sharingRepository := repo.NewSharingRepository()

	serverService, err := server.NewService(repo, serverRepository)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}
	sharingService := sharing.NewService(historyRepository, syncRepository, sharingRepository)

	serverHandler := quicshttp3.NewServerHandler(serverService)
	sharingHandler := quicshttp.NewSharingHandler(sharingService)

	mux := http.NewServeMux()
	serverHandler.SetupRoutes(mux)
	sharingHandler.SetupRoutes(mux)

	restServer := &http3.Server{
		Addr:       config.GetHttp3ServerAddress(ip, port),
		QuicConfig: &quic.Config{},
		Handler:    mux,
	}

	// get directory path for certification
	quicsDir := utils.GetQuicsDirPath()
	certFileDir := filepath.Join(quicsDir, config.GetViperEnvVariables("QUICS_CERT_NAME"))
	keyFileDir := filepath.Join(quicsDir, config.GetViperEnvVariables("QUICS_KEY_NAME"))

	// load the certificate and the key from the files
	_, err = tls.LoadX509KeyPair(certFileDir, keyFileDir)
	if err != nil {
		err = config.CreateSecurityFiles()
		if err != nil {
			log.Println("quics: ", err)
			return nil, err
		}
	}

	fmt.Println("************************************************************")
	fmt.Println("                           Start                            ")
	fmt.Println("************************************************************")

	return &App{
		certFileDir: certFileDir,
		keyFileDir:  keyFileDir,
		restServer:  restServer,
	}, nil
}

func (a *App) Start() error {
	err := a.restServer.ListenAndServeTLS(a.certFileDir, a.keyFileDir)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
	return nil
}

func (a *App) Close() error {
	// define system call actions
	interruptCh := make(chan os.Signal, 1) // buffered channel
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	// if pressed ctrl + c, then stop server with closing database
	go func() {
		<-interruptCh

		os.Exit(0)
	}()

	return nil
}

func (a *App) Stop() error {
	return nil
}
