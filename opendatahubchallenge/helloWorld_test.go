package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestResponse(t *testing.T) {

	// t.Run("successful response", func(t *testing.T) {
	// 	response, _ := getAccommodation()
	// 	got := response.statusCode
	// 	want := 200
	// 	if got != want {
	// 		t.Errorf("got %q want %q", got, want)
	// 	}
	// })

	// t.Run("compare response body", func(t *testing.T) {
	// 	currentResponse, _ := getAccommodation()
	// 	got := extractFloatFromJson(t, currentResponse.body, "TotalResults")

	// 	sampleResponse := readFile(t, "../response-samples/tourism-accommodations.json")
	// 	want := extractFloatFromJson(t, string(sampleResponse), "TotalResults")

	// 	if got != want {
	// 		t.Errorf("got %v want %v", got, want)
	// 	}
	// })
}

func extractFloatFromJson(t *testing.T, currentResponse string, keyField string) float64 {
	var accomResponse map[string]interface{}
	err := json.Unmarshal([]byte(currentResponse), &accomResponse)
	if err != nil {
		t.Errorf("Error parsing Json String in got %v", err)
	}
	got, ok := accomResponse[keyField].(float64)
	if !ok {
		t.Errorf("key field %q is not an integer", keyField)
	}
	return got
}

func readFile(t *testing.T, fileName string) []byte {
	data, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("File reading error %v", err)
		return nil
	}
	return data
}
