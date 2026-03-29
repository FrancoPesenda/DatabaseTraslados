package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type InfraConfig struct {
	HTTP HTTPConfig  `json:"http"`
	MySQL MySQLConfig `json:"mysql"`
}

type HTTPConfig struct {
	Addr string `json:"addr"`
}

type MySQLConfig struct {
	DSN string `json:"dsn"`
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
	b, err := os.ReadFile(p)
	if err != nil {
		return InfraConfig{}, fmt.Errorf("read config file %q: %w", p, err)
	}

	var cfg InfraConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		return InfraConfig{}, fmt.Errorf("parse config file %q: %w", p, err)
	}

	if cfg.HTTP.Addr == "" {
		cfg.HTTP.Addr = ":8080"
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

