package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"

	"mxshop_api/user-web/global"
	"mxshop_api/user-web/proto"
)

func InitSrvConn() error {
	c := global.ServerConfig.ConsulConfig
	conn, err := grpc.Dial(fmt.Sprintf("consul://%s:%d/%s?wait=14s", c.Host, c.Port, global.ServerConfig.UserSrvConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		return err
	}

	client := proto.NewUserClient(conn)
	global.UserClient = client
	return nil
}
