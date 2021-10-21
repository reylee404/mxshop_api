package initialize

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop_api/user-web/config"
	"mxshop_api/user-web/global"
)

func getBoolEnvInfo(name string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(name)
}

func InitConfig() {
	dev := getBoolEnvInfo("MX_SHOP_DEV")
	configFileName := "./config_pro.yaml"
	if dev {
		configFileName = "./config_dev.yaml"
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	readAndUnmarshalConfig(v, global.ServerConfig)

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Info("config changed", zap.String("name", in.Name))
		readAndUnmarshalConfig(v, global.ServerConfig)
	})

}

func readAndUnmarshalConfig(v *viper.Viper, serverConfig *config.ServerConfig) {
	if err := v.ReadInConfig(); err != nil {
		zap.L().Error("config not found", zap.Error(err))
		return
	}
	if err := v.Unmarshal(serverConfig); err != nil {
		zap.L().Error("config unmarshal failed", zap.Error(err))
		return
	}
}
