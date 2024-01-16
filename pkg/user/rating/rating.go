package rating

import (
	"encoding/json"
	"io/ioutil"

	"go-boilerplate-api/config"
	httpPkg "go-boilerplate-api/pkg/clients/http"
)

// Rater is implemented by any value that contains the required methods
type Rater interface {
	Get(GetRequest) (*GetResponse, error)
}

// NewRating creates a new instance of Rating
func NewRating(conf config.IConfig, httpRequester httpPkg.IRequest) Rater {
	return &Rating{config: conf, httpRequester: httpRequester}
}

// Rating contains methods to the perform operations on ratings
type Rating struct {
	config        config.IConfig
	httpRequester httpPkg.IRequest
}

// Get makes the request to get the ratings
func (r *Rating) Get(req GetRequest) (*GetResponse, error) {

	httpReq, err := r.httpRequester.New(r.config.Get().User.RatingsUrl)
	if err != nil {
		return nil, err
	}

	// here do whatever you want
	// httpReq.Req.Header("", "")
	// with adding features to the request object

	resp, err := r.httpRequester.Get(httpReq)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := GetResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
