package initialize

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"mxshop_api/user-web/global"
)

func InitRedis(db interface{}) {
	db = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisConfig.Host, global.ServerConfig.RedisConfig.Port),
	})
}
