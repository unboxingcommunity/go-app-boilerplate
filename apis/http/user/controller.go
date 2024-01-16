package user

import (
	"net/http"

	"go-boilerplate-api/apis/http/utils"
	"go-boilerplate-api/apm"
	"go-boilerplate-api/config"
	"go-boilerplate-api/pkg/clients/db"
	grpcPkg "go-boilerplate-api/pkg/clients/grpc"
	httpPkg "go-boilerplate-api/pkg/clients/http"
	user "go-boilerplate-api/pkg/user"
	"go-boilerplate-api/pkg/user/favourite"
	"go-boilerplate-api/pkg/user/rating"
	userRepo "go-boilerplate-api/pkg/user/repo"

	"github.com/gin-gonic/gin"
	"github.com/ralstan-vaz/go-errors"
)

// Service contains the methods required to perfom operation's on users
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

func (service *Service) getAll(ctx *gin.Context) {
	var err error
	defer utils.HandleError(ctx, &err)

	users, err := service.user.GetAll()
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (service *Service) getOne(ctx *gin.Context) {
	var err error
	defer utils.HandleError(ctx, &err)

	userID := ctx.Param("userID")
	users, err := service.user.GetOne(userID)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (service *Service) getWithInfo(ctx *gin.Context) {
	var err error
	defer utils.HandleError(ctx, &err)

	userID := ctx.Param("userID")
	users, err := service.user.GetWithInfo(userID)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (service *Service) insert(ctx *gin.Context) {
	var err error
	defer utils.HandleError(ctx, &err)

	var user user.User
	if err = ctx.ShouldBindJSON(&user); err != nil {
		err = errors.NewBadRequest("Could not bind request to model").SetCode("APIS.HTTP.USER.REQUEST_BIND_FAILD")
		return
	}

	err = service.user.Insert(user)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
