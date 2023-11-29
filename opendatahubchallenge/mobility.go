package main

import (
	"fmt"
	"io"
	"net/http"
)

type TourismService struct {
}

func (TourismService) ExecuteRequest(path string) Response {
	tourismPath := "https://tourism.opendatahub.com" + path

	response, err := getRequest(tourismPath)
	fmt.Println("Error", err)
	return response
}

func getRequest(tourismPath string) (Response, error) {

	request, err := http.NewRequest("GET", tourismPath, nil)

	if err != nil {
		return Response{}, err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

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
