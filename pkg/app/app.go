package app

import (
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

// New initialize program
func New(repo *badger.Badger) (*App, error) {
	// define system call actions
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

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

	err = proto.RecvTransactionHandleFunc(types.REGISTERCLIENT, registrationHandler.RegisterClient)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}
	err = proto.RecvTransactionHandleFunc(types.REGISTERROOTDIR, registrationHandler.RegisterRootDir)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}
	err = proto.RecvTransactionHandleFunc(types.GETROOTDIRS, registrationHandler.GetRemoteDirs)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	return &App{
		repo:  repo,
		Proto: proto,
		SigCh: sigCh,
	}, nil
}

func (a *App) Start() error {
	// start quics protocol server
	err := a.Proto.Start()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

func (a *App) Close() error {
	// if pressed ctrl + c, then stop server with closing database
	<-a.SigCh
	err := a.repo.Close()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}
