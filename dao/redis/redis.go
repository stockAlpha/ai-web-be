package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"strings"
)

var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	return redisClient
}
func Init() {

	str := conf.Handler.GetString("redis.uri")
	// 分割字符串
	splitStr := strings.Split(str, "@")
	// 获取密码和地址
	password := strings.TrimPrefix(splitStr[0], "redis://:")
	address := splitStr[1]
	redisOption := &redis.Options{
		Addr:     address,
		Password: password,
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
