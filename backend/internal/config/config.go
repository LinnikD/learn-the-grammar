// Package config loads backend configuration from defaults, an optional
// YAML file, and environment variables, in that ascending order of priority.
package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config holds all backend configuration values.
type Config struct {
	Port int `yaml:"port"`
}

func defaultConfig() Config {
	return Config{
		Port: 8080,
	}
}

// Load builds the configuration by starting from defaults, applying values
// from the YAML file at path (if path is non-empty), and finally applying
// overrides from environment variables.
//
// A non-empty path that cannot be read, or a file that is not valid YAML,
// results in an error. Fields absent from the file keep their default
// value; fields present in the file but not recognized by Config are
// ignored.
func Load(path string) (Config, error) {
	cfg := defaultConfig()

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return Config{}, fmt.Errorf("reading config file %q: %w", path, err)
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parsing config file %q: %w", path, err)
		}
	}

	if v, ok := os.LookupEnv("LTG_PORT"); ok {
		port, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("parsing LTG_PORT %q: %w", v, err)
		}
		cfg.Port = port
	}

	return cfg, nil
}
