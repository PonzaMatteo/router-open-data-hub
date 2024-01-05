package router

import (
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

		_, err := router.EntryPoint("/test_keyword", "GET")
		assert.NoError(t, err)
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

		response, err := router.EntryPoint(path, method)
		assert.NoError(t, err)
		assert.Equal(t, "response from test service", response.Body)
	})
}

type testService struct {
	executeRequest func(method string, path string, body []byte) (service.Response, error)
}

func (t testService) ExecuteRequest(method string, path string, body []byte) (service.Response, error) {
	return t.executeRequest(method, path, body)
}
