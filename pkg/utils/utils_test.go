package utils_test

import (
	"testing"

	"github.com/elliot40404/acc/pkg/utils"
)

func TestSplitDateRange(t *testing.T) {
	// test >date
	r := utils.SplitDateRange(":2020-01-01")
	if len(r) != 1 {
		t.Error("expected 1, got", len(r))
	}
	// test <date
	r = utils.SplitDateRange("2020-01-01:")
	if len(r) != 1 {
		t.Error("expected 1, got", len(r))
	}
	// test date<>date
	r = utils.SplitDateRange("2020-01-01:2020-01-01")
	if len(r) != 2 {
		t.Error("expected 2, got", len(r))
	}
}

func TestPadDate(t *testing.T) {
	// test 1
	r := utils.PadDate("1")
	if r != "01" {
		t.Error("expected 01, got", r)
	}
	// test 01
	r = utils.PadDate("01")
	if r != "01" {
		t.Error("expected 01, got", r)
	}
	// test 10
	r = utils.PadDate("10")
	if r != "10" {
		t.Error("expected 10, got", r)
	}
	// test 11
	r = utils.PadDate("11")
	if r != "11" {
		t.Error("expected 11, got", r)
	}
	r = utils.PadDate("2020-1-1")
	if r != "2020-01-01" {
		t.Error("expected 2020-01-01, got", r)
	}
}