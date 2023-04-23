package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
)

var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	return redisClient
}
func Init() {
	dsn := conf.Handler.GetString("redis.uri")
	fmt.Println("dsn", dsn)
	redisOption := &redis.Options{
		Addr: dsn,
	}
	if conf.Handler.GetString("redis.password") != "" {
		redisOption.Password = conf.Handler.GetString("redis.password")
	}
	if conf.Handler.GetString("redis.default_db") != "" {
		redisOption.DB = conf.Handler.GetInt("redis.default_db")
	}

	tlog.Handler.Infof(nil, consts.SLTagRedisSuccess, redisOption.Addr+","+redisOption.Password)
	redisClient = redis.NewClient(redisOption)
	if resp := redisClient.Ping(context.Background()); resp.Err() != nil {
		panic("redis init error," + redisOption.Addr + "," + redisOption.Password)
	}

	tlog.Handler.Infof(nil, consts.SLTagRedisSuccess, redisOption.Addr+","+redisOption.Password)
}

func Close() {
	_ = redisClient.Close()
}
