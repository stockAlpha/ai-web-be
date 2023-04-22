package public

import (
	"context"
	"net/http"
	"stock-web-be/dao/db"
	"stock-web-be/dao/redis"
	"sync"
	"time"

	"stock-web-be/controller"

	"github.com/gin-gonic/gin"
)

var timeout = time.Second * 2

func Check(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, c.RemoteIP()+"\t\n"+CheckAll())
	return
}
func CheckAll() (resp string) {
	var wg sync.WaitGroup
	redisResp := ""
	mysqlResp := ""
	wg.Add(1)
	go checkRedis(&wg, &redisResp)
	wg.Add(1)
	go checkMysql(&wg, &mysqlResp)
	wg.Wait()
	resp += redisResp + "\t\n"
	resp += mysqlResp + "\t\n"
	return
}
func checkRedis(wg *sync.WaitGroup, res *string) {
	var resp string
	defer wg.Done()
	firstDate := time.Now()
	redisCLient := redis.GetRedisClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout))
	defer cancel()
	go func() {
		redisCLient.Ping(ctx)
		ctx.Done()
	}()
	select {
	case <-ctx.Done():
		resp = "redis done"
	case <-time.After(time.Duration(2 * timeout)):
		resp = "redis timeout"
	}
	*res = resp + " with time " + time.Now().Sub(firstDate).String()
}
func checkMysql(wg *sync.WaitGroup, res *string) {
	var resp string
	defer wg.Done()
	firstDate := time.Now()
	mysql, _ := db.DbIns.DB.DB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout))
	defer cancel()
	go func() {
		mysql.PingContext(ctx)
		ctx.Done()
	}()
	select {
	case <-ctx.Done():
		resp = "mysql done"
	case <-time.After(time.Duration(2 * timeout)):
		resp = "mysql timeout"
	}
	*res = resp + " with time " + time.Now().Sub(firstDate).String()
}
