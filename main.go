package main

import (
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/preload"
	"stock-web-be/router"
)

func main() {
	preload.Init()
	tlog.Handler.Infof(nil, consts.SLTagSeverStart, "stock web be is Starting...")
	//支持多平台编译
	err := preload.Run(router.Register)
	if err != nil {
		tlog.Handler.Warnf(nil, consts.SLTagSeverFail, "stock web be preload Error: %v", err)
	}
	tlog.Handler.Infof(nil, consts.SLTagSeverStart, "Quit the stock web be !")
}
