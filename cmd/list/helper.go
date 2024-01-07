package list

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/elliot40404/acc/pkg/database"
	"github.com/elliot40404/acc/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

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

func tableWriter(transactions []database.Transaction, totalTx int, queryConfig database.TransactionConfig, interactive bool) string {
	t := table.NewWriter()
	header := getTableHeader(queryConfig)
	t.AppendHeader(header)
	rows := getTableRows(transactions, queryConfig)
	t.AppendRows(rows)
	t.SetStyle(table.StyleLight)
	if !interactive {
		t.SetOutputMirror(os.Stdout)
		t.Render()
		printTableSummary(queryConfig, totalTx, len(transactions), false)
		return ""
	}
	result := t.Render()
	result += printTableSummary(queryConfig, totalTx, len(transactions), true)
	return result
}

func getTableHeader(queryConfig database.TransactionConfig) table.Row {
	if len(queryConfig.Columns) == 0 {
		return table.Row{
			"#",
			"Type",
			"Amt",
			"Desc",
			"Date",
		}
	}
	header := table.Row{}
	for _, column := range queryConfig.Columns {
		for _, validColumn := range validColumns {
			if column == validColumn {
				header = append(header, column)
			}
		}
	}
	return header
}

func getTableRows(transactions []database.Transaction, queryConfig database.TransactionConfig) []table.Row {
	rows := []table.Row{}
	for _, transaction := range transactions {
		if queryConfig.IsHRTime {
			transaction.CreatedAt = utils.HRTime(transaction.CreatedAt)
		}
		row := getTableRow(transaction, queryConfig)
		rows = append(rows, row)
	}
	return rows
}

func getTableRow(transaction database.Transaction, queryConfig database.TransactionConfig) table.Row {
	if len(queryConfig.Columns) == 0 {
		return table.Row{
			strconv.Itoa(transaction.ID),
			transaction.Type,
			transaction.Amount,
			transaction.Description,
			transaction.CreatedAt,
		}
	}
	row := table.Row{}
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
	return row
}

func printTableSummary(queryConfig database.TransactionConfig, totalTx int, transactions int, interactive bool) string {
	totalPages := int(math.Ceil(float64(totalTx) / float64(queryConfig.Limit)))
	if queryConfig.All {
		totalPages = 1
	}
	if !interactive {
		fmt.Printf(
			"Page: %d of %d | Results: %d | Total: %d\n",
			queryConfig.Page,
			totalPages,
			transactions,
			totalTx,
		)
		return ""
	}
	return fmt.Sprintf(
		"\nPage: %d of %d | Results: %d | Total: %d\n",
		queryConfig.Page,
		totalPages,
		transactions,
		totalTx,
	)
}
