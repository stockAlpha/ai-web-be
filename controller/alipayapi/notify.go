package alipayapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
)

// @Tags	alipay支付相关接口
// @Summary	异步通知
// @Router		/api/v1/alipay/notify [post]
func Notify(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	fmt.Println("req body", c.Request.Body)
	var req alipay.TradeNotification

	if err := c.ShouldBind(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	fmt.Println("req1:", req)
	fmt.Println("order:", req.OutTradeNo)
	c.String(http.StatusOK, "success")
}
