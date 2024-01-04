package database

import (
	"fmt"

	"github.com/elliot40404/acc/pkg/utils"
)

var colMap = map[string]string{
	"id":   "id",
	"type": "type",
	"amt":  "amount",
	"desc": "description",
	"date": "created_at",
}

func (q *TQuery) Build() string {
	return q.Query
}

func (q *TQuery) AddColumns() {
	if len(q.Config.Columns) != 0 && !q.isCount {
		q.Query = " SELECT "
		var cols string
		for i, col := range q.Config.Columns {
			if i != len(q.Config.Columns)-1 {
				cols += colMap[col] + ", "
			} else {
				cols += colMap[col]
			}
		}
		if cols != "" {
			q.Query += cols
		} else {
			q.Query += "*"
		}
		q.Query += " FROM transactions"
	}
}

func (q *TQuery) AddType() {
	if q.Config.TxType != "" {
		q.Query += fmt.Sprintf(" WHERE type = '%s'", q.Config.TxType)
	}
}

func (q *TQuery) AddDate() {
	if q.Config.Date != "" {
		var dateQuery string
		if utils.IsValueRange(q.Config.Date) {
			dateQuery += buildDateRangeQuery(q.Config.Date)
		} else {
			dateQuery += buildDateQuery(q.Config.Date)
		}
		if q.Config.TxType != "" && dateQuery != "" {
			dateQuery = " AND" + dateQuery
		} else if dateQuery != "" {
			dateQuery = " WHERE" + dateQuery
		}
		q.Query += dateQuery
	}
}

func (q *TQuery) AddAmount() {
	if q.Config.Amount != "" {
		var amountQuery string
		if utils.IsValueRange(q.Config.Amount) {
			amountQuery += buildAmountRangeQuery(q.Config.Amount)
		} else {
			amountQuery += buildAmountQuery(q.Config.Amount)
		}
		if q.Config.TxType != "" && amountQuery != "" {
			amountQuery = " AND" + amountQuery
		} else if amountQuery != "" {
			amountQuery = " WHERE" + amountQuery
		}
		q.Query += amountQuery
	}
}

func (q *TQuery) AddDesc() {
	if q.Config.Desc != "" {
		var descQuery string
		if q.Config.TxType != "" || q.Config.Date != "" || q.Config.Amount != "" {
			descQuery += " AND"
		} else {
			descQuery += " WHERE"
		}
		descQuery += fmt.Sprintf(" description LIKE '%%%s%%'", q.Config.Desc)
		q.Query += descQuery
	}
}

func (q *TQuery) AddLimit() {
	if q.Config.Limit != 0 {
		q.Query += fmt.Sprintf(" LIMIT %d", q.Config.Limit)
	}
}

func (q *TQuery) AddOffset() {
	if q.Config.Page != 0 {
		q.Query += fmt.Sprintf(" OFFSET %d", (q.Config.Page-1)*q.Config.Limit)
	}
}

func (q *TQuery) AddSort() {
	if q.Config.Sort != "" {
		if q.Config.Sort == "date" {
			q.Config.Sort = "created_at"
		} else if q.Config.Sort == "amt" {
			q.Config.Sort = "amount"
		}
		q.Query += fmt.Sprintf(" ORDER BY %s", q.Config.Sort)
		if q.Config.SortAsc {
			q.Query += " ASC"
		} else {
			q.Query += " DESC"
		}
	}
}

func buildDateQuery(date string) string {
	switch date {
	case "today":
		return " DATE(created_at) = DATE('now')"
	case "yesterday":
		return " DATE(created_at) = DATE('now', '-1 day')"
	case "thisweek":
		return " strftime('%W', created_at) = strftime('%W', 'now')"
	case "lastweek":
		return " strftime('%W', created_at) = strftime('%W', 'now', '-7 days')"
	case "thismonth":
		return " strftime('%m', created_at) = strftime('%m', 'now')"
	case "lastmonth":
		return " strftime('%m', created_at) = strftime('%m', 'now', '-1 month')"
	case "thisyear":
		return " strftime('%Y', created_at) = strftime('%Y', 'now')"
	case "lastyear":
		return " strftime('%Y', created_at) = strftime('%Y', 'now', '-1 year')"
	default:
		return " DATE(created_at) = DATE('" + utils.ConvertToDateFormat(date) + "')"
	}
}

func buildDateRangeQuery(date string) string {
	if date[0] == ':' {
		return " DATE(created_at) <= DATE('" + utils.ConvertToDateFormat(date[1:]) + "')"
	} else if date[len(date)-1] == ':' {
		return " DATE(created_at) >= DATE('" + utils.ConvertToDateFormat(date[:len(date)-1]) + "')"
	} else {
		dates := utils.SplitDateRange(date)
		// NOTE: can throw a warning if dates[0] > dates[1]
		return " DATE(created_at) BETWEEN DATE('" + dates[0] + "') AND DATE('" + dates[1] + "')"
	}
}

func buildAmountQuery(amount string) string {
	return " amount = " + amount
}

func buildAmountRangeQuery(amount string) string {
	if amount[0] == ':' {
		return " amount <= " + amount[1:]
	} else if amount[len(amount)-1] == ':' {
		return " amount >= " + amount[:len(amount)-1]
	} else {
		amounts := utils.SplitAmountRange(amount)
		// NOTE: can throw a warning if amounts[0] > amounts[1]
		return " amount BETWEEN " + amounts[0] + " AND " + amounts[1]
	}
}
