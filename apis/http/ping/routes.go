package ping

import "github.com/gin-gonic/gin"

// NewPingRoute Creates and initializes ping route
func NewPingRoute(router *gin.Engine) {
	bindRoutes(router)
}

func bindRoutes(router *gin.Engine) {
	service := NewPingService()
	routerAPI := router.Group("/ping")
	{
		routerAPI.GET("/", service.Get)

	}
}
