package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

type Route struct {
	Keyword string
	Service string
	Mapping *map[string]string
}

type Config struct {
	Routes []Route
}

func FromFile(fileName string) (*Config, error) {
	var configData Config
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	extension := strings.ToLower(path.Ext(fileName))
	if extension != ".json" {
		return nil, fmt.Errorf("unsupported configuration file extension (`%s`): %s", extension, fileName)
	}

	err = json.Unmarshal([]byte(data), &configData)
	if err != nil {
		return nil, err
	}
	return &configData, nil
}

var defaultConfig *Config

func GetDefault() *Config {
	if defaultConfig != nil {
		return defaultConfig
	}

	const defaultConfigFile = "config.json"
	c, err := FromFile(defaultConfigFile)
	if err != nil {
		panic(fmt.Errorf("default configuration is not valid (%s): %w", defaultConfigFile, err))
	}
	defaultConfig = c
	return defaultConfig
}
