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
	config       *Config
	serviceTypes map[string]service.Service
}

func NewDefaultRouter() Router {
	router := NewRouter("config.json")
	router.AddService("tourism", tourism.TourismService{})
	router.AddService("mobility", mobility.MobilityService{})
	return router
}

func NewRouter(fileName string) Router {
	var config, err = readConfigFromFile(fileName)
	if err != nil {
		panic(err)
	}
	var router = Router{
		config:       config,
		serviceTypes: make(map[string]service.Service),
	}
	return router
}

func (r *Router) AddService(serviceID string, serviceType service.Service) {
	r.serviceTypes[serviceID] = serviceType
}

func (r *Router) EntryPoint(path string, method string) (*service.Response, error) {
	for _, route := range r.config.Routes {
		if strings.Contains(path, route.Keyword) {
			s := r.serviceTypes[route.Service]
			var err error // declare error first to avoid shadowing
			response, err := s.ExecuteRequest(method, path, nil)
			if err != nil {
				return nil, err
			}
			return &response, nil
		}
	}

	response, err := r.AttemptRequest(method, path)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *Router) AttemptRequest(method string, path string) (*service.Response, error) {
	var response service.Response
	for _, serviceType := range r.serviceTypes {
		var err error
		response, err = serviceType.ExecuteRequest(method, path, nil)
		if err != nil {
			return nil, err
		}
		if response.StatusCode == 200 {
			break
		}
	}
	return &response, nil
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
