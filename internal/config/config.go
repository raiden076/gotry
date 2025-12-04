package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Workspace WorkspaceConfig `mapstructure:"workspace"`
	Git       GitConfig       `mapstructure:"git"`
}

type WorkspaceConfig struct {
	Path string `mapstructure:"path"`
}

type GitConfig struct {
	AutoInit      bool `mapstructure:"auto_init"`
	InitialCommit bool `mapstructure:"initial_commit"`
}

func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		Workspace: WorkspaceConfig{
			Path: filepath.Join(homeDir, "tries"),
		},
		Git: GitConfig{
			AutoInit:      true,
			InitialCommit: true,
		},
	}
}

func Load() (*Config, error) {
	cfg := DefaultConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil // Return defaults if home dir unavailable
	}

	configDir := filepath.Join(homeDir, ".config", "gotry")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return cfg, nil // Config file not found, use defaults
		}
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Expand ~ in workspace path
	if cfg.Workspace.Path[:2] == "~/" {
		cfg.Workspace.Path = filepath.Join(homeDir, cfg.Workspace.Path[2:])
	}

	return cfg, nil
}

func (c *Config) EnsureWorkspaceExists() error {
	return os.MkdirAll(c.Workspace.Path, 0755)
}
