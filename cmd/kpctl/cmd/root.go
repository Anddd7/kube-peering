/*
Copyright © 2023 Anddd7 <liaoad_space@sina.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kpctl",
	Short: "setup your local environment and connect with remote kpeering",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
