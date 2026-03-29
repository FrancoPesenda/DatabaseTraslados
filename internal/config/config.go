package config

import (
	"os"
)

type Config struct {
	HTTP  HTTPConfig
	MySQL MySQLConfig
}

type HTTPConfig struct {
	Addr string
}

type MySQLConfig struct {
	DSN string
}

func FromEnv() Config {
	return Config{
		HTTP: HTTPConfig{
			Addr: envOr("HTTP_ADDR", ":8080"),
		},
		MySQL: MySQLConfig{
			DSN: envOr("MYSQL_DSN", "root:password@tcp(127.0.0.1:3306)/app?parseTime=true"),
		},
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

