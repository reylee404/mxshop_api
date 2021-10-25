package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"mxshop_api/user-web/global/response"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	digitDriver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(digitDriver, store)
	id, b64s, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, err.Error()))
		return
	}

	captchaResponse := response.CaptchaResponse{
		CaptchaId:   id,
		CaptchaPath: b64s,
	}
	c.JSON(http.StatusOK, response.NewSuccessResponse(captchaResponse))
}
