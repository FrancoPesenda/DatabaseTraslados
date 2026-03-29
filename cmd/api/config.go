package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type InfraConfig struct {
	HTTP  HTTPConfig  `json:"http"`
	MySQL MySQLConfig `json:"mysql"`
}

type HTTPConfig struct {
	Addr string `json:"addr"`
}

type MySQLConfig struct {
	DSN             string `json:"dsn"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime string `json:"conn_max_lifetime"`
}

type LoadConfigOptions struct {
	ConfigDir string
	Scope     string
}

func LoadInfraConfig(opts LoadConfigOptions) (InfraConfig, error) {
	if opts.ConfigDir == "" {
		opts.ConfigDir = envOr("CONFIG_DIR", "config")
	}
	if opts.Scope == "" {
		opts.Scope = envOr("CONFIG_SCOPE", "local")
	}

	p := filepath.Join(opts.ConfigDir, opts.Scope, "infrastructure_config.json")
	var cfg InfraConfig
	if b, err := os.ReadFile(p); err == nil {
		if err := json.Unmarshal(b, &cfg); err != nil {
			return InfraConfig{}, fmt.Errorf("parse config file %q: %w", p, err)
		}
	} else {
		if !errors.Is(err, os.ErrNotExist) {
			return InfraConfig{}, fmt.Errorf("read config file %q: %w", p, err)
		}
	}

	// Environment variables override file config.
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		cfg.HTTP.Addr = v
	}
	if v := os.Getenv("MYSQL_DSN"); v != "" {
		cfg.MySQL.DSN = v
	}
	if v := os.Getenv("MYSQL_MAX_OPEN_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.MySQL.MaxOpenConns = n
		}
	}
	if v := os.Getenv("MYSQL_MAX_IDLE_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.MySQL.MaxIdleConns = n
		}
	}
	if v := os.Getenv("MYSQL_CONN_MAX_LIFETIME"); v != "" {
		cfg.MySQL.ConnMaxLifetime = v
	}

	if cfg.HTTP.Addr == "" {
		cfg.HTTP.Addr = ":8080"
	}
	if cfg.MySQL.DSN == "" {
		return InfraConfig{}, fmt.Errorf("missing MYSQL_DSN (set it in environment or config file %q)", p)
	}
	if cfg.MySQL.MaxOpenConns == 0 {
		cfg.MySQL.MaxOpenConns = 25
	}
	if cfg.MySQL.MaxIdleConns == 0 {
		cfg.MySQL.MaxIdleConns = 25
	}
	if cfg.MySQL.ConnMaxLifetime == "" {
		cfg.MySQL.ConnMaxLifetime = "5m"
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
