package cmd

import (
	"fmt"

	"github.com/elliot40404/acc/cmd/list"
	"github.com/elliot40404/acc/pkg/database"
	"github.com/elliot40404/acc/pkg/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all the transactions",
	Run:   List,
}
var columns = []string{}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("interactive", "i", false, "interactive mode")
	listCmd.Flags().Bool("date-help", false, "show help for supported date formats")
	listCmd.Flags().Bool("pretty", false, "pretty print json")
	listCmd.Flags().BoolP("all", "x", false, "print all transactions")
	listCmd.Flags().BoolP("asc", "A", false, "sort in ascending order")
	listCmd.Flags().BoolP("htime", "H", false, "human friendly time format")
	listCmd.Flags().IntP("page", "p", 1, "page number")
	listCmd.Flags().IntP("limit", "l", 10, "limit per page")
	listCmd.Flags().StringP("date", "d", "", "filter by date or date-ranges")
	listCmd.Flags().StringP("type", "t", "", "filter by type (income, expense)")
	listCmd.Flags().StringP("amount", "a", "", "filter by amount")
	listCmd.Flags().StringP("sort", "s", "", "sort by date, amount")
	listCmd.Flags().StringP("desc", "D", "", "filter by description")
	listCmd.Flags().StringP("format", "f", "table", "print in table/json/csv format")
	listCmd.Flags().StringSliceVarP(&columns, "columns", "c", []string{}, "columns to print (id, type, amt, desc, date) (default: all) (only works with table format) (example: -c 'id,type' or -c id -c type)")
}

func List(cmd *cobra.Command, args []string) {
	queryConfig := database.TransactionConfig{
		Dry:      cmd.Flag("dry").Value.String() == "true",
		Verbose:  cmd.Flag("verbose").Value.String() == "true",
		TxType:   cmd.Flag("type").Value.String(),
		Page:     utils.ToInt(cmd.Flag("page").Value.String()),
		Limit:    utils.ToInt(cmd.Flag("limit").Value.String()),
		All:      cmd.Flag("all").Value.String() == "true",
		Date:     cmd.Flag("date").Value.String(),
		Amount:   cmd.Flag("amount").Value.String(),
		Sort:     cmd.Flag("sort").Value.String(),
		SortAsc:  cmd.Flag("asc").Value.String() == "true",
		Desc:     cmd.Flag("desc").Value.String(),
		Columns:  columns,
		IsHRTime: cmd.Flag("htime").Value.String() == "true",
		Format:   cmd.Flag("format").Value.String(),
		IsPretty: cmd.Flag("pretty").Value.String() == "true",
	}
	if cmd.Flag("date-help").Value.String() == "true" {
		printDateHelp()
		return
	}
	if queryConfig.Verbose {
		fmt.Printf("%#v\n", queryConfig)
	}
	err := list.ValidateConfig(&queryConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	isInteractive := cmd.Flag("interactive").Value.String() == "true"
	if isInteractive {
		list.InteractiveListRenderer(queryConfig)
		return
	}
	list.NonInteractiveListRenderer(queryConfig)
}

func printDateHelp() {
	fmt.Println("Supported date formats:")
	fmt.Println("Builtins: today, yesterday, thisweek, lastweek, thismonth, lastmonth, thisyear, lastyear")
	fmt.Println("Specific date formats: YYYY-MM-DD, DD-MM-YYYY, MM-DD-YYYY")
	fmt.Println("Date ranges:")
	fmt.Println(":date    - all transactions before date")
	fmt.Println("date: 	 - all transactions after date")
	fmt.Println("date:date - all transactions between dates")
}
