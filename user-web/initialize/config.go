package initialize

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func getBoolEnvInfo(name string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(name)
}

// MustInitConfig 第一次初始化 config 时必须成功，否则 panic;
// 监听到变化后导致重新解析 config 失败不 panic，只打印错误日志
func MustInitConfig(config interface{}) {
	v, err := initConfig(config)
	if err != nil {
		panic(err)
	}
	watchConfigFile(v, func(in fsnotify.Event) {
		err = readAndUnmarshalConfig(v, config)
		if err != nil {
			zap.L().Error("readAndUnmarshalConfig failed", zap.Error(err))
		}
	})
}

func initConfig(config interface{}) (*viper.Viper, error) {
	dev := getBoolEnvInfo("MX_SHOP_DEV")
	configFileName := "./config_pro.yaml"
	if dev {
		configFileName = "./config_dev.yaml"
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	err := readAndUnmarshalConfig(v, config)
	if err != nil {
		return nil, err
	}
	return v, nil

}

func watchConfigFile(v *viper.Viper, run func(in fsnotify.Event)) {
	v.WatchConfig()
	v.OnConfigChange(run)
}

func readAndUnmarshalConfig(v *viper.Viper, config interface{}) error {
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	if err := v.Unmarshal(config); err != nil {
		return err
	}
	return nil
}
