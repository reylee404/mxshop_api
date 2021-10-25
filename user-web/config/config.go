package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	SignKey string `mapstructure:"key"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name          string        `mapstructure:"name"`
	Port          int           `mapstructure:"port"`
	UserSrvConfig UserSrvConfig `mapstructure:"user_srv"`
	JWTConfig     JWTConfig     `mapstructure:"jwt"`
	RedisConfig   RedisConfig   `mapstructure:"redis"`
}
