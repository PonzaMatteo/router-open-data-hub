package router

import (
	"opendatahubchallenge/pkg/service"
	"opendatahubchallenge/pkg/tourism"
)

func EntryPoint(keywords []string, value string) service.Response {
	var response service.Response
	for _, keyword := range keywords {
		if keyword == "Accommodation" {
			if value != "" {
				response = GetAccommodation(value)
			}
		}
	}
	return response

}

func GetAccommodation(value string) service.Response {
	var service tourism.TourismService
	response := service.ExecuteRequest("GET", "/v1/Accommodation/"+value, nil)
	return response
}
