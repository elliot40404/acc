package database

import (
	_ "embed"
	"errors"
	"log/slog"

	"github.com/elliot40404/acc/pkg/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

var DBPATH = utils.DBPATH()

func GetDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", DBPATH)
	if err != nil {
		slog.Error("DB: failed to open database", "Error", err.Error())
		return nil, errors.New("failed to open database")
	}
	return db, nil
}

func InitApplication() error {
	path := utils.DBPATH()
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		slog.Error("DB: failed to open database", "Error", err.Error())
		return errors.New("failed to open database")
	}
	_, err = db.Exec(schema)
	if err != nil {
		slog.Error("DB: failed to initialize database schema", "Error", err.Error())
		return errors.New("failed to initialize database schema")
	}
	return nil
}