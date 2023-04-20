package public

import (
	"net/http"

	"stock-web-be/controller"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, c.RemoteIP())
	return
}
