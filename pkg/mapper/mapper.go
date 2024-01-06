package mapper

import (
	"encoding/json"
	"strings"
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
			switch value := value.(type) {
			case map[string]interface{}:
				outputKey, extractedValues := m.extractInternalMapping(inputKey, value)
				outputResponse[outputKey] = extractedValues
			case []interface{}:
				var extractedValues []interface{}
				for _, v := range value {
					outputKey, extractedValue := m.extractInternalMapping(inputKey, v)
					extractedValues = append(extractedValues, extractedValue)
					outputResponse[outputKey] = extractedValues
				}
			default:	
				if outputKey, ok := m.mapping[inputKey]; ok {
					outputResponse[outputKey] = value
				}
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

func (m Mapper) extractInternalMapping(previousInputKey string, inputResponse interface{}) (string, interface{}) {
	switch inputResponse := inputResponse.(type) {
	case map[string]interface{}:

		outputResponse := make(map[string]interface{})
		previousOutputKey := ""
		
		for inputKey, value := range inputResponse {
			newInputKey := previousInputKey + "." + inputKey
			if outputKey, ok := m.mapping[newInputKey]; ok {
				newOutputKey := strings.Split(outputKey, ".")[1]
				previousOutputKey = strings.Split(outputKey, ".")[0]
				outputResponse[newOutputKey] = value
			}
		}
		
		return previousOutputKey, outputResponse
	// case []interface{}:
	// 	var outputResponse []interface{}
	// 	for _, data := range inputResponse {
	// 		outputResponse = append(outputResponse, m.extractMapping(data))
	// 	}
	// 	return outputResponse
	default:
		return "", inputResponse
	}
}

func (m *Mapper) AddMapping(inputKey string, outputKey string) {
	m.mapping[inputKey] = outputKey
}
