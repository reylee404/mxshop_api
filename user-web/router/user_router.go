package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user-web/api"
)

func RegisterUserRouter(routerGroup *gin.RouterGroup) {
	group := routerGroup.Group("/user")
	{
		group.GET("/list", api.GetUserList)
		group.POST("/pwd_login", api.PasswordLogin )
	}
}
