package favourite

import (
	user "go-boilerplate-api/apis/grpc/generated/user"

	"google.golang.org/grpc"
)

// grpcConnectioner contains methods to retrieve grpc connections
// this makes it possible to pass multiple grpc connectiions incase the package needs it
type grpcConnectioner interface {
	GetFavourite() *grpc.ClientConn
}

// Used to provide interface for calling proto client funcs
type grpcClient interface {
	GetFav(grpcCon grpcConnectioner) (*user.Users, error)
}

// dummy struct used for interfacing
type gclient struct {
}

// GetRequest request for getting a favourite
type GetRequest struct {
	ID string
}

// GetResponse response after getting a favourite
type GetResponse struct {
	ID    string   `json:"id,omitempty"`
	Beers []string `json:"beers,omitempty"`
}
