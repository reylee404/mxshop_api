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
	global.Dev = utils.GetBoolEnvInfo("MX_SHOP_DEV")
	configFileName := "./config_pro.yaml"
	if global.Dev {
		configFileName = "./config_dev.yaml"
	}
	initialize.MustInitConfig(configFileName, global.NacosConfig)
	initialize.MustInitConfigFromNacos(initialize.NacosManagerConfig{
		IP:        global.NacosConfig.IP,
		Port:      global.NacosConfig.Port,
		NameSpace: global.NacosConfig.NameSpace,
		DataId:    global.NacosConfig.DataId,
		Group:     global.NacosConfig.Group,
		Listen:    global.NacosConfig.Listen,
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
