package preload

import (
	"github.com/gin-contrib/pprof"
	"net/http"

	"stock-web-be/controller/public"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/middleware"

	"github.com/gin-gonic/gin"
)

type server struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
}

func Run(inner func(r *gin.Engine) *gin.Engine) error {
	s := &server{
		Address:      conf.Handler.GetString(`http.addr`),
		ReadTimeout:  conf.Handler.GetInt64(`http.readTimeout`),
		WriteTimeout: conf.Handler.GetInt64(`http.writeTimeout`),
	}

	tlog.Handler.Infof(nil, consts.SLTagSeverStart, "HttpServer starting from  ....%s", s.Address)
	appEngine := gin.New()
	appEngine.GET("ping", public.Ping)
	appEngine.GET("check/aalcsekudqkqgskctkjowuvhbiixmujd", public.Check)
	pprof.Register(appEngine, "pprof/aalcsekudqkqgskctkjowuvhbiixmujd")
	appEngine.Use(middleware.Recovery())
	appEngine.Use(middleware.GinLogger(middleware.LoggerConfig{}))
	appEngine.Use(middleware.ValidUser())

	//初始化路由
	appRouterConfig := inner(appEngine)
	err := appRouterConfig.Run(s.Address)
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
