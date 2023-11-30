package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type TourismService struct{}
type Message struct {
	Body string
}

func (TourismService) ExecuteRequest(method string, path string, body []byte) Response {
	tourismPath := "https://tourism.opendatahub.com" + path

	response, err := request(tourismPath, method, body)
	fmt.Println("Error", err)
	return response
}

func request(tourismPath string, method string, body []byte) (Response, error) {

	request, err := http.NewRequest(method, tourismPath, bytes.NewBuffer(body))

	if err != nil {
		return Response{}, err
	}

	request.Header.Set("Content-Type", "application/json")

	var client *http.Client = http.DefaultClient
	response, err := client.Do(request)

	if err != nil {
		return Response{}, err
	}

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		return Response{}, err
	}

	var result = Response{body: string(responseBody), statusCode: response.StatusCode}

	// clean up memory after execution
	defer response.Body.Close()
	return result, nil
}
