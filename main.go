package main

import (
	"os"
	"os/signal"
	"stock-web-be/dao/redis"
	"syscall"
	"time"

	"stock-web-be/async"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/preload"
	"stock-web-be/router"

	"github.com/stockAlpha/gopkg/common/safego"
)

func main() {
	preload.Init()
	tlog.Handler.Infof(nil, consts.SLTagSeverStart, "stock web be is Starting...")
	//支持多平台编译
	safego.SafeGo(func() {
		err := preload.Run(router.Register)
		if err != nil {
			tlog.Handler.Warnf(nil, consts.SLTagSeverFail, "stock web be preload Error: %v", err)
			panic(err)
		}
	})

	// 一些后台任务依赖优雅结束
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	tlog.Handler.Infof(nil, consts.SLTagSeverStop, "Quit the stock web be !")
	//释放资源
	redis.Close()
	// 停止async的所有线程
	async.PreStop()
	safego.SafeGoWait()
	time.Sleep(time.Second)
}
