package config

import (
	"fmt"
	"time"

	wbconfig "github.com/wb-go/wbf/config"
)

type Config struct {
	App   App      `mapstructure:"app"`
	DB    Database `mapstructure:"database"`
	Retry Retry    `mapstructure:"retry"`
}

type App struct {
	Port            string        `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	MigrationDir    string        `mapstructure:"migration_dir"`
}

type Database struct {
	URL             string        `mapstructure:"url"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type Retry struct {
	Attempts int           `mapstructure:"attempts"`
	Delay    time.Duration `mapstructure:"delay"`
	Backoff  float64       `mapstructure:"backoff"`
}

func Load(configFilePath string) (*Config, error) {
	c := wbconfig.New()

	c.EnableEnv("")

	err := c.LoadConfigFiles(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	cfg := new(Config)
	if err := c.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
