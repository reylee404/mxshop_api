package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	OK    = 200
	Error = 300
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Response(c *gin.Context, res interface{}) {
	c.JSON(http.StatusOK, res)
}
