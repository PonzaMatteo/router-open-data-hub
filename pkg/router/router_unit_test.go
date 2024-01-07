package router

import (
	"errors"
	"opendatahubchallenge/pkg/config"
	"opendatahubchallenge/pkg/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterWithInjectedConfiguration(t *testing.T) {
	testConfig := &config.Config{
		Routes: []config.Route{
			{
				Keyword: "test_keyword",
				Service: "test_service",
			},
		},
	}

	t.Run("Router should connect to mock service", func(t *testing.T) {
		ts := testService{
			executeRequest: func(method, path string, body []byte) (service.Response, error) {
				return service.Response{Body: "response from test service"}, nil
			},
		}

		router := NewRouter(testConfig)
		router.AddService("test_service", ts)

		actual, err := router.RouteRequest("/test_keyword", "GET")
		expected := &service.Response{Body: "response from test service"}
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Router should connect to correct service", func(t *testing.T) {

		ts := testService{
			executeRequest: func(method, path string, body []byte) (service.Response, error) {
				return service.Response{Body: "response from test service", StatusCode: 200}, nil
			},
		}
		ts1 := testService{
			executeRequest: func(method, path string, body []byte) (service.Response, error) {
				return service.Response{Body: "response from test service 1", StatusCode: 200}, nil
			},
		}

		path := "/test_keyword"
		method := "GET"

		router := NewRouter(testConfig)

		router.AddService("test_service_1", ts1)
		router.AddService("test_service", ts)

		response, err := router.RouteRequest(path, method)
		assert.NoError(t, err)
		assert.Equal(t, "response from test service", response.Body)
	})

	t.Run("Router should map request correctly even with different case in the keywords", func(t *testing.T) {

		var requestedPaths []string

		ts := testService{
			executeRequest: func(method, path string, body []byte) (service.Response, error) {
				requestedPaths = append(requestedPaths, path)
				return service.Response{Body: "response from test service", StatusCode: 200}, nil
			},
		}

		ts1 := testService{
			executeRequest: func(method, path string, body []byte) (service.Response, error) {
				requestedPaths = append(requestedPaths, path)
				return service.Response{Body: "response from test service 1", StatusCode: 200}, nil
			},
		}

		router := NewRouter(testConfig)

		router.AddService("test_service_1", ts1)
		router.AddService("test_service", ts)

		response, err := router.RouteRequest("/TEST_KEYWORD", "GET")
		assert.NoError(t, err)
		assert.Equal(t, "response from test service", response.Body)

		// should not change the case of the path in the actual request
		assert.Equal(t, []string{"/TEST_KEYWORD"}, requestedPaths)
	})

	testConfigWithMapping := &config.Config{
		Routes: []config.Route{
			{
				Keyword: "test_keyword",
				Service: "test_service",
				Mapping: &map[string]string{"evuuid": "id"},
			},
		},
	}

	t.Run("Router with mapping should connect to mock service", func(t *testing.T) {
		ts := testService{
			executeRequest: func(method, path string, body []byte) (service.Response, error) {
				return service.Response{Body: `{
					"evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
				}`}, nil
			},
		}

		router := NewRouter(testConfigWithMapping)
		router.AddService("test_service", ts)

		actual, err := router.RouteRequest("/test_keyword", "GET")
		expected := `{
			"id": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
		}`
		assert.NoError(t, err)
		assert.JSONEq(t, expected, actual.Body)
	})

	t.Run("Router should return error if response body for the mapper is not a valid json", func(t *testing.T) {
		ts := testService{
			executeRequest: func(method, path string, body []byte) (service.Response, error) {
				return service.Response{
					Body: `this is not a valid json`,
				}, nil
			},
		}

		router := NewRouter(testConfigWithMapping)
		router.AddService("test_service", ts)

		_, err := router.RouteRequest("/test_keyword", "GET")
		assert.Error(t, err)
	})

	t.Run("Router should return error if service in the mapping returns error", func(t *testing.T) {
		ts := testService{
			executeRequest: func(method, path string, body []byte) (service.Response, error) {
				return service.Response{}, errors.New("TESTING ERROR")
			},
		}

		router := NewRouter(testConfigWithMapping)
		router.AddService("test_service", ts)

		_, err := router.RouteRequest("/test_keyword", "GET")
		assert.Error(t, err)
	})
}

type testService struct {
	executeRequest func(method string, path string, body []byte) (service.Response, error)
}

func (t testService) ExecuteRequest(method string, path string, body []byte) (service.Response, error) {
	return t.executeRequest(method, path, body)
}
