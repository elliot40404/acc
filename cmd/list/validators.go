package list

import (
	"errors"

	"github.com/elliot40404/acc/pkg/database"
	"github.com/elliot40404/acc/pkg/utils"
)

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

func ValidateConfig(config *database.TransactionConfig) error {
	if err := validateTxType(config.TxType); err != nil {
		return err
	}
	if err := validatePageAndLimit(config.Page, config.Limit); err != nil {
		return err
	}
	if err := validateDate(config.Date); err != nil {
		return err
	}
	if err := validateAmount(config.Amount); err != nil {
		return err
	}
	if err := validateSort(config.Sort); err != nil {
		return err
	}
	return nil
}

func validateTxType(txType string) error {
	if txType != "" && txType != "income" && txType != "expense" {
		return errors.New("invalid type. type must be either income or expense")
	}
	return nil
}

func validatePageAndLimit(page, limit int) error {
	if page < 1 || limit < 1 {
		return errors.New("invalid page or limit. page and limit must be greater than 0")
	}
	return nil
}

func validateDate(date string) error {
	if date != "" {
		err := utils.Checkdate(&date)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateAmount(amount string) error {
	if amount != "" {
		err := utils.CheckAmount(&amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateSort(sort string) error {
	if sort != "" {
		for _, sortable := range sortables {
			if sort == sortable {
				return nil
			}
		}
		return errors.New("invalid sort. sort must be either 'date' or 'amt'")
	}
	return nil
}