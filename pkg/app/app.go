package app

import (
	"crypto/tls"
	"errors"
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
	"github.com/quic-s/quics/pkg/fs"
	quicshttp "github.com/quic-s/quics/pkg/network/http"
	"github.com/quic-s/quics/pkg/repository/badger"
	"github.com/quic-s/quics/pkg/utils"
)

type App struct {
	certFileDir   string
	keyFileDir    string
	serverService server.Service
	entryServer   *http.Server
	restServer    *http3.Server
}

// New initialize program
func New(ip string, port string, port3 string) (*App, error) {
	err := config.SetServerAddress(ip, port, port3)
	if err != nil {
		err = errors.New("[App.New] setting server address: " + err.Error())
		return nil, err
	}

	repo, err := badger.NewBadgerRepository()
	if err != nil {
		err = errors.New("[App.New] initializing badger repository: " + err.Error())
		return nil, err
	}

	serverRepository := repo.NewServerRepository()
	historyRepository := repo.NewHistoryRepository()
	syncRepository := repo.NewSyncRepository()
	sharingRepository := repo.NewSharingRepository()

	syncDirAdapter := fs.NewSyncDir(utils.GetQuicsSyncDirPath())

	serverService, err := server.NewService(repo, serverRepository, syncDirAdapter)
	if err != nil {
		err = errors.New("[App.New] initializing server service: " + err.Error())
		return nil, err
	}

	sharingService := sharing.NewService(historyRepository, syncRepository, sharingRepository, syncDirAdapter)

	serverHandler := quicshttp.NewServerHandler(serverService)
	sharingHandler := quicshttp.NewSharingHandler(sharingService)

	mux := http.NewServeMux()
	serverHandler.SetupRoutes(mux)
	sharingHandler.SetupRoutes(mux)

	restServer := &http3.Server{
		Addr:       "0.0.0.0:" + config.GetViperEnvVariables("REST_SERVER_H3_PORT"),
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
			err = errors.New("[App.New] creating security files: " + err.Error())
			return nil, err
		}
	}

	// set legacy http for first connection
	entryServer := &http.Server{
		Addr:    "0.0.0.0:" + config.GetViperEnvVariables("REST_SERVER_PORT"),
		Handler: mux,
	}

	return &App{
		certFileDir:   certFileDir,
		keyFileDir:    keyFileDir,
		serverService: serverService,
		entryServer:   entryServer,
		restServer:    restServer,
	}, nil
}

func (a *App) StartRestServer() error {
	fmt.Println("************************************************************")
	fmt.Println("                   Start Rest Server                        ")
	fmt.Println("************************************************************")
	go func() {
		err := a.entryServer.ListenAndServeTLS(a.certFileDir, a.keyFileDir)
		if err != nil {
			err = errors.New("[App.Start] starting rest server: " + err.Error())
			log.Fatalln("quics err: ", err)
		}
	}()
	err := a.restServer.ListenAndServeTLS(a.certFileDir, a.keyFileDir)
	if err != nil {
		err = errors.New("[App.Start] starting rest server: " + err.Error())
		log.Fatalln("quics err: ", err)

		return err
	}
	return nil
}

func (a *App) Run() error {
	go a.StartRestServer()
	err := a.serverService.ListenProtocol()
	if err != nil {
		err = errors.New("[App.Run] listening protocol: " + err.Error())
		log.Fatalln("quics err: ", err)
		return err
	}
	err = a.Close()
	if err != nil {
		err = errors.New("[App.Run] closing: " + err.Error())
		log.Fatalln("quics err: ", err)
		return err
	}
	return nil
}

func (a *App) Close() error {
	// define system call actions
	interruptCh := make(chan os.Signal, 1) // buffered channel
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	// if pressed ctrl + c, then stop server with closing database
	<-interruptCh
	a.serverService.StopServer()

	fmt.Println("************************************************************")
	fmt.Println("                           Close                            ")
	fmt.Println("************************************************************")
	os.Exit(0)

	return nil
}

func (a *App) Stop() error {
	return nil
}
