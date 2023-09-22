package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// EnvPrefix is the prefix to use for environment variables.
const EnvPrefix = "BREWDAY_"

func LoadConfig(path string) (*Config, error) {
	k := koanf.New(".") // New creates a new instance of Koanf.
	if path != "" {
		ext := filepath.Ext(path)
		parser, err := getParser(ext)
		if err != nil {
			return nil, err
		}
		err = k.Load(file.Provider(path), parser)
		if err != nil {
			return nil, err
		}
	}
	err := k.Load(env.ProviderWithValue(EnvPrefix, ".", formatEnvVariables), nil)
	if err != nil {
		return nil, err
	}
	var config Config
	err = k.Unmarshal("", &config)
	if err != nil {
		return nil, err
	}
	err = validateConfig(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// getParser returns a parser for the given file given the format.
// It returns an error if the format is not supported.
// Supported formats are: .yaml, .yml.
// Support can be extended to .env, .toml, .hcl, .ini, .json
func getParser(format string) (koanf.Parser, error) {
	switch format {
	case ".yaml", ".yml":
		return yaml.Parser(), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// formatEnvVariables formats the environment variables to help with parsing.
// It removes the prefix from the variable name and replaces _ with .
func formatEnvVariables(s string, v string) (string, interface{}) {
	key := strings.Replace(strings.ToLower(strings.TrimPrefix(s, EnvPrefix)), "_", ".", -1)
	return key, v
}

// validateConfig validates the config.
// It returns an error if the config is invalid.
func validateConfig(config *Config) error {
	if config.App.Port == 0 {
		return fmt.Errorf("port is missing")
	}
	if config.Notification.Enabled {
		if config.Notification.AppToken == "" {
			return fmt.Errorf("notification is enabled but app-token is missing")
		}
		if config.Notification.GotifyURL == "" {
			return fmt.Errorf("notification is enabled but gotify-url is missing")
		}
	}
	return nil
}
