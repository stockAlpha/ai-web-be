package public

import (
	"log"
	"stock-web-be/dao/db"
	"stock-web-be/dao/redis"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/tlog"
	"testing"
)

func TestCheck(t *testing.T) {
	conf.Init("../../conf/app.local.toml")
	tlog.Init()
	db.InitDB()
	redis.Init()
	log.Println(CheckAll())
}
