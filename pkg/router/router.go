package router

import (
	"log"
	"opendatahubchallenge/pkg/config"
	"opendatahubchallenge/pkg/mapper"
	"opendatahubchallenge/pkg/mobility"
	"opendatahubchallenge/pkg/service"
	"opendatahubchallenge/pkg/tourism"
	"strings"
)

type Router struct {
	config      *config.Config
	serviceByID map[string]service.Service
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
		config:      config,
		serviceByID: make(map[string]service.Service),
	}
}

func (r *Router) AddService(serviceID string, serviceType service.Service) {
	r.serviceByID[serviceID] = serviceType
}

func (r *Router) RouteRequest(path string, method string) (*service.Response, error) {
	route, ok := r.identifyRoute(path)
	if !ok {
		return r.attemptRequestUntilFound(method, path)
	}

	s := r.serviceByID[route.Service]
	response, err := s.ExecuteRequest(method, path, nil)
	if err != nil {
		return nil, err
	}

	m := createMapper(route)
	newBody, err := m.MapJSON(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body = newBody

	return &response, nil
}

func (r *Router) identifyRoute(path string) (config.Route, bool) {
	for _, route := range r.config.Routes {
		if strings.Contains(strings.ToLower(path), strings.ToLower(route.Keyword)) {
			log.Println("[router]: identified service ", route.Service, " for serving the request", path)
			return route, true
		}
	}
	return config.Route{}, false
}

func (r *Router) attemptRequestUntilFound(method string, path string) (*service.Response, error) {
	for id, serviceType := range r.serviceByID {
		log.Println("[router]: attempting to contact service ", id, "for request", method, path)

		var err error
		response, err := serviceType.ExecuteRequest(method, path, nil)
		if err != nil {
			log.Println("[router]: service", id, "responded with error, skip and trying with the next one")
			continue
		}

		if response.StatusCode == 200 {
			return &response, nil
		}
	}

	return &service.Response{
		StatusCode: 404,
		Body:       "",
	}, nil
}

func createMapper(route config.Route) mapper.Mapper {
	var m mapper.Mapper
	if route.Mapping == nil {
		m = mapper.EmptyMapper()
	} else {
		m = mapper.NewMapper(*route.Mapping)
	}
	return m
}
