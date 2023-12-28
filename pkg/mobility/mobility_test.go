package mobility

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestMobilityService(t *testing.T) {

	// t.SkipNow()

	t.Run("Mobility service should connect to an existing API", func(t *testing.T) {
		var service = MobilityService{}

		response := service.ExecuteRequest(http.MethodGet, "/v2/tree,node", nil)

		expected := 200
		actual := response.StatusCode
		assert.Equal(t, expected, actual, "Wrong Status Code")

	})

	t.Run("Mobility service should fail with wrong path", func(t *testing.T) {
		var service = MobilityService{}

		response := service.ExecuteRequest(http.MethodGet, "/v/path-not-exist", nil)

		expected := 404
		actual := response.StatusCode
		assert.Equal(t, expected, actual, "Wrong Status Code")

	})
}

func TestMobilityServiceUsingMock(t *testing.T) {
	t.Run("Mobility GET service should connect to API using mock", func(t *testing.T) {

		defer gock.Off()

		gock.New("https://mobility.api.opendatahub.com").
			Get("/v2/tree,node").
			Reply(200).
			JSON(map[string]string{"value": "fixed"})

		var service = MobilityService{}
		response := service.ExecuteRequest(http.MethodGet, "/v2/tree,node", nil)

		assert.Equal(t, 200, response.StatusCode, "Wrong Status Code")

		assert.Equal(t, "{\"value\":\"fixed\"}\n", response.Body, "Wrong Status Code")
	})
}
