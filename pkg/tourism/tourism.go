package tourism

import (
	"bytes"
	"io"
	"net/http"
	"opendatahubchallenge/pkg/service"
)

type TourismService struct{}
type Message struct {
	Body string
}

func (TourismService) ExecuteRequest(method string, path string, body []byte) service.Response {
	tourismPath := "https://tourism.opendatahub.com" + path

	response, _ := request(tourismPath, method, body)
	return response
}

func request(tourismPath string, method string, body []byte) (service.Response, error) {

	request, err := http.NewRequest(method, tourismPath, bytes.NewBuffer(body))

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
