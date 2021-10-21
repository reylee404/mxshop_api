package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"

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

	host := global.ServerConfig.UserSrvConfig.Host
	port := global.ServerConfig.UserSrvConfig.Port
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("GetUserList", zap.String("dial", err.Error()))
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
		c.JSON(http.StatusOK, response.NewSuccessResponse(map[string]string{
			"token": "aaaa",
		}))
	} else {
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(400, "用户名或密码错误"))
	}

}

func GetUserList(c *gin.Context) {
	page := utils.DefaultAtoi(c.Query("page"), 1)
	pageSize := utils.DefaultAtoi(c.Query("page_size"), 10)

	host := global.ServerConfig.UserSrvConfig.Host
	port := global.ServerConfig.UserSrvConfig.Port
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("GetUserList", zap.String("dial", err.Error()))
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, err.Error()))
		return
	}
	client := proto.NewUserClient(conn)
	userList, err := client.GetUserList(context.Background(), &proto.PageInfo{
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
	data := response.UserListResponse{
		Total:    userList.Total,
		UserList: list,
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(data))
}
