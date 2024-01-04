package mapper

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

type Mapper struct {
	mapping map[string]string
}

type Route struct {
	Keyword string
	Service string
	Mapping *map[string]string
}

type Config struct {
	Routes []Route
}

func NewMapper() Mapper {
	return Mapper{
		mapping: make(map[string]string),
	}
}

func NewMapperFromFile(fileName string, keyword string) (Mapper, error) {
	var config, err = readConfigFromFile(fileName)
	if err != nil {
		return Mapper{}, err
	}
	for _, route := range config.Routes {
		if (route.Keyword == keyword) {
			var mapper = Mapper{
				mapping: *route.Mapping,
			}
			return mapper, nil
		}
	}
	return NewMapper(), nil
}

func (m *Mapper) Transform(input string) (string, error) {

	var inputResponse map[string]interface{}
	var outputResponse = make(map[string]interface{})

	err := json.Unmarshal([]byte(input), &inputResponse)
	if err != nil {
		return "", err
	}

	for inputKey, value := range inputResponse {
		if outputKey, ok := m.mapping[inputKey]; ok {
			outputResponse[outputKey] = value
		}
	}

	modifiedJSON, err := json.Marshal(outputResponse)
	if err != nil {
		return "", err
	}

	return string(modifiedJSON), nil
}

func (m *Mapper) AddMapping(inputKey string, outputKey string) {
	m.mapping[inputKey] = outputKey
}

func readConfigFromFile(fileName string) (*Config, error) {
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
