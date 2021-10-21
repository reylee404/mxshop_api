package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	SignKey string `mapstructure:"key"`
}

type ServerConfig struct {
	Name          string        `mapstructure:"name"`
	Port          int           `mapstructure:"port"`
	UserSrvConfig UserSrvConfig `mapstructure:"user_srv"`
	JWTConfig     JWTConfig     `mapstructure:"jwt"`
}
