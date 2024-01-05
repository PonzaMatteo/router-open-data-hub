package mapper

import (
	"encoding/json"
)

type Mapper struct {
	mapping map[string]string
}

func EmptyMapper() Mapper {
	return NewMapper(make(map[string]string))
}

func NewMapper(mapping map[string]string) Mapper {
	return Mapper{
		mapping: mapping,
	}
}

func (m *Mapper) Transform(input string) (string, error) {
	if len(m.mapping) == 0 {
		return input, nil
	}

	var inputResponse interface{}

	err := json.Unmarshal([]byte(input), &inputResponse)
	if err != nil {
		return "", err
	}

	outputResponse := m.extractMapping(inputResponse)

	modifiedJSON, err := json.Marshal(outputResponse)
	if err != nil {
		return "", err
	}

	return string(modifiedJSON), nil
}

func (m Mapper) extractMapping(inputResponse interface{}) interface{} {
	switch inputResponse := inputResponse.(type) {
	case map[string]interface{}:
		outputResponse := make(map[string]interface{})
		for inputKey, value := range inputResponse {
			if outputKey, ok := m.mapping[inputKey]; ok {
				outputResponse[outputKey] = value
			}
		}
		return outputResponse
	case []interface{}:
		var outputResponse []interface{}
		for _, data := range inputResponse {
			outputResponse = append(outputResponse, m.extractMapping(data))
		}
		return outputResponse
	default:
		return inputResponse
	}
}

func (m *Mapper) AddMapping(inputKey string, outputKey string) {
	m.mapping[inputKey] = outputKey
}
