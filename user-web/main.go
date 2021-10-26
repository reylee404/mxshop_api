package main

import (
	"fmt"
	"mxshop_api/user-web/global"

	"go.uber.org/zap"

	"mxshop_api/user-web/initialize"
)

func main() {

	initialize.MustInitConfig(global.ServerConfig)
	initialize.MustInitLogger()
	global.Trans = initialize.MustInitTrans("zh")
	initialize.MustInitValidators(&global.Trans)
	initialize.InitRedis(global.RedisClient)

	engine := initialize.InitGinAndRouters()

	addr := fmt.Sprintf(":%d", global.ServerConfig.Port)
	zap.L().Info("starting gin web", zap.String("name", global.ServerConfig.Name), zap.String("addr", addr))
	err := engine.Run(addr)
	if err != nil {
		zap.L().Panic("start Failed", zap.Error(err))
	}
}
