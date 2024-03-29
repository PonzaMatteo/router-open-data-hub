package tourism

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestTourismService(t *testing.T) {
	t.Run("Tourism service should connect to an existing API", func(t *testing.T) {
		var service = TourismService{}

		response, err := service.ExecuteRequest(http.MethodGet, "/v1/Accommodation", nil)
		assert.NoError(t, err)

		expected := 200
		actual := response.StatusCode
		assert.Equal(t, expected, actual, "Wrong Status Code")

	})

	t.Run("Tourism service should fail with wrong path", func(t *testing.T) {
		var service = TourismService{}

		response, err := service.ExecuteRequest(http.MethodGet, "/v1/path-not-exist", nil)
		assert.NoError(t, err)

		expected := 404
		actual := response.StatusCode
		assert.Equal(t, expected, actual, "Wrong Status Code")

	})

	t.Run("Tourism POST service to an existing API without authorization", func(t *testing.T) {
		var service = TourismService{}

		body, err := json.Marshal("string")
		assert.NoError(t, err)
		response, err := service.ExecuteRequest(http.MethodPost, "/v1/AccommodationAvailable", body)
		assert.NoError(t, err)

		expected := 401
		actual := response.StatusCode
		assert.Equal(t, expected, actual, "Wrong Status Code")

	})
}

func TestTourismServiceUsingMock(t *testing.T) {
	t.Run("Tourism GET service should connect to API using mock", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://tourism.opendatahub.com").
			Get("/v1/Accommodation").
			Reply(200).
			JSON(map[string]string{"value": "fixed"})

		var service = TourismService{}
		response, err := service.ExecuteRequest(http.MethodGet, "/v1/Accommodation", nil)
		assert.NoError(t, err)

		expected := 200
		actual := response.StatusCode
		assert.Equal(t, expected, actual, "Wrong Status Code")

		assert.Equal(t, "{\"value\":\"fixed\"}\n", response.Body, "Wrong Status Code")
	})

	t.Run("Tourism POST service to API without authorization using mock", func(t *testing.T) {
		defer gock.Off()
		gock.Observe(gock.DumpRequest)

		gock.New("https://tourism.opendatahub.com").
			Post("/v1/AccommodationAvailable").
			MatchType("json").
			JSON(map[string]string{"message": "hello"}).
			Reply(401).
			JSON(map[string]string{"value": "fixed"})

		body, err := json.Marshal(map[string]string{"message": "hello"})
		assert.NoError(t, err)
		var service = TourismService{}
		response, err := service.ExecuteRequest(http.MethodPost, "/v1/AccommodationAvailable", body)
		assert.NoError(t, err)

		expected := 401
		actual := response.StatusCode
		assert.Equal(t, expected, actual, "Wrong Status Code")

		assert.Equal(t, "{\"value\":\"fixed\"}\n", response.Body, "Wrong Status Code")
	})
}
