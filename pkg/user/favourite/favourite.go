package favourite

import (
	"context"
	user "go-boilerplate-api/apis/grpc/generated/user"
	"go-boilerplate-api/config"
	log "go-boilerplate-api/pkg/utils/logger"
)

// FavouriteInterface is implemented by any value that contains the required methods
// makes Favourite mockable
type FavouriteInterface interface {
	Get(GetRequest) (*GetResponse, error)
}

// NewFavourite create an instance of favourite
func NewFavourite(conf config.IConfig, grpcConn grpcConnectioner) FavouriteInterface {
	return &Favourite{config: conf, grpcCon: grpcConn, client: new(gclient)}
}

// Favourite contains methods to peform operations on users favourites
type Favourite struct {
	config  config.IConfig
	grpcCon grpcConnectioner
	client  grpcClient
}

// Get Make the request to favourites
func (f *Favourite) Get(req GetRequest) (*GetResponse, error) {
	resp, err := f.client.GetFav(f.grpcCon)
	if err != nil {
		return nil, err
	}

	log.Info("Response : ", resp)
	res := GetResponse{}
	res.Beers = []string{"Moon Shine", "Bira", "Simba"}
	return &res, nil
}

// The reason for this is to avoid calling the actual grpc functions during testing , need to find a better way around
func (c *gclient) GetFav(grpcCon grpcConnectioner) (*user.Users, error) {
	favGrpcCon := grpcCon.GetFavourite()
	cli := user.NewUserServiceClient(favGrpcCon)
	return cli.GetAll(context.Background(), &user.UserGetRequest{})
}
