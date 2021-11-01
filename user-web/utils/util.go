package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
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
			return response.NewBaseResponse(300, err.Error(), nil), false
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

func GetServerConnFromConsul(serverName string) (*grpc.ClientConn, error) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	var host string
	var port int

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", serverName))
	if err != nil {
		panic(err)
	}


	for _, s := range data {
		host = s.Address
		port = s.Port
		break
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func GetBoolEnvInfo(name string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(name)
}
