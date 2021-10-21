package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user-web/api"
	"mxshop_api/user-web/middlewares"
)

func RegisterUserRouter(routerGroup *gin.RouterGroup) {
	group := routerGroup.Group("/user")
	{
		group.POST("/pwd_login", api.PasswordLogin)

		group.Use(middlewares.JWTAuth())
		group.GET("/list", middlewares.IsAdminAuth(), api.GetUserList)
	}
}
