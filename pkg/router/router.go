package router

import (
	"encoding/json"
	"fmt"
	"opendatahubchallenge/pkg/mobility"
	"opendatahubchallenge/pkg/service"
	"opendatahubchallenge/pkg/tourism"
	"os"
	"path"
	"strings"
)

type Route struct {
	Keyword string
	Service string
}

type Config struct {
	Routes []Route
}

type Router struct {
	config *Config
	// "tourism"
}

func NewDefaultRouter() Router {
	return NewRouter("config.json")
}

func NewRouter(fileName string) Router {
	var config, err = readConfigFromFile(fileName)
	if err != nil {
		panic(err)
	}
	var router = Router{
		config: config,
	}
	return router
}

func (r Router) EntryPoint(path string, method string) service.Response {
	var response service.Response
	configurations := r.config
	var s service.Service

	// TODO: would a map[string] service.Service help us here?
	for _, route := range configurations.Routes {
		if strings.Contains(path, route.Keyword) {
			switch route.Service {
			case "tourism":
				// TODO: inject the services rather than "creating them" here
				s = tourism.TourismService{}
			case "mobility":
				s = mobility.MobilityService{}
			default:
				continue
			}

			var err error
			response, err = s.ExecuteRequest(method, path, nil)
			if err != nil {
				panic(err)
			}
			break
		}
	}

	if response.StatusCode != 200 {
		response = AttemptRequest(response, method, path)
	}

	return response
}

func AttemptRequest(response service.Response, method string, path string) service.Response {
	var serviceTypes = []service.Service{tourism.TourismService{}, mobility.MobilityService{}}
	for _, serviceType := range serviceTypes {
		var err error
		response, err = serviceType.ExecuteRequest(method, path, nil)
		if response.StatusCode == 200 {
			break
		}

		if err != nil {
			panic(err)
		}
	}
	return response
}

func readConfigFromFile(fileName string) (*Config, error) {
	var configData Config
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	extension := strings.ToLower(path.Ext(fileName))
	if extension != ".json" {
		return nil, fmt.Errorf("unsupported configuration file extension (`%s`): %s", extension, fileName)
	}

	err = json.Unmarshal([]byte(data), &configData)
	if err != nil {
		return nil, err
	}
	return &configData, nil
}
