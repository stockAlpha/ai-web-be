package router

import (
	"github.com/gin-gonic/gin"
	"stock-web-be/controller/stockapi/auth"
	"stock-web-be/controller/stockapi/echo"
)

func Register(r *gin.Engine) *gin.Engine {
	stock := r.Group("/apis/stock/web/")

	stock.POST("/echo", echo.Echo)

	stock.POST("/verify/send_code", auth.SendVerificationCode)
	stock.POST("/user/register", auth.Register)
	stock.POST("/user/login", auth.Login)
	return r
}
