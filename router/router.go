package router

import (
	"github.com/gin-gonic/gin"
	"stock-web-be/controller/stockapi/auth"
	"stock-web-be/controller/stockapi/echo"
	"stock-web-be/gocommon/consts"
)

func Register(r *gin.Engine) *gin.Engine {
	stock := r.Group(consts.Prefix)

	stock.POST("/echo", echo.Echo)

	stock.POST(consts.SendVerificationCodeApi, auth.SendVerificationCode)
	stock.POST(consts.RegisterApi, auth.Register)
	stock.POST(consts.LoginApi, auth.Login)
	return r
}
