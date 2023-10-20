package main

import (
	"log"

	"github.com/quic-s/quics/pkg/app"
	"github.com/spf13/cobra"
)

/**
* Commands
*
* `qis`: Root command (meaning quic-s)
*
* `qis start`: Start quic-s server
* `qis stop`: Stop quic-s server
* `qis listen`: Listen quic-s protocol
*
* `qis show`: Show quic-s server information (needed options)
* `qis show client --id <client-UUID>`: Show client information
* `qis show client --all`: Show all clients information
* `qis show dir <directory-path>`: Show directory information
* `qis show dir --all`: Show all directories information
* `qis show file <file-path>`: Show file information
* `qis show file --all`: Show all files information
*
* `qis disconnect`: Initialize quic-s server (needed options)
* `qis disconnect client --id <client-UUID>`: Initialize client
* `qis disconnect client --all`: Initialize all clients
* `qis disconnect dir <directory-path>`: Initialize directory
* `qis disconnect dir --all`: Initialize all directories
* `qis disconnect file <file-path>`: Initialize file
* `qis disconnect file --all`: Initialize all files
 */

/**
* Options & Short Options
*
* `--all`: All option
* `-a`: All short option
*
* `--id`: ID option
* `-i`: ID short option
 */

const (
	// root
	RootCommand = "qis"

	// server
	StartCommand      = "start"
	StopCommand       = "stop"
	ListenCommand     = "listen"
	ShowCommand       = "show"
	DisconnectCommand = "disconnect"

	ClientCommand = "client"
	DirCommand    = "dir"
	FileCommand   = "file"
)

const (
	// --all, -a
	AllOption      = "all"
	AllShortOption = "a"

	// --id, -i
	IDOption       = "id"
	IDShortCommand = "i"
)

var (
	all bool   = false
	id  string = ""
)

var rootCmd = &cobra.Command{
	Use:   RootCommand,
	Short: "qis is a CLI for interacting with the quics server",
}

var (
	startServerCmd      *cobra.Command
	stopServerCmd       *cobra.Command
	listenCmd           *cobra.Command
	showCmd             *cobra.Command
	showClientCmd       *cobra.Command
	showDirCmd          *cobra.Command
	showFileCmd         *cobra.Command
	disconnectCmd       *cobra.Command
	disconnectClientCmd *cobra.Command
	disconnectDirCmd    *cobra.Command
	disconnectFileCmd   *cobra.Command
)

// Run initializes and executes commands using cobra library
func Run() int {
	// initialize
	startServerCmd = initStartServerCmd()
	stopServerCmd = initStopServerCmd()
	listenCmd = initListenCmd()

	showCmd = initShowCmd()
	showClientCmd = initShowClientCmd()
	showDirCmd = initShowDirCmd()
	showFileCmd = initShowFileCmd()

	disconnectCmd = initDisconnectCmd()
	disconnectClientCmd = initDisconnectClientCmd()
	disconnectDirCmd = initDisconnectDirCmd()
	disconnectFileCmd = initDisconnectFileCmd()

	// set flags (= options)
	showClientCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Show all status")
	showClientCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Show status by ID")
	showDirCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Show all status")
	showDirCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Show status by ID")
	showFileCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Show all status")
	showFileCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Show status by ID")

	disconnectClientCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Initialize all data")
	disconnectClientCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Initialize by ID")
	disconnectDirCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Initialize all data")
	disconnectDirCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Initialize by ID")
	disconnectFileCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Initialize all data")
	disconnectFileCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Initialize by ID")

	// add command
	rootCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(stopServerCmd)
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(disconnectCmd)

	showCmd.AddCommand(showClientCmd)
	showCmd.AddCommand(showDirCmd)
	showCmd.AddCommand(showFileCmd)

	disconnectCmd.AddCommand(disconnectClientCmd)
	disconnectCmd.AddCommand(disconnectDirCmd)
	disconnectCmd.AddCommand(disconnectFileCmd)

	// execute command
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
			quicsApp, err := app.New()
			if err != nil {
				return err
			}

			err = quicsApp.Start()
			if err != nil {
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
			url := "/api/v1/server/stop"

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil) // /server/stop
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
			url := "/api/v1/server/listen"

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil) // /server/listen
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

func initShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   ShowCommand,
		Short: "show quic-s server data",
	}
}

func initDisconnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   DisconnectCommand,
		Short: "initialize quic-s server",
	}
}

func initShowClientCmd() *cobra.Command {
	return &cobra.Command{
		Use:   ClientCommand,
		Short: "show client information",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(showClientCmd)

			url := "/api/v1/server/logs/clients"
			url = getUrlWithQueryString(url)

			restClient := NewRestClient()

			_, err := restClient.GetRequest(url) // /clients
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

func initShowDirCmd() *cobra.Command {
	return &cobra.Command{
		Use:   DirCommand,
		Short: "show directory information",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(showDirCmd)

			url := "/api/v1/server/logs/directories"
			url = getUrlWithQueryString(url)

			restClient := NewRestClient()

			_, err := restClient.GetRequest(url) // /directories
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

func initShowFileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   FileCommand,
		Short: "show file information",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(showFileCmd)

			url := "/api/v1/server/logs/files"
			url = getUrlWithQueryString(url)

			restClient := NewRestClient()

			_, err := restClient.GetRequest(url) // /files
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

func initDisconnectClientCmd() *cobra.Command {
	return &cobra.Command{
		Use:   ClientCommand,
		Short: "disconnect client",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(disconnectClientCmd)

			url := "/api/v1/server/disconnections/clients"
			url = getUrlWithQueryString(url)

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil) // /server/init/client
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

func initDisconnectDirCmd() *cobra.Command {
	return &cobra.Command{
		Use:   DirCommand,
		Short: "initialize directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(disconnectDirCmd)

			url := "/api/v1/server/disconnections/directories"
			url = getUrlWithQueryString(url)

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil) // /server/init/directory
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

func initDisconnectFileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   FileCommand,
		Short: "initialize file",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(disconnectFileCmd)

			url := "/api/v1/server/disconnections/files"
			url = getUrlWithQueryString(url)

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil) // /server/init/file
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

func validateOptionByCommand(command *cobra.Command) {
	if !all && id == "" {
		log.Println("quics: ", "Please enter only one option")
		command.Help()
		return
	}
}

func getUrlWithQueryString(url string) string {
	if all {
		url += "?all=true"
		return url
	}
	if id != "" {
		url += "?id=" + id
		return url
	}

	return url
}
