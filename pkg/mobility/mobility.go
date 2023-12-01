package mobility

import (
	"opendatahubchallenge/pkg/service"
)

type MobilityService struct{}
type Message struct {
	Body string
}

func (MobilityService) ExecuteRequest(method string, path string, body []byte) service.Response {
	// mobilityPath := "https://mobility.opendatahub.com" + path

	return service.Response{}
}
