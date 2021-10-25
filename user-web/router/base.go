package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user-web/api"
)

func RegisterBaseRouter(routerGroup *gin.RouterGroup) {
	group := routerGroup.Group("base")
	group.GET("captcha", api.GetCaptcha)
}
