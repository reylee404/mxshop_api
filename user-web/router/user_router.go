package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user-web/api"
	"mxshop_api/user-web/middlewares"
)

func RegisterUserRouter(routerGroup *gin.RouterGroup) {
	group := routerGroup.Group("/user")
	group.Use(middlewares.Cors())
	{
		group.POST("/register", api.Register)
		group.POST("/pwd_login", api.PasswordLogin)
	}

	group.Use(middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		group.GET("/list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
	}
}
