package service

type Response struct {
	Body       string
	StatusCode int
}

type Service interface {
	ExecuteRequest(method string, path string, body []byte) Response
}
