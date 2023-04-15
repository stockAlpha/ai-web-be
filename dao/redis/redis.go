package redis

import (
	"context"

	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	return redisClient
}
func Init() {
	dsn := conf.Handler.GetString("redis.host") + ":" + conf.Handler.GetString("redis.port")
	redisOption := &redis.Options{
		Addr: dsn,
	}
	if conf.Handler.GetString("redis.password") != "" {
		redisOption.Password = conf.Handler.GetString("redis.password")
	}
	tlog.Handler.Infof(nil, consts.SLTagRedisSuccess, redisOption.Addr+","+redisOption.Password)
	redisClient = redis.NewClient(redisOption)
	defer redisClient.Close()
	if resp := redisClient.Ping(context.Background()); resp.Err() != nil {
		panic("redis init error," + redisOption.Addr + "," + redisOption.Password)
	}
	tlog.Handler.Infof(nil, consts.SLTagRedisSuccess, redisOption.Addr+","+redisOption.Password)
}
