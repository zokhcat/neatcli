package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	AppName    = "neatcli"
	ConfigDir  = ".neatlogs"
	ConfigFile = "config.yaml"
	EnvAPIKey  = "NEATLOGS_API_KEY"
)

type Config struct {
	APIKey  string `yaml:"api_key" mapstructure:"api_key"`
	BaseURL string `yaml:"base_url" mapstructure:"base_url"`
	Project string `yaml:"project" mapstructure:"project"`
}

func DefaultConfig() Config {
	return Config{
		BaseURL: "https://ingest.neatlogs.com",
	}
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	return filepath.Join(home, ConfigDir, ConfigFile), nil
}

func WorkspaceDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ConfigDir), nil
}

func Load() (Config, error) {
	cfg := DefaultConfig()

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}

	if envKey := os.Getenv(EnvAPIKey); envKey != "" {
		cfg.APIKey = envKey
	}

	return cfg, nil
}

func Save(cfg Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ConfigDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, ConfigFile)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(cfg)
}

func InitViper() error {
	cfgPath, err := ConfigPath()
	if err != nil {
		return err
	}

	viper.SetConfigFile(cfgPath)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("NEATLOGS")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}
