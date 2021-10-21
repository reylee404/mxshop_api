package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/status"
	"mxshop_api/user-web/global"
	"mxshop_api/user-web/global/response"
	"strconv"
	"strings"
)

func RequestBind(c *gin.Context, wantBind interface{}) (errInfo interface{}, ok bool) {
	if err := c.ShouldBind(wantBind); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return response.NewBaseResponse(300, "参数错误", err.Error()), false
		}
		return response.NewBaseResponse(300, "参数错误",
			removeTopStruct(errs.Translate(global.Trans))), false
	}

	return nil, true
}

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func GrpcErrToHttpMessage(err error) string {
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return s.String()
		}
		return "Unknown Error"

	}
	return "OK"

}

func DefaultAtoi(value string, defaultValue int) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		result = defaultValue
	}
	return result
}

