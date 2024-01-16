package grpc

import (
	"go-boilerplate-api/config"

	"google.golang.org/grpc"
)

// IGrpcConnections ...
type IGrpcConnections interface {
	initialize() error
	favouriteInit() error
	GetFavourite() *grpc.ClientConn
}

// GrpcConnections contains all the GRPC connections this app uses
type GrpcConnections struct {
	favouriteConnection *grpc.ClientConn
	conf                config.IConfig
}

// NewConnections creates an instance of initialized GrpcConnections
func NewConnections(conf config.IConfig) (IGrpcConnections, error) {
	grpcCons := newGrpcConnections(conf)
	err := grpcCons.initialize()
	if err != nil {
		return nil, err
	}
	return grpcCons, nil
}

// newGrpcConnections creates an instance of GrpcConnections
// It does not initialize the connections.
func newGrpcConnections(conf config.IConfig) *GrpcConnections {
	return &GrpcConnections{conf: conf}
}

// Initialize ..
func (g *GrpcConnections) initialize() error {
	err := g.favouriteInit()
	if err != nil {
		return err
	}
	return nil
}

// FavouriteInit ..
func (g *GrpcConnections) favouriteInit() error {
	favouriteCon, err := grpc.Dial(g.conf.Get().User.FavouritesUrl, grpc.WithInsecure())
	if err != nil {
		return err
	}
	g.favouriteConnection = favouriteCon
	return nil
}

// GetFavourite return the grpc connection for the favourite service
func (g *GrpcConnections) GetFavourite() *grpc.ClientConn {
	return g.favouriteConnection
}
