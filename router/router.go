package router

import (
	"github.com/gin-gonic/gin"
	"stock-web-be/controller/stockapi/echo"
)

func Register(r *gin.Engine) *gin.Engine {
	stock := r.Group("/apis/stock/web/")

	stock.POST("/echo", echo.Echo)
	return r
}
