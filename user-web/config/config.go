package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SignKey string `mapstructure:"key" json:"key"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name          string        `mapstructure:"name" json:"name"`
	Port          int           `mapstructure:"port" json:"port"`
	UserSrvConfig UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTConfig     JWTConfig     `mapstructure:"jwt" json:"jwt"`
	RedisConfig   RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulConfig  ConsulConfig  `mapstructure:"consul" json:"consul"`
}

type NacosConfig struct {
	IP        string `mapstructure:"ip"`
	Port      uint64 `mapstructure:"port"`
	NameSpace string `mapstructure:"name_space"`
	DataId    string `mapstructure:"data_id"`
	Group     string `mapstructure:"group"`
	Listen    bool   `mapstructure:"listen"`
}
