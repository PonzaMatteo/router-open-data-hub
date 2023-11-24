package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	fmt.Println("Hello, world.")
	getAccommodation()
}

type Response struct {
	body   string
	status int
}

func getAccommodation() int {

	fmt.Println("Getting Accommodation...")

	accommodationUrl := "https://tourism.opendatahub.com/v1/Accommodation"

	request, error := http.NewRequest("GET", accommodationUrl, nil)

	if error != nil {
		fmt.Println(error)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
		fmt.Println(error)
	}

	responseBody, error := io.ReadAll(response.Body)

	if error != nil {
		fmt.Println(error)
	}

	formattedData := string(responseBody)

	fmt.Println("Status: ", response.Status)
	fmt.Println("Response body: ", formattedData)

	// clean up memory after execution
	defer response.Body.Close()
	return response.StatusCode
}
