package api

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"mxshop_api/user-web/middlewares"
	"mxshop_api/user-web/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop_api/user-web/forms"
	"mxshop_api/user-web/global"
	"mxshop_api/user-web/global/response"
	"mxshop_api/user-web/proto"
	"mxshop_api/user-web/utils"
)

func PasswordLogin(c *gin.Context) {
	var passwordLoginForm forms.PasswordLoginForm
	bind, ok := utils.RequestBind(c, &passwordLoginForm)
	if !ok {
		c.JSON(http.StatusOK, bind)
		return
	}

	verify := store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Answer, true)
	if !verify {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(400, "验证码错误"))
		return
	}

	conn, err := utils.GetServerConnFromConsul(global.ServerConfig.UserSrvConfig.Name)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, err.Error()))
		return
	}

	client := proto.NewUserClient(conn)

	user, err := client.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(400, utils.GrpcErrToHttpMessage(err)))
		return
	}

	checkPassword, err := client.CheckPassword(context.Background(), &proto.CheckPasswordInfo{
		Password:          passwordLoginForm.Password,
		EncryptedPassword: user.Password,
	})
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(400, utils.GrpcErrToHttpMessage(err)))
		return
	}

	if checkPassword.Success {
		j := middlewares.NewJWT()
		token, err := j.CreateToken(models.CustomClaims{
			Id:          uint(user.Id),
			NickName:    user.NickName,
			AuthorityId: uint(user.Role),
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
				Issuer:    "Lynn",
			},
		})

		if err != nil {
			c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, err.Error()))
			return
		}
		c.JSON(http.StatusOK, response.NewSuccessResponse(response.PasswordLoginResponse{
			Id:       user.Id,
			NickName: user.NickName,
			Token:    token,
		}))
	} else {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(400, "用户名或密码错误"))
	}

}

func GetUserList(c *gin.Context) {
	page := utils.DefaultAtoi(c.Query("page"), 1)
	pageSize := utils.DefaultAtoi(c.Query("page_size"), 10)

	conn, err := utils.GetServerConnFromConsul(global.ServerConfig.UserSrvConfig.Name)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, err.Error()))
		return
	}

	grpcClient := proto.NewUserClient(conn)
	userList, err := grpcClient.GetUserList(context.Background(), &proto.PageInfo{
		PIndex: uint32(page),
		PSize:  uint32(pageSize),
	})
	if err != nil {
		zap.L().Error("GetUserList", zap.String("SrvClientErr", err.Error()))
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, utils.GrpcErrToHttpMessage(err)))
		return
	}

	list := make([]response.UserResponse, 0, 10)
	for _, u := range userList.Data {
		userResponse := response.UserResponse{
			Id:       u.Id,
			Mobile:   u.Mobile,
			NickName: u.NickName,
			Birthday: u.Birthday,
		}
		list = append(list, userResponse)
	}
	r := response.UserListResponse{
		Total:    userList.Total,
		UserList: list,
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(r))
}

func Register(c *gin.Context) {
	form := forms.RegisterForm{}
	if bind, ok := utils.RequestBind(c, &form); !ok {
		c.JSON(http.StatusOK, bind)
		return
	}

	if result, err := global.RedisClient.Get(context.Background(), form.Mobile).Result(); err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusOK, response.NewFailedBaseResponse(400, "验证码不存在"))
		} else {
			c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, err.Error()))
		}
		return
	} else {
		if result != form.Code {
			c.JSON(http.StatusOK, response.NewFailedBaseResponse(400, "验证码错误"))
			return
		}
	}

	conn, err := utils.GetServerConnFromConsul(global.ServerConfig.UserSrvConfig.Name)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, err.Error()))
		return
	}
	client := proto.NewUserClient(conn)

	user, err := client.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   form.Mobile,
		Password: form.Password,
		NickName: form.Mobile,
	})
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(400, err.Error()))
		return
	}

	newJWT := middlewares.NewJWT()
	token, err := newJWT.CreateToken(models.CustomClaims{
		Id:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
			Issuer:    "Lynn",
		},
	})
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(response.PasswordLoginResponse{
		Id:       user.Id,
		NickName: user.NickName,
		Token:    token,
	}))

}
