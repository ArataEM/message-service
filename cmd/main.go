package main

import (
	"log/slog"

	"github.com/ArataEM/message-service/internal/config"
	"github.com/ArataEM/message-service/internal/db"
	"github.com/ArataEM/message-service/internal/server"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMysqlStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		slog.Error(err.Error())
		panic(1)
	}

	server := server.NewServer(":8080", db)
	if err := server.Run(); err != nil {
		slog.Error(err.Error())
	}
}
