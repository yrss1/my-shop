package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	defaultAppMode    = "dev"
	defaultAppPort    = "8080"
	defaultAppPath    = "/"
	defaultAppTimeout = 60 * time.Second
)

type (
	Configs struct {
		APP AppConfig
		API ApiConfig
	}

	AppConfig struct {
		Mode    string `required:"true"`
		Port    string
		Path    string
		Timeout time.Duration
	}

	ApiConfig struct {
		Order   string
		Product string
		Payment string
		User    string
	}
)

func New() (cfg Configs, err error) {
	root, err := os.Getwd()
	if err != nil {
		return
	}
	godotenv.Load(filepath.Join(root, ".env"))

	cfg.APP = AppConfig{
		Mode:    defaultAppMode,
		Port:    defaultAppPort,
		Path:    defaultAppPath,
		Timeout: defaultAppTimeout,
	}

	if err = envconfig.Process("APP", &cfg.APP); err != nil {
		return
	}

	if err = envconfig.Process("API", &cfg.API); err != nil {
		return
	}

	return
}
