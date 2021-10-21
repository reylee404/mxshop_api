package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"mxshop_api/user-web/global"
	"mxshop_api/user-web/global/response"
	"mxshop_api/user-web/proto"
	"net/http"
	"strconv"
)

func GrpcErrToHttpMessage(err error) string {
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return s.String()
		}
		return "Unknown Error"

	}
	return "OK"

}

func GetUserList(c *gin.Context) {
	page := DefaultAtoi(c.Query("page"), 1)
	pageSize := DefaultAtoi(c.Query("page_size"), 10)

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
		c.JSON(http.StatusOK, response.NewFailedBaseResponse(500, GrpcErrToHttpMessage(err)))
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

func DefaultAtoi(value string, defaultValue int) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		result = defaultValue
	}
	return result
}
