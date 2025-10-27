package config 

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerAddress string `json:"server_address" yaml:"server_address"`
	MaxWorkers    int    `json:"max_workers" yaml:"max_workers"`
	LogLevel      string `json:"log_level" yaml:"log_level"`
	TimeoutSeconds int `json:"timeout_seconds" yaml:"timeout_seconds"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &Config{}

	switch ext := filepath.Ext(path); ext {
	case ".json":
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse JSON: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse YAML: %w", err)
		}
	default:
		return nil, errors.New("unsupported config format: must be .json or .yaml")
	}

	return cfg, nil
}