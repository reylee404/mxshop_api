package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"mxshop_api/user-web/global/response"
	"mxshop_api/user-web/proto"
	"net/http"
)

func GrpcErrToHttpMessage(err error) string {
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return s.Message()
		}
		return "Unknown Error"

	}
	return "OK"

}

func GetUserList(c *gin.Context) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusOK, response.NewUserListResponse(500, err.Error()))
		return
	}
	client := proto.NewUserClient(conn)
	list, err := client.GetUserList(context.Background(), &proto.PageInfo{
		PIndex: 1,
		PSize:  2,
	})
	if err != nil {
		c.JSON(http.StatusOK, response.NewUserListResponse(500, GrpcErrToHttpMessage(err)))
		return
	}

	result := response.NewSuccessUserListResponse()
	for _, u := range list.Data {
		userResponse := response.UserResponse{
			Id:       u.Id,
			Mobile:   u.Mobile,
			NickName: u.NickName,
			Birthday: u.Birthday,
		}
		result.Data = append(result.Data, userResponse)
	}
	c.JSON(http.StatusOK, result)
}
