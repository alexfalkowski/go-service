package config

import (
	"errors"
	"os"
)

// ErrMissingConfigFile for config.
var ErrMissingConfigFile = errors.New("missing config file")

// File config location.
func File() (string, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		return "", ErrMissingConfigFile
	}

	return configFile, nil
}

// ReadFile from config location.
func ReadFile() ([]byte, error) {
	configFile, err := File()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(configFile)
}
