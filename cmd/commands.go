package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qis",
	Short: "qis is a CLI for interacting with the quics server",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {

}
