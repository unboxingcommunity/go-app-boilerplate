package user

import (
	"context"

	pb "go-boilerplate-api/apis/grpc/generated/user"
	"go-boilerplate-api/apis/grpc/utils"
	"go-boilerplate-api/apm"
	"go-boilerplate-api/config"
	"go-boilerplate-api/pkg/clients/db"
	grpcPkg "go-boilerplate-api/pkg/clients/grpc"
	httpPkg "go-boilerplate-api/pkg/clients/http"
	"go-boilerplate-api/pkg/user"
	"go-boilerplate-api/pkg/user/favourite"
	"go-boilerplate-api/pkg/user/rating"
	userRepo "go-boilerplate-api/pkg/user/repo"
	pkgUtils "go-boilerplate-api/pkg/utils"
)

// Service contains the methods required to perfom operation's on users (proto definition)
type Service struct {
	user user.UsersInterface
}

// NewUserService Create a new instance of a Service with the given dependencies.
func NewUserService(conf config.IConfig, db *db.Instances, apm apm.HandlerInterface, httpReq httpPkg.IRequest, grpcConn grpcPkg.IGrpcConnections) *Service {
	userRating := rating.NewRating(conf, httpReq)
	userFavourites := favourite.NewFavourite(conf, grpcConn)
	userRepo := userRepo.NewUserRepo(conf, db)
	userService := user.NewUser(conf, userRepo, apm, userRating, userFavourites)

	return &Service{user: userService}
}

// GetAll gets all users
func (service *Service) GetAll(ctx context.Context, req *pb.UserGetRequest) (res *pb.Users, err error) {
	defer utils.HandleError(&err)

	users, err := service.user.GetAll()
	if err != nil {
		return nil, err
	}

	res = &pb.Users{}
	err = pkgUtils.Bind(users, &res.Users)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetOne gets one users
func (service *Service) GetOne(ctx context.Context, req *pb.UserGetRequest) (res *pb.User, err error) {
	defer utils.HandleError(&err)

	userReq := user.User{}
	// Need to decode to user.User since User is an embedded struct
	err = pkgUtils.Bind(req, &userReq)
	if err != nil {
		return nil, err
	}

	users, err := service.user.GetOne(userReq.ID)
	if err != nil {
		return nil, err
	}

	res = &pb.User{}
	err = pkgUtils.Bind(users, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Insert stores a user in the datastore
func (service *Service) Insert(ctx context.Context, req *pb.User) (res *pb.User, err error) {
	defer utils.HandleError(&err)

	userReq := user.User{}
	// Need to decode to user.User since User is an embedded struct
	err = pkgUtils.Bind(req, &userReq)
	if err != nil {
		return nil, err
	}

	err = service.user.Insert(userReq)
	if err != nil {
		return nil, err
	}
	res = nil

	return res, nil
}

// GetWithInfo gets a user from the database along with rating and favourites
func (service *Service) GetWithInfo(ctx context.Context, req *pb.UserGetRequest) (res *pb.User, err error) {
	defer utils.HandleError(&err)

	userReq := user.User{}
	// Need to decode to user.User since User is an embedded struct
	err = pkgUtils.Bind(req, &userReq)
	if err != nil {
		return nil, err
	}

	users, err := service.user.GetWithInfo(userReq.ID)
	if err != nil {
		return nil, err
	}

	res = &pb.User{}
	err = pkgUtils.Bind(users, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
