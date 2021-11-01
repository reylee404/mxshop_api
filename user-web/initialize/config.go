package initialize

import (
	"encoding/json"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"mxshop_api/user-web/global"
)

type NacosManagerConfig struct {
	IP        string
	Port      uint64
	NameSpace string
	DataId    string
	Group     string
	Listen    bool
}

type NacosConfigManager struct {
	ip        string
	port      uint64
	namespace string
	dataId    string
	group     string
	client    config_client.IConfigClient
	listen    bool
}

func MustInitConfigFromNacos(nConfig NacosManagerConfig, config interface{}) {
	err := InitConfigFromNacos(nConfig, config)
	if err != nil {
		panic(err)
	}
}

func InitConfigFromNacos(nConfig NacosManagerConfig, config interface{}) error {
	manager := NewNacosConfigManager(nConfig)
	err := manager.GetConfig(config)
	if err != nil {
		return err
	}
	if nConfig.Listen {
		err := manager.ListenConfig(config)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewNacosConfigManager(config NacosManagerConfig) *NacosConfigManager {
	return &NacosConfigManager{
		ip:        config.IP,
		port:      config.Port,
		namespace: config.NameSpace,
		dataId:    config.DataId,
		group:     config.Group,
		listen:    config.Listen,
	}
}

func (n *NacosConfigManager) GetConfig(config interface{}) error {
	err := n.initConfigClient()
	if err != nil {
		return err
	}

	err = n.getConfigAndUnmarshal(config)
	if err != nil {
		return err
	}

	return nil
}

func (n *NacosConfigManager) ListenConfig(config interface{}) error {
	err := n.client.ListenConfig(vo.ConfigParam{
		DataId: n.dataId,
		Group:  n.group,
		OnChange: func(namespace, group, dataId, data string) {
			err := n.UnmarshalConfig([]byte(data), config)
			if err != nil {
				zap.L().Error("UnmarshalConfig", zap.Error(err))
			}
			zap.L().Info("ListenConfig", zap.String("name", global.ServerConfig.Name))
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *NacosConfigManager) initConfigClient() (err error) {
	sc := []constant.ServerConfig{
		{
			IpAddr: n.ip,
			Port:   n.port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         n.namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	n.client, err = clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &cc,
		ServerConfigs: sc,
	})
	if err != nil {
		return err
	}

	return nil
}

func (n *NacosConfigManager) getConfigAndUnmarshal(config interface{}) error {
	configContent, err := n.client.GetConfig(vo.ConfigParam{
		DataId: n.dataId,
		Group:  n.group,
	})
	if err != nil {
		return err
	}

	err = n.UnmarshalConfig([]byte(configContent), config)
	if err != nil {
		return err
	}
	return nil
}

func (n *NacosConfigManager) UnmarshalConfig(data []byte, config interface{}) error {
	err := json.Unmarshal(data, config)
	if err != nil {
		return err
	}
	return nil
}

// MustInitConfig 第一次初始化 config 时必须成功，否则 panic;
// 监听到变化后导致重新解析 config 失败不 panic，只打印错误日志
func MustInitConfig(configFileName string, config interface{}) {
	v, err := initConfig(configFileName, config)
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

func initConfig(configFileName string, config interface{}) (*viper.Viper, error) {
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
