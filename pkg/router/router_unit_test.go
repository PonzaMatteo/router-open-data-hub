package router

import (
	"opendatahubchallenge/pkg/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterWithMock(t *testing.T) {
	// we create a router that have the following config
	// routes:
	//   - keyword: "test_keyword"
	//     service: "test_service
	// test_service is a mock
	t.Run("Router should connect to mock service", func(t *testing.T) {

		ts := testService{
			executeRequest: func(method, path string, body []byte) (service.Response, error) {
				return service.Response{Body: "response from test service"}, nil
			},
		}

		path := "/test_keyword"
		method := "GET"

		var router = NewRouter("mockfile.json")
		router.AddService("test_service", ts)

		_ = router.EntryPoint(path, method)

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

		var router = NewRouter("mockfile.json")
		router.AddService("test_service_1", ts1)
		router.AddService("test_service", ts)

		response := router.EntryPoint(path, method)
		assert.Equal(t, "response from test service", response.Body)

	})
}

type testService struct {
	executeRequest func(method string, path string, body []byte) (service.Response, error)
}

func (t testService) ExecuteRequest(method string, path string, body []byte) (service.Response, error) {
	return t.executeRequest(method, path, body)
}
