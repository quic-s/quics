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
* `qis show dir --id <directory-path>`: Show directory information
* `qis show dir --all`: Show all directories information
* `qis show file --id <file-path>`: Show file information
* `qis show file --all`: Show all files information
* `qis show history --id <file-history-key>`: Show history information
* `qis show history --all`: Show all history information
*
* `qis remove`: Initialize quic-s server (needed options)
* `qis remove client --id <client-UUID>`: Initialize client
* `qis remove client --all`: Initialize all clients
* `qis remove dir --id <directory-path>`: Initialize directory
* `qis remove dir --all`: Initialize all directories
* `qis remove file --id <file-path>`: Initialize file
* `qis remove file --all`: Initialize all files
*
* `qis download file --path --version --target`: Download certain file
 */

/**
* Options & Short Options
*
* `--all`: All option
* `-a`: All short option
*
* `--id`: ID option
* `-i`: ID short option
*
* `--path`: Path option
* `-p`: Path short option
*
* `--version`: Version option
* `-v`: Version short option
*
* `--target`: Target(=destination directory) option
* `--t`: Target short option
 */

const (
	RootCommand     = "qis"
	StartCommand    = "start"
	StopCommand     = "stop"
	ListenCommand   = "listen"
	ShowCommand     = "show"
	RemoveCommand   = "remove"
	DownloadCommand = "download"

	ClientCommand  = "client"
	DirCommand     = "dir"
	FileCommand    = "file"
	HistoryCommand = "history"
)

const (
	// --all, -a
	AllOption      = "all"
	AllShortOption = "a"

	// --id, -i
	IDOption       = "id"
	IDShortCommand = "i"

	// --path, -p
	PathOption       = "path"
	PathShortCommand = "p"

	// --version, -v
	VersionOption       = "version"
	VersionShortCommand = "v"

	// --target, -t
	TargetOption       = "target"
	TargetShortCommand = "t"
)

var (
	all     bool   = false
	id      string = ""
	path    string = ""
	version uint64 = 0
	target  string = ""
)

var rootCmd = &cobra.Command{
	Use:   RootCommand,
	Short: "qis is a CLI for interacting with the quics server",
}

var (
	startServerCmd  *cobra.Command
	stopServerCmd   *cobra.Command
	listenCmd       *cobra.Command
	showCmd         *cobra.Command
	showClientCmd   *cobra.Command
	showDirCmd      *cobra.Command
	showFileCmd     *cobra.Command
	showHistoryCmd  *cobra.Command
	removeCmd       *cobra.Command
	removeClientCmd *cobra.Command
	removeDirCmd    *cobra.Command
	removeFileCmd   *cobra.Command
	downloadCmd     *cobra.Command
	downloadFileCmd *cobra.Command
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
	showHistoryCmd = initShowHistoryCmd()
	removeCmd = initRemoveCmd()
	removeClientCmd = initRemoveClientCmd()
	removeDirCmd = initRemoveDirCmd()
	removeFileCmd = initRemoveFileCmd()
	downloadCmd = initDownloadCmd()
	downloadFileCmd = initDownloadFileCmd()

	// set flags (= options)
	// qis show client --id, qis show client --all
	showClientCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Show all status")
	showClientCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Show status by ID")
	// qis show dir --id, qis show dir --all
	showDirCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Show all status")
	showDirCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Show status by ID")
	// qis show file --id, qis show file --all
	showFileCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Show all status")
	showFileCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Show status by ID")
	// qis show history --id, qis show history --all
	showHistoryCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Show all status")
	showHistoryCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Show status by ID")
	// qis remove client --id, qis remove client --all
	removeClientCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Initialize all data")
	removeClientCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Initialize by ID")
	// qis remove dir --id, qis remove dir --all
	removeDirCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Initialize all data")
	removeDirCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Initialize by ID")
	// qis remove file --id, qis remove file --all
	removeFileCmd.Flags().BoolVarP(&all, AllOption, AllShortOption, false, "Initialize all data")
	removeFileCmd.Flags().StringVarP(&id, IDOption, IDShortCommand, "", "Initialize by ID")
	// qis download file --path --version
	downloadFileCmd.Flags().StringVarP(&path, PathOption, PathShortCommand, "", "Download a file by path")
	downloadFileCmd.Flags().Uint64VarP(&version, VersionOption, VersionShortCommand, 0, "Download a file by version")
	downloadFileCmd.Flags().StringVarP(&target, TargetOption, TargetShortCommand, "", "Download location")

	// add command to root command
	rootCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(stopServerCmd)
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(downloadCmd)

	// add command to show command
	showCmd.AddCommand(showClientCmd)
	showCmd.AddCommand(showDirCmd)
	showCmd.AddCommand(showFileCmd)
	showCmd.AddCommand(showHistoryCmd)

	// add command to remove command
	removeCmd.AddCommand(removeClientCmd)
	removeCmd.AddCommand(removeDirCmd)
	removeCmd.AddCommand(removeFileCmd)

	// add command to download command
	downloadCmd.AddCommand(downloadFileCmd)

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

func initShowHistoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   HistoryCommand,
		Short: "show history information",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(showHistoryCmd)

			url := "/api/v1/server/logs/histories"
			url = getUrlWithQueryString(url)

			restClient := NewRestClient()

			_, err := restClient.GetRequest(url) // /history
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

func initRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   RemoveCommand,
		Short: "initialize quic-s server",
	}
}

func initRemoveClientCmd() *cobra.Command {
	return &cobra.Command{
		Use:   ClientCommand,
		Short: "remove client",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(removeClientCmd)

			url := "/api/v1/server/remove/clients"
			url = getUrlWithQueryString(url)

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil)
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

func initRemoveDirCmd() *cobra.Command {
	return &cobra.Command{
		Use:   DirCommand,
		Short: "initialize directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(removeDirCmd)

			url := "/api/v1/server/remove/directories"
			url = getUrlWithQueryString(url)

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil)
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

func initRemoveFileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   FileCommand,
		Short: "initialize file",
		RunE: func(cmd *cobra.Command, args []string) error {
			validateOptionByCommand(removeFileCmd)

			url := "/api/v1/server/remove/files"
			url = getUrlWithQueryString(url)

			restClient := NewRestClient()

			_, err := restClient.PostRequest(url, "application/json", nil)
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

func initDownloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   DownloadCommand,
		Short: "download certain file",
	}
}

func initDownloadFileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   FileCommand,
		Short: "download certain file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" || version == 0 || target == "" {
				log.Println("quics: ", "Please enter both path and version")
				cmd.Help()
				return nil
			}

			url := "/api/v1/server/download/files"
			url = getUrlWithQueryString(url)

			restClient := NewRestClient()

			_, err := restClient.GetRequest(url)
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

// ********************************************************************************
//                                  Private Logic
// ********************************************************************************

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
