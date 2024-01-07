package mobility

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/PonzaMatteo/router-open-data-hub/pkg/service"
)

type MobilityService struct{}
type Message struct {
	Body string
}

func (MobilityService) ExecuteRequest(method string, path string, body []byte) (service.Response, error) {
	mobilityPath := "https://mobility.api.opendatahub.com" + path
	response, err := request(mobilityPath, method, body)
	if err != nil {
		return service.Response{}, fmt.Errorf("failed to execute request to mobility service: %w", err)
	}
	return response, nil
}

func request(mobilityPath string, method string, body []byte) (service.Response, error) {
	request, err := http.NewRequest(method, mobilityPath, bytes.NewBuffer(body))

	if err != nil {
		return service.Response{}, err
	}

	request.Header.Set("Content-Type", "application/json")

	var client *http.Client = http.DefaultClient
	response, err := client.Do(request)

	if err != nil {
		return service.Response{}, err
	}

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		return service.Response{}, err
	}

	var result = service.Response{Body: string(responseBody), StatusCode: response.StatusCode}

	// clean up memory after execution
	defer response.Body.Close()
	return result, nil
}
