package middlewares

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user-web/global/response"
	"mxshop_api/user-web/models"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, exists := context.Get("claims")
		if !exists {
			context.JSON(http.StatusOK, response.NewFailedBaseResponse(400, "无权限"))
			context.Abort()
			return
		}
		customClaims := claims.(*models.CustomClaims)
		if customClaims.AuthorityId != 2 {
			context.JSON(http.StatusOK, response.NewFailedBaseResponse(400, "无权限"))
			context.Abort()
			return
		}
		context.Next()

	}

}
