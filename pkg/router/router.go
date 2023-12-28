package router

import (
	"encoding/json"
	"opendatahubchallenge/pkg/mobility"
	"opendatahubchallenge/pkg/service"
	"opendatahubchallenge/pkg/tourism"
	"os"
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
	for _, route := range configurations.Routes {
		if strings.Contains(path, route.Keyword) {
			switch route.Service {
			case "tourism":
				s = tourism.TourismService{}
				break
			case "mobility":
				s = mobility.MobilityService{}
				break
			default:
				continue
			}
		}
	}
	response = s.ExecuteRequest(method, path, nil)

	return response
}

func readConfigFromFile(fileName string) (*Config, error) {
	var configData Config

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &configData)
	if err != nil {
		return nil, err
	}

	return &configData, nil
}
