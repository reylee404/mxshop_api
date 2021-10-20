package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"mxshop_api/user-web/initialize"
)

func main() {

	port := flag.Int("port", 8080, "端口号")
	flag.Parse()

	initialize.MustInitLogger()

	engine := initialize.InitGinAndRouters()
	err := engine.Run(fmt.Sprintf(":%d", *port))
	if err != nil {
		zap.L().Panic("启动失败", zap.Error(err))
	}
}
