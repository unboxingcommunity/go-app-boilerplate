package http

import (
	"net/http"
)

// IRequest ...
type IRequest interface {
	New(URL string) (*InnerRequest, error)
	Get(r *InnerRequest) (*http.Response, error)
}

// NewRequest Creates an instance if a request
func NewRequest() IRequest {
	return &Request{}
}

// Request , is an invoker struct for the interface
type Request struct {
}

// InnerRequest contains a method to perform an HTTP request
type InnerRequest struct {
	Req *http.Request
}

// New ...
func (r *Request) New(URL string) (*InnerRequest, error) {
	var err error
	newIReq := InnerRequest{Req: nil}
	newIReq.Req, err = http.NewRequest("", URL, nil)
	if err != nil {
		return nil, err
	}
	return &newIReq, nil
}

// Get makes an http get request
// Interceptors can be added here for tracing etc
func (r *Request) Get(req *InnerRequest) (*http.Response, error) {
	req.Req.Method = "GET"
	return do(req)
}

func do(r *InnerRequest) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(r.Req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
