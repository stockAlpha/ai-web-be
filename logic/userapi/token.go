package userapi

import (
	"context"
	"stock-web-be/dao/redis"
	"stock-web-be/utils"
	"time"
)

func AddTokenToBlackList(token string, expiration time.Duration) error {
	//1.避免token过长将token进行hash
	hashToken, err := utils.Md5(token)
	if err != nil {
		return err
	}

	//2.对token进行setNx,并设置过期时间
	redisClient := redis.GetRedisClient()
	if err = redisClient.SetNX(context.Background(), hashToken, nil, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func IsOnBlackList(token string) (bool, error) {
	//1.对token进行hash
	hashToken, err := utils.Md5(token)
	if err != nil {
		return false, err
	}

	//2.判断token是否存在在黑名单中
	redisClient := redis.GetRedisClient()
	exist, err := redisClient.Exists(context.Background(), hashToken).Result()
	if err != nil {
		return false, err
	}
	return exist == 1, nil
}
