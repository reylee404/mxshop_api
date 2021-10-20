package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"mxshop_api/user-web/global/response"
	"mxshop_api/user-web/proto"
	"net/http"
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
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		zap.L().Error("GetUserList", zap.String("dial", err.Error()))
		c.JSON(http.StatusOK, response.NewFailedUserListResponse(500, err.Error()))
		return
	}
	client := proto.NewUserClient(conn)
	userList, err := client.GetUserList(context.Background(), &proto.PageInfo{
		PIndex: 1,
		PSize:  2,
	})
	if err != nil {
		zap.L().Error("GetUserList", zap.String("SrvClientErr", err.Error()))
		c.JSON(http.StatusOK, response.NewFailedUserListResponse(500, GrpcErrToHttpMessage(err)))
		return
	}

	data := make([]response.UserResponse, 0, 10)
	for _, u := range userList.Data {
		userResponse := response.UserResponse{
			Id:       u.Id,
			Mobile:   u.Mobile,
			NickName: u.NickName,
			Birthday: u.Birthday,
		}
		data = append(data, userResponse)
	}
	c.JSON(http.StatusOK, response.NewSuccessUserListResponse(data))
}
