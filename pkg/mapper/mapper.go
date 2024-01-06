package mapper

import (
	"bytes"
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

func (m *Mapper) AddMapping(inputKey string, outputKey string) {
	m.mapping[inputKey] = outputKey
}

func (m *Mapper) Transform(input string) (string, error) {
	if len(m.mapping) == 0 {
		return input, nil
	}

	var inputResponse interface{}
	if err := json.NewDecoder(strings.NewReader(input)).Decode(&inputResponse); err != nil {
		return "", err

	}

	_, outputResponse := m.transformResponse("", inputResponse)
	var output bytes.Buffer
	if err := json.NewEncoder(&output).Encode(&outputResponse); err != nil {
		return "", err
	}
	return output.String(), nil
}

func (m *Mapper) transformResponse(previousInputPath string, inputResponse interface{}) (string, interface{}) {
	if outputKey, ok := m.mapping[previousInputPath]; ok {
		return outputKey, inputResponse
	}

	switch inputResponse := inputResponse.(type) {
	case map[string]interface{}:
		outputResponse := make(map[string]interface{})
		currentOutputPath := ""

		for inputKey, input := range inputResponse {
			outputPath, output := m.transformResponse(determinePath(previousInputPath, inputKey), input)
			if outputPath != "" {
				previousPath, outputKey := splitPath(outputPath)
				assertSingleOutputPath(currentOutputPath, previousPath)
				currentOutputPath = previousPath
				outputResponse[outputKey] = output
			}
		}
		return currentOutputPath, outputResponse
	case []interface{}:
		var extractedValues []interface{}
		var outputKey string
		for _, v := range inputResponse {
			var value interface{}
			outputKey, value = m.transformResponse(previousInputPath, v)
			extractedValues = append(extractedValues, value)
		}
		return outputKey, extractedValues

	default:
		return "", nil
	}
}

func assertSingleOutputPath(currentOutputPath string, previousPath string) {
	if currentOutputPath != "" && currentOutputPath != previousPath {
		panic("Error due to different mapping: `" + currentOutputPath + "` --- `" + previousPath + "`")
	}
}

func splitPath(path string) (string, string) {
	if strings.Contains(path, ".") {
		keys := strings.Split(path, ".")
		lastKey := keys[len(keys)-1]
		prefix := strings.Join(keys[0:len(keys)-1], ".")
		return prefix, lastKey
	}
	return "", path
}

func determinePath(previousInputKey string, inputKey string) string {
	if previousInputKey == "" {
		return inputKey
	}
	return previousInputKey + "." + inputKey
}
