package cmd

import (
	"os"

	"github.com/elliot40404/acc/pkg/database"
	"github.com/elliot40404/acc/pkg/utils"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Initializes the application",
	Example: `acc init`,
	Run:     InitApp,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func InitApp(cmd *cobra.Command, args []string) {
	// check if the application has already been initialized
	if utils.IsInitialized() {
		utils.PrintError(nil, "Application has already been initialized", false)
		os.Exit(0)
	}

	dbPath := utils.DBPATH()
	appDir := utils.APPDIR()
	debugMode, _ := cmd.Flags().GetBool("verbose")
	// make .acc directory if it doesn't exist
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err = os.MkdirAll(appDir, 0755)
		if err != nil {
			utils.PrintError(err, "Failed to create application directory", debugMode)
		}
	}
	// create database if it doesn't exist
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		_, err := os.Create(dbPath)
		if err != nil {
			utils.PrintError(err, "Failed to create database", debugMode)
		}
	}
	// load schema.sql into database
	err := database.InitApplication()
	if err != nil {
		utils.PrintError(err, "Failed to initialize database", debugMode)
	}
	// create marker to indicate that the database has been initialized
	_, err = os.Create(appDir + "/.initialized")
	if err != nil {
		utils.PrintError(err, "Failed to create init file", debugMode)
	}
}
