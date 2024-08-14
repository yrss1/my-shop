package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	defaultAppMode     = "dev"
	defaultAppPort     = "8004"
	defaultAppGRPCPort = "9004"
	defaultAppPath     = "/"
	defaultAppTimeout  = 60 * time.Second
)

type (
	Configs struct {
		APP      AppConfig
		POSTGRES StoreConfig
		EPAY     CredentialsConfig
	}

	AppConfig struct {
		Mode     string `required:"true"`
		Port     string
		GRPCPort string
		Path     string
		Timeout  time.Duration
	}

	StoreConfig struct {
		DSN string
	}
	CredentialsConfig struct {
		URL      string
		Login    string
		Password string

		OAuthURL       string
		PaymentPageURL string
	}
)

func New() (cfg Configs, err error) {
	root, err := os.Getwd()
	if err != nil {
		return
	}
	godotenv.Load(filepath.Join(root, ".env"))

	cfg.APP = AppConfig{
		Mode:     defaultAppMode,
		Port:     defaultAppPort,
		GRPCPort: defaultAppGRPCPort,
		Path:     defaultAppPath,
		Timeout:  defaultAppTimeout,
	}

	if err = envconfig.Process("APP", &cfg.APP); err != nil {
		return
	}

	if err = envconfig.Process("POSTGRES", &cfg.POSTGRES); err != nil {
		return
	}

	if err = envconfig.Process("EPAY", &cfg.EPAY); err != nil {
		return
	}

	return
}
