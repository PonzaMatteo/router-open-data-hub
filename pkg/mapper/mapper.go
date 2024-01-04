package mapper

import "encoding/json"

type Mapper struct {
	mapping map[string]string
}

func NewMapper() Mapper {
	return Mapper{
		mapping: make(map[string]string),
	}
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
