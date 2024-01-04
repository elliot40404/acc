package cmd

import (
	"fmt"
	"os"

	"github.com/elliot40404/acc/pkg/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "acc",
	Short: "A simple cli to manage income and expenses",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = "0.0.1-alpha"
	rootCmd.PersistentFlags().Bool("verbose", false, "Prints debug messages")
	rootCmd.PersistentFlags().Bool("trace", false, "Prints trace messages")
	rootCmd.PersistentFlags().Bool("dry", false, "Dry run")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "init" {
			checkInitialized()
		}
	}
	// TODO: I should be able to see the TRACES in debug mode
}

func checkInitialized() {
	if !utils.IsInitialized() {
		fmt.Println("acc has not been initialized. Run 'acc init' to initialize the application")
		os.Exit(0)
	}

}
