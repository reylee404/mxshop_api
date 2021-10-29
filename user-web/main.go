package main

import (
	"fmt"
	"mxshop_api/user-web/global"
	"mxshop_api/user-web/utils"

	"go.uber.org/zap"

	"mxshop_api/user-web/initialize"
)

func main() {
	initialize.MustInitLogger()
	initialize.MustInitConfigFromNacos(initialize.NacosManagerConfig{
		IP:        "127.0.0.1",
		Port:      8848,
		NameSpace: "d4842894-685a-42ae-b012-0776520abaab",
		DataId:    "user-web",
		Group:     "dev",
		Listen: true,
	}, global.ServerConfig)
	global.Trans = initialize.MustInitTrans("zh")
	initialize.MustInitValidators(&global.Trans)
	initialize.InitRedis(global.RedisClient)
	_ = initialize.InitSrvConn()

	engine := initialize.InitGinAndRouters()

	if !global.Dev {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	addr := fmt.Sprintf(":%d", global.ServerConfig.Port)
	zap.L().Info("starting gin web", zap.String("name", global.ServerConfig.Name), zap.String("addr", addr))
	err := engine.Run(addr)
	if err != nil {
		zap.L().Panic("start Failed", zap.Error(err))
	}
}
