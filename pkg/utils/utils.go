package utils

import (
	"errors"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/itlightning/dateparse"
)

const (
	DateSyntaxError        = "invalid date syntax. date must be in a valid format or one of the predefined builtins like today, yesterday"
	DateRangeSyntaxError   = "invalid date range synxtax. Must be one of the following: :date, date:, date:date"
	AmountSyntaxError      = "invalid amount syntax. amount must be in a valid format"
	AmountRangeSyntaxError = "invalid amount range syntax. must be one of the following: :amount, amount:, amount:amount"
)

var ValidDateBuiltins = []string{
	"today",
	"yesterday",
	"thisweek",
	"lastweek",
	"thismonth",
	"lastmonth",
	"thisyear",
	"lastyear",
}

func PrintError(err error, message string, debug bool) {
	slog.Error(message)
	if debug {
		slog.Debug(err.Error())
	}
}

func IsInitialized() bool {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	if _, err := os.Stat(homedir + "/.acc/.initialized"); err == nil {
		return true
	}
	return false
}

func DBPATH() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return homedir + "/.acc/acc.db"
}

func APPDIR() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return homedir + "/.acc"
}

// conver time to human readable time
func HRTime(date string) string {
	t, err := dateparse.ParseAny(date)
	if err != nil {
		return date
	}
	return t.Format("02 Jan 2006 15:04:05")
}

func PadDate(date string) string {
	re := regexp.MustCompile(`(\b\d\b)`)
	date = re.ReplaceAllString(date, "0$1")
	return date
}

func SplitDateRange(dateRange string) []string {
	var filtered []string
	for _, date := range strings.Split(dateRange, ":") {
		if date != "" {
			date = strings.Replace(date, "/", "-", -1)
			filtered = append(filtered, ConvertToDateFormat(date))
		}
	}
	return filtered
}

func ToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

func Checkdate(dates ...*string) error {
	for _, date := range dates {
		if IsValueRange(*date) {
			err := CheckdateRange(date)
			if err != nil {
				return err
			}
			continue
		}
		// if date contains / then replace it with -
		if strings.Contains(*date, "/") {
			*date = strings.ReplaceAll(*date, "/", "-")
		}
		// replace all single digits with 0 prefixed double digits with regex
		*date = PadDate(*date)
		if !IsValidDateFormat(*date) {
			return errors.New(DateSyntaxError)
		}
	}
	return nil
}

func CheckdateRange(date *string) error {
	if strings.HasPrefix(*date, ":") || strings.HasSuffix(*date, ":") {
		dateRange := SplitDateRange(*date)
		if len(dateRange) != 1 {
			return errors.New(DateRangeSyntaxError)
		}
		err := Checkdate(&dateRange[0])
		if err != nil {
			return err
		}
	} else if strings.Contains(*date, ":") {
		dateRange := SplitDateRange(*date)
		if len(dateRange) != 2 {
			return errors.New(DateRangeSyntaxError)
		}
		err := Checkdate(&dateRange[0], &dateRange[1])
		if err != nil {
			return err
		}
	} else {
		return errors.New(DateRangeSyntaxError)
	}
	return nil
}

func IsValidDateFormat(date string) bool {
	// check if date in builtins
	for _, builtin := range ValidDateBuiltins {
		if date == builtin {
			return true
		}
	}
	_, err := dateparse.ParseAny(date)
	return err == nil
}

func ConvertToDateFormat(date string) string {
	t, _ := dateparse.ParseAny(date)
	// convert to YYYY-MM-DD format
	return t.Format("2006-01-02")
}

func IsValueRange(value string) bool {
	return strings.Contains(value, ":")
}

func CheckAmount(amounts ...*string) error {
	for _, amount := range amounts {
		if IsValueRange(*amount) {
			err := CheckAmountRange(amount)
			if err != nil {
				return err
			}
			continue
		}
		if !IsValidAmountFormat(*amount) {
			return errors.New(AmountSyntaxError)
		}
	}
	return nil
}

func CheckAmountRange(amount *string) error {
	if strings.HasPrefix(*amount, ":") || strings.HasSuffix(*amount, ":") {
		amountRange := SplitAmountRange(*amount)
		if len(amountRange) != 1 {
			return errors.New(AmountRangeSyntaxError)
		}
		err := CheckAmount(&amountRange[0])
		if err != nil {
			return err
		}
	} else if strings.Contains(*amount, ":") {
		amountRange := SplitAmountRange(*amount)
		if len(amountRange) != 2 {
			return errors.New(AmountRangeSyntaxError)
		}
		err := CheckAmount(&amountRange[0], &amountRange[1])
		if err != nil {
			return err
		}
	} else {
		return errors.New(AmountRangeSyntaxError)
	}
	return nil
}

func IsValidAmountFormat(amount string) bool {
	// parse amount to float64
	_, err := strconv.ParseFloat(amount, 64)
	return err == nil
}

func SplitAmountRange(amountRange string) []string {
	var filtered []string
	for _, amount := range strings.Split(amountRange, ":") {
		if amount != "" {
			filtered = append(filtered, amount)
		}
	}
	return filtered
}
