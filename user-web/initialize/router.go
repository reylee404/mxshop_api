package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user-web/router"
)

func InitGinAndRouters() *gin.Engine {
	engine := gin.Default()

	userV1Group := engine.Group("/u/v1")
	{
		router.RegisterUserRouter(userV1Group)
		router.RegisterBaseRouter(userV1Group)
	}
	return engine
}
