package user

import (
	"go-boilerplate-api/shared"

	"github.com/gin-gonic/gin"
)

// NewUserRoute Creates and initializes user routes
func NewUserRoute(router *gin.Engine, deps *shared.Deps) {
	bindRoutes(router, deps)
}

func bindRoutes(router *gin.Engine, deps *shared.Deps) {
	service := NewUserService(deps.Config, deps.Database, deps.Apm, deps.HTTPRequester, deps.GrpcConn)
	userAPI := router.Group("/users")
	{
		userAPI.GET("/", service.getAll)
		userAPI.GET("/:userId", service.getOne)
		userAPI.GET("/:userId/rating", service.getWithInfo)
		userAPI.POST("/", service.insert)
	}
}
