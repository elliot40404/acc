package cmd

import (
	"fmt"
	"log/slog"

	"github.com/elliot40404/acc/pkg/database"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a new income or expense",
	Example: `acc add -t income -d "Paycheck" -a 1000`,
	Run:     Add,
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("type", "t", "", "Type of transaction (income or expense)")
	addCmd.Flags().StringP("description", "d", "", "Description")
	addCmd.Flags().Float64P("amount", "a", 0, "Amount")
	addCmd.MarkFlagRequired("type")
	addCmd.MarkFlagRequired("description")
	addCmd.MarkFlagRequired("amount")
}

func Add(cmd *cobra.Command, args []string) {
	transactionType, _ := cmd.Flags().GetString("type")
	description, _ := cmd.Flags().GetString("description")
	amount, _ := cmd.Flags().GetFloat64("amount")
	if transactionType != "income" && transactionType != "expense" {
		slog.Error("Invalid transaction type", "Type", transactionType)
		fmt.Println("Supported transaction types: income, expense")
		return
	}
	fmt.Printf("Adding %s transaction: %s for $%.2f\n", transactionType, description, amount)
	db := database.NewTransactionRepository()
	err := db.CreateTransaction(database.Transaction{
		Type:        transactionType,
		Description: description,
		Amount:      amount,
	})
	if err != nil {
		slog.Error("Failed to create transaction", "Error", err.Error())
	}
}
