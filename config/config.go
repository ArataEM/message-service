package config

import (
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	RedisAddress string
	ServerPort   uint16
}

func LoadConfig() Config {
	port, err := strconv.ParseUint(getEnv("SERVER_PORT", "8080"), 10, 16)
	if err != nil {
		slog.Error("failed to parse env SERVER_PORT")
		panic(1)
	}
	cfg := Config{
		RedisAddress: getEnv("REDIS_ADDR", "localhost:6379"),
		ServerPort:   uint16(port),
	}
	return cfg
}

func getEnv(name, fallback string) string {
	env, exists := os.LookupEnv(name)
	if exists {
		return env
	}
	return fallback
}
