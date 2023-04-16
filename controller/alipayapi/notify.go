package alipayapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/smartwalle/alipay/v3"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/payapi"
)

// @Tags	alipay支付相关接口
// @Summary	通知
// @Router		/api/v1/alipay/notify [post]
func Notify(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req alipay.TradeNotification

	if err := c.ShouldBindWith(&req, binding.Form); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	fmt.Println("req1:", req)

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	fmt.Println("req2:", req)
	response := payapi.AlipayResponse{
		Code: "10000",
		Msg:  "Success",
	}
	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}
