package main

import (
	"io"
	"net/http"
)

func main() {
	getAccommodation()
}

type Response struct {
	body       string
	statusCode int
}

func getAccommodation() (Response, error) {

	accommodationUrl := "https://tourism.opendatahub.com/v1/Accommodation"

	request, err := http.NewRequest("GET", accommodationUrl, nil)

	if err != nil {
		return Response{}, err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
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
