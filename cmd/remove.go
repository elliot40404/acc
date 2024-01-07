package cmd

import (
	"fmt"

	"github.com/elliot40404/acc/pkg/database"
	"github.com/elliot40404/acc/pkg/utils"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove transactions by id or all transactions",
	Run:   Remove,
}

var ids []string

func init() {
	RootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolP("all", "A", false, "Remove all transactions")
	removeCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, "Remove a transaction by id. Example: -i 1 -i 2 or --id\"1,2\"")
	removeCmd.MarkFlagsMutuallyExclusive("all", "id")
}

func Remove(cmd *cobra.Command, args []string) {
	db := database.NewTransactionRepository()
	all, _ := cmd.Flags().GetBool("all")
	dc := database.DeleteConfig{
		Dry:     cmd.Flags().Changed("dry"),
		Verbose: cmd.Flags().Changed("verbose"),
		All:     all,
		Ids:     ids,
	}
	if dc.Verbose {
		fmt.Printf("%+v\n", dc)
	}
	if utils.PromptConfirmation() {
		err := db.DeleteTransactions(dc)
		if err != nil {
			panic(err)
		}
		return
	} else {
		fmt.Println("Aborted")
	}
	fmt.Println("Please specify an id or use --all to remove all transactions")
	// TODO: maybe have a backup option
}
