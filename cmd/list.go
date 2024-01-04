package cmd

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/elliot40404/acc/pkg/database"
	"github.com/elliot40404/acc/pkg/utils"
	"github.com/spf13/cobra"

	"github.com/jedib0t/go-pretty/v6/table"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all the transactions",
	Run:   List,
}

var sortables = []string{
	"date",
	"amt",
}

var validFormats = []string{
	"table",
	"json",
	"csv",
}

var validColumns = []string{
	"id",
	"type",
	"amt",
	"desc",
	"date",
}

var columns = []string{}

func init() {
	rootCmd.AddCommand(listCmd)
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
	// TODO: I should be able change pages interactively ?? use arrow keys to navigate or vim like navigation
	db := database.NewTransactionRepository()
	queryConfig := database.TransactionConfig{
		Dry:     cmd.Flag("dry").Value.String() == "true",
		Verbose: cmd.Flag("verbose").Value.String() == "true",
		TxType:  cmd.Flag("type").Value.String(),
		Page:    utils.ToInt(cmd.Flag("page").Value.String()),
		Limit:   utils.ToInt(cmd.Flag("limit").Value.String()),
		All:     cmd.Flag("all").Value.String() == "true",
		Date:    cmd.Flag("date").Value.String(),
		Amount:  cmd.Flag("amount").Value.String(),
		Sort:    cmd.Flag("sort").Value.String(),
		SortAsc: cmd.Flag("asc").Value.String() == "true",
		Desc:    cmd.Flag("desc").Value.String(),
		Columns: columns,
	}

	if cmd.Flag("date-help").Value.String() == "true" {
		printDateHelp()
		return
	}
	if queryConfig.Verbose {
		fmt.Printf("%#v\n", queryConfig)
	}
	err := validateConfig(&queryConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	transactions, err := db.GetTransactionsWithConfig(queryConfig)
	totalTx, _ := db.GetTransactionCountWithConfig(queryConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	// exit if dry run
	if cmd.Flag("dry").Value.String() == "true" {
		return
	}
	switch cmd.Flag("format").Value.String() {
	case "json":
		isPretty := cmd.Flag("pretty").Value.String() == "true"
		jsonWriter(transactions, isPretty)
	case "csv":
		csvWriter(transactions)
	default:
		tableWriter(transactions, cmd, totalTx, queryConfig)
	}
}

func validateConfig(config *database.TransactionConfig) error {
	if config.TxType != "" && config.TxType != "income" && config.TxType != "expense" {
		return errors.New("invalid type. type must be either income or expense")
	}
	if config.Page < 1 {
		return errors.New("invalid page number. page number must be greater than 0")
	}
	if config.Limit < 1 {
		return errors.New("invalid limit. limit must be greater than 0")
	}
	if config.Date != "" {
		err := utils.Checkdate(&config.Date)
		if err != nil {
			return err
		}
	}
	if config.Amount != "" {
		err := utils.CheckAmount(&config.Amount)
		if err != nil {
			return err
		}
	}
	if config.Sort != "" {
		for _, sortable := range sortables {
			if config.Sort == sortable {
				return nil
			}
		}
		return errors.New("invalid sort. sort must be either 'date' or 'amt'")
	}
	return nil
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

func tableWriter(transactions []database.Transaction, cmd *cobra.Command, totalTx int, queryConfig database.TransactionConfig) {
	t := table.NewWriter()
	var header table.Row
	if len(queryConfig.Columns) == 0 {
		header = table.Row{
			"#",
			"Type",
			"Amt",
			"Desc",
			"Date",
		}
	} else {
		for _, column := range queryConfig.Columns {
			for _, validColumn := range validColumns {
				if column == validColumn {
					header = append(header, column)
				}
			}
		}
	}
	t.AppendHeader(header)
	var rows []table.Row
	for _, transaction := range transactions {
		if cmd.Flag("htime").Value.String() == "true" {
			transaction.CreatedAt = utils.HRTime(transaction.CreatedAt)
		}
		var row table.Row
		if len(queryConfig.Columns) == 0 {
			row = table.Row{
				strconv.Itoa(transaction.ID),
				transaction.Type,
				transaction.Amount,
				transaction.Description,
				transaction.CreatedAt,
			}
		} else {
			for _, column := range queryConfig.Columns {
				switch column {
				case "id":
					row = append(row, strconv.Itoa(transaction.ID))
				case "type":
					row = append(row, transaction.Type)
				case "amt":
					row = append(row, transaction.Amount)
				case "desc":
					row = append(row, transaction.Description)
				case "date":
					row = append(row, transaction.CreatedAt)
				}
			}
		}
		rows = append(rows, row)
	}
	t.AppendRows(rows)
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()
	totalPages := int(math.Ceil(float64(totalTx) / float64(queryConfig.Limit)))
	if queryConfig.All {
		totalPages = 1
	}
	fmt.Printf(
		"Page: %d of %d | Results: %d | Total: %d\n",
		queryConfig.Page,
		totalPages,
		len(transactions),
		totalTx,
	)
}

func jsonWriter(transactions []database.Transaction, pretty bool) {
	var b []byte
	var err error
	if pretty {
		b, err = json.MarshalIndent(transactions, "", "  ")
	} else {
		b, err = json.Marshal(transactions)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func csvWriter(transactions []database.Transaction) {
	csvWriter := csv.NewWriter(os.Stdout)
	var rows [][]string
	rows = append(rows, []string{
		"ID",
		"Type",
		"Amt",
		"Desc",
		"CreatedAt",
		"UpdatedAt",
	})
	for _, transaction := range transactions {
		rows = append(rows, []string{
			strconv.Itoa(transaction.ID),
			transaction.Type,
			strconv.FormatFloat(transaction.Amount, 'f', 2, 64),
			transaction.Description,
			transaction.CreatedAt,
			transaction.UpdatedAt,
		})
	}
	csvWriter.WriteAll(rows)
}
