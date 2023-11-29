package main

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestTourismService(t *testing.T) {
	t.Run("Tourism service should connect to an existing API", func(t *testing.T) {
		var service = TourismService{}

		var response Response = service.ExecuteRequest("/v1/Accommodation")

		expected := 200
		actual := response.statusCode
		assert.Equal(t, expected, actual, "Wrong Status Code")

	})

	t.Run("Tourism service should fail with wrong path", func(t *testing.T) {
		var service = TourismService{}

		var response Response = service.ExecuteRequest("/v1/path-not-exist")

		expected := 404
		actual := response.statusCode
		assert.Equal(t, expected, actual, "Wrong Status Code")

	})
}

func TestGetFixedValue(t *testing.T) {
	defer gock.Off()

	gock.New("https://tourism.opendatahub.com").
		Get("/v1/Accommodation").
		Reply(200).
		JSON(map[string]string{"value": "fixed"})

	var service = TourismService{}
	var response Response = service.ExecuteRequest("/v1/Accommodation")

	expected := 200
	actual := response.statusCode
	assert.Equal(t, expected, actual, "Wrong Status Code")

	assert.Equal(t, "{\"value\":\"fixed\"}\n", response.body, "Wrong Status Code")
}
