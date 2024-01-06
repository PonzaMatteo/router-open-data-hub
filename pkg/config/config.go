package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type Route struct {
	Keyword string
	Service string
	Mapping *map[string]string
}

type Config struct {
	Routes []Route
}

type decoder interface {
	Decode(v any) error
}

func FromFile(fileName string) (*Config, error) {
	var configData Config
	data, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	var d decoder
	switch extension := strings.ToLower(path.Ext(fileName)); extension {
	case ".yaml", ".yml":
		d = yaml.NewDecoder(data)
	case ".json":
		d = json.NewDecoder(data)
	default:
		return nil, fmt.Errorf("unsupported configuration file extension (`%s`): %s", extension, fileName)
	}

	if err := d.Decode(&configData); err != nil {
		return nil, err
	}
	return &configData, nil
}

var defaultConfig *Config

func GetDefault() *Config {
	if defaultConfig != nil {
		return defaultConfig
	}

	const defaultConfigFile = "config.yaml"
	c, err := FromFile(defaultConfigFile)
	if err != nil {
		panic(fmt.Errorf("default configuration is not valid (%s): %w", defaultConfigFile, err))
	}
	defaultConfig = c
	return defaultConfig
}
