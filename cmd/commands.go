package main

import (
	"fmt"
	"github.com/quic-s/quics/pkg/server"
	"github.com/spf13/cobra"
)

var id string
var all bool

// Root command
var rootCmd = &cobra.Command{
	Use:   "qis",
	Short: "qis is a CLI for interacting with the quics server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("If you need some help, then enter 'qis help'")
	},
}

// Start server
var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start server",
	RunE: func(cmd *cobra.Command, args []string) error {
		//fmt.Println("Start cobra server")
		// rest start
		return server.StartServer()
	},
}

// Stop server
var stopServerCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Stop server")
		server.StopServer()
		return nil
	},
}

// Reboot server
var rebootServerCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Reboot server")
		server.RebootServer()
		return nil
	},
}

// Shutdown server
var shutdownServerCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shutdown server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Shutdown server")
		server.ShutdownServer()
		return nil
	},
}

// Show server logs
var serverLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Check server log",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Check server log")
		return nil
	},
}

// Show client status
var getClientConnectionStatusCmd = &cobra.Command{
	Use:   "clients",
	Short: "Show client status",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Show client status")
		return nil
	},
}

// Show all clients
var getAllClientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "Show status of all clients",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Get all clients")
		return nil
	},
}

// Disconnect client
var disconnectClientCmd = &cobra.Command{
	Use:   "discc",
	Short: "Disconnect client",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Disconnect client")
		return nil
	},
}

// Disconnect all clients
var disconnectAllClientsCmd = &cobra.Command{
	Use:   "discc",
	Short: "Disconnect all clients",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Disconnect all clients")
		return nil
	},
}

// Force to disconnect file sync
var disconnectFileSyncCmd = &cobra.Command{
	Use:   "discf",
	Short: "Disconnect file sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Disconnect file sync")
		return nil
	},
}

// Force to disconnect sync of all files
var disconnectAllFileSyncCmd = &cobra.Command{
	Use:   "discf",
	Short: "Disconnect all files sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Disconnect all files sync")
		return nil
	},
}

// Show dashboard
var getDashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Open dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Open dashboard")
		return nil
	},
}

func init() {
	// Setup flags per command
	getClientConnectionStatusCmd.Flags().StringVarP(&id, "id", "i", "", "Show client status")
	getAllClientsCmd.Flags().BoolVarP(&all, "all", "a", false, "Show all clients")
	disconnectClientCmd.Flags().StringVarP(&id, "id", "i", "", "Disconnect client")
	disconnectAllClientsCmd.Flags().BoolVarP(&all, "all", "a", false, "Disconenct all clients")
	disconnectFileSyncCmd.Flags().StringVarP(&id, "id", "i", "", "Disconnect this file sync")
	disconnectAllFileSyncCmd.Flags().BoolVarP(&all, "all", "a", false, "Disconnect all files sync")

	// Commands
	rootCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(stopServerCmd)
	rootCmd.AddCommand(rebootServerCmd)
	rootCmd.AddCommand(shutdownServerCmd)
	rootCmd.AddCommand(serverLogCmd)
	rootCmd.AddCommand(getClientConnectionStatusCmd)
	rootCmd.AddCommand(getAllClientsCmd)
	rootCmd.AddCommand(disconnectClientCmd)
	rootCmd.AddCommand(disconnectAllClientsCmd)
	rootCmd.AddCommand(disconnectFileSyncCmd)
	rootCmd.AddCommand(disconnectAllFileSyncCmd)
	rootCmd.AddCommand(getDashboardCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
