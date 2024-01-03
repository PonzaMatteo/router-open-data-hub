package service

type Response struct {
	Body       string
	StatusCode int
}

type Service interface {
	// TODO: we should change the response type to (Response, error)
	ExecuteRequest(method string, path string, body []byte) Response
}
