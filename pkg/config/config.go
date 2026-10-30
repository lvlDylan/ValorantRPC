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
	RefreshInterval int           `json:"refresh_interval"`
	Logging         LoggingConfig `json:"logging"`
	I18n            I18nConfig    `json:"i18n"`
}

func defaultAppConfig() *AppConfig {
	return &AppConfig{
		RefreshInterval: 15,
		Logging: LoggingConfig{
			Level:      "info",
			FilePath:   "",
			JSONFormat: false,
		},
		I18n: I18nConfig{
			Language: "en",
		},
	}
}

func Load(configPath string) (*AppConfig, error) {
	cfg := defaultAppConfig()

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
