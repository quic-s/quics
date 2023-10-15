package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"path/filepath"

	"github.com/quic-go/quic-go"
	http3 "github.com/quic-go/quic-go/http3"
	"github.com/quic-s/quics/pkg/app"
	"github.com/quic-s/quics/pkg/config"
	"github.com/quic-s/quics/pkg/core/server"
	customhttp3 "github.com/quic-s/quics/pkg/network/http3"
	"github.com/quic-s/quics/pkg/repository/badger"
	"github.com/quic-s/quics/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	RootCommand             = "qis"
	StartCommand            = "start"
	StopCommand             = "stop"
	ListenCommand           = "listen"
	RebootCommand           = "reboot"
	ShutdownCommand         = "shutdown"
	LogCommand              = "log"
	ClientCommand           = "clients"
	DisconnectClientCommand = "discc"
	DisconnectFileCommand   = "discf"
	DashboardCommand        = "dashboard"

	ClientIDCommand      = "id"
	ClientIDShortCommand = "i"

	AllCommand      = "all"
	AllShortCommand = "a"
)

var rootCmd = &cobra.Command{
	Use:   RootCommand,
	Short: "qis is a CLI for interacting with the quics server",
}

func Run() int {
	// initialize
	startServerCmd := initStartServerCmd()
	stopServerCmd := initStopServerCmd()
	listenCmd := initListenCmd()

	// add command
	rootCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(stopServerCmd)
	rootCmd.AddCommand(listenCmd)

	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}

// initStartServerCmd start quic-s server (`qis start`)
func initStartServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   StartCommand,
		Short: "start quic-s server",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := badger.NewBadgerRepository()
			if err != nil {
				log.Println("quics: ", err)
				return err
			}

			initApp, err := app.New(repo)
			if err != nil {
				return err
			}

			serverRepository := repo.NewServerRepository()
			serverService := server.NewService(initApp, serverRepository)
			serverHandler := customhttp3.NewServerHandler(serverService)

			handler := serverHandler.SetupRoutes()

			restServer := &http3.Server{
				Addr:       "0.0.0.0:" + config.GetViperEnvVariables("REST_SERVER_PORT"),
				QuicConfig: &quic.Config{},
				Handler:    handler,
			}

			// get directory path for certification
			quicsDir := utils.GetQuicsDirPath()
			certFileDir := filepath.Join(quicsDir, config.GetViperEnvVariables("QUICS_CERT_NAME"))
			keyFileDir := filepath.Join(quicsDir, config.GetViperEnvVariables("QUICS_KEY_NAME"))

			// load the certificate and the key from the files
			_, err = tls.LoadX509KeyPair(certFileDir, keyFileDir)
			if err != nil {
				config.SecurityFiles()
			}

			fmt.Println("************************************************************")
			fmt.Println("                           Start                            ")
			fmt.Println("************************************************************")

			err = restServer.ListenAndServeTLS(certFileDir, keyFileDir)
			if err != nil {
				log.Println("quics: ", err)
				return err
			}

			return nil
		},
	}
}

// initStopServerCmd stop quic-s server (`qis stop`)
func initStopServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   StopCommand,
		Short: "stop quic-s server",
		RunE: func(cmd *cobra.Command, args []string) error {
			restClient := NewRestClient()

			_, err := restClient.PostRequest("/api/v1/server/stop", "application/json", nil) // /server/stop
			if err != nil {
				log.Println("quics: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics: ", err)
				return err
			}

			return nil
		},
	}
}

// initListenCmd listen quic-s protocol (`qis listen`)
func initListenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   ListenCommand,
		Short: "listen quic-s protocol",
		RunE: func(cmd *cobra.Command, args []string) error {
			restClient := NewRestClient()

			_, err := restClient.PostRequest("/api/v1/server/listen", "application/json", nil) // /server/listen
			if err != nil {
				log.Println("quics: ", err)
				return err
			}

			err = restClient.Close()
			if err != nil {
				log.Println("quics: ", err)
				return err
			}

			return nil
		},
	}
}
