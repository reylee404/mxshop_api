package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"mxshop_api/user-web/forms"
	"mxshop_api/user-web/global"
	"mxshop_api/user-web/global/response"
	"mxshop_api/user-web/utils"
	"net/http"
	"strings"
	"time"
)

func generateSMSCode(length int) string {
	numeric := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	sb := strings.Builder{}
	for i := 0; i < length; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSMS(c *gin.Context) {
	var sendSMSForm forms.SendSMSForm
	if bind, ok := utils.RequestBind(c, &sendSMSForm); !ok {
		c.JSON(http.StatusOK, bind)
		return
	}

	result, err := global.RedisClient.Exists(context.Background(), sendSMSForm.Mobile).Result()
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, err.Error()))
		return
	}
	if result == 1 {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(300, "发送太过频繁"))
		return
	}

	code := generateSMSCode(6)
	zap.L().Info("验证码", zap.String(sendSMSForm.Mobile, code))
	set := global.RedisClient.Set(context.Background(), sendSMSForm.Mobile, code, 300*time.Second)
	if set.Err() != nil {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, set.Err().Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewBaseResponse(200, "发送成功", nil))
}
