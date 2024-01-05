package router

import (
	"opendatahubchallenge/pkg/config"
	"opendatahubchallenge/pkg/mapper"
	"opendatahubchallenge/pkg/mobility"
	"opendatahubchallenge/pkg/service"
	"opendatahubchallenge/pkg/tourism"
	"strings"
)

type Router struct {
	config       *config.Config
	serviceTypes map[string]service.Service
}

func NewDefaultRouter() (Router, error) {
	defaultConfig := config.GetDefault()
	router := NewRouter(defaultConfig)
	router.AddService("tourism", tourism.TourismService{})
	router.AddService("mobility", mobility.MobilityService{})
	return router, nil
}

func NewRouter(config *config.Config) Router {
	return Router{
		config:       config,
		serviceTypes: make(map[string]service.Service),
	}
}

func (r *Router) AddService(serviceID string, serviceType service.Service) {
	r.serviceTypes[serviceID] = serviceType
}

func (r *Router) EntryPoint(path string, method string) (*service.Response, error) {
	for _, route := range r.config.Routes {
		if strings.Contains(strings.ToLower(path), strings.ToLower(route.Keyword)) {
			s := r.serviceTypes[route.Service]
			var err error // declare error first to avoid shadowing
			response, err := s.ExecuteRequest(method, path, nil)
			if err != nil {
				return nil, err
			}

			// TODO: review here, important part!
			if route.Mapping != nil {
				var mapper = mapper.NewMapperWithMapping(*route.Mapping)
				var newBody, err = mapper.Transform(response.Body)
				if err != nil {
					panic(err)
				}
				response.Body = newBody
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
