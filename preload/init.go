package preload

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/logic/stockapi/notify"
	"strconv"
)

/*资源预加载*/
//var LogPath = "./log"
//var ConfPath = "./conf"

func Init() {
	conf.Init("")
	//logic类init start
	notify.Init()
	//logic类init end
	tlog.Init()
	initGOProcs()
	initPProf()
	db.InitDB()
}

func initPProf() {
	if addr := conf.Handler.GetString("pprof.port"); addr != "" {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					tlog.Handler.Panicf(context.TODO(), consts.SLTagPanic, "errmsg=%s||stack info=%s", err, strconv.Quote(string(debug.Stack())))
				}
			}()

			if err := http.ListenAndServe(addr, nil); err != nil {
				tlog.Handler.Warnf(context.TODO(), consts.SLTagPprofFail, "PProf ListenAndServe Error:%v", err)
			}
		}()
	}
}

func initGOProcs() {
	cpus := os.Getenv("CPUS")
	cpusNum, _ := strconv.Atoi(cpus)
	if cpusNum > 0 {
		runtime.GOMAXPROCS(cpusNum)
		tlog.Handler.Infof(context.TODO(), consts.SLTagSeverStart, "CPUS:%v", cpusNum)
	}
}
