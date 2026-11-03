package config

import (
	"ValorantRPC/pkg/errors"
	"encoding/json"
	"fmt"
	"os"
)

type LoggingConfig struct {
	Level      string `json:"level"`
	FilePath   string `json:"file_path"`
	JSONFormat bool   `json:"json_format"`
}

type I18nConfig struct {
	Language string `json:"language"`
}

type AppConfig struct {
	PollingInterval int           `json:"refresh_interval"`
	Logging         LoggingConfig `json:"logging"`
	Language        string        `json:"language"`
}

func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		PollingInterval: 15,
		Logging: LoggingConfig{
			Level:      "info",
			FilePath:   "",
			JSONFormat: false,
		},
		Language: "en",
	}
}

func Load(configPath string) (*AppConfig, error) {
	cfg := DefaultAppConfig()

	if configPath == "" {
		configPath = "config.json"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		cause := fmt.Errorf("%w: %v", errors.ErrConfigFileUnknown, err)
		return nil, errors.New("config.load", cause)
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		cause := fmt.Errorf("%w: %v", errors.ErrConfigFileMalformed, err)
		return nil, errors.New("config.Load", cause)
	}

	return cfg, nil
}

func Save(configPath string, cfg *AppConfig) error {
	if configPath == "" {
		configPath = "config.json"
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		cause := fmt.Errorf("%w: %v", errors.ErrConfigFileMalformed, err)
		return errors.New("config.Save", cause)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return errors.New("config.Save", err)
	}

	return nil
}
