package db

import (
	"database/sql"
	"log/slog"

	"github.com/go-sql-driver/mysql"
)

func NewMysqlStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return db, nil
}
