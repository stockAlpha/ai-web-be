package router

import (
	"github.com/gin-gonic/gin"
	"stock-web-be/controller/userapi/auth"
	"stock-web-be/controller/userapi/echo"
	"stock-web-be/gocommon/consts"
)

func Register(r *gin.Engine) *gin.Engine {
	stock := r.Group(consts.ChatPrefix)

	stock.POST("/echo", echo.Echo)

	stock.POST(consts.SendVerificationCodeApi, auth.SendVerificationCode)
	stock.POST(consts.RegisterApi, auth.Register)
	stock.POST(consts.LoginApi, auth.Login)
	return r
}
