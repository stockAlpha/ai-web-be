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
// @Summary	异步通知
// @Router		/api/v1/alipay/notify [post]
func Notify(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	fmt.Println("req", c.Request)
	fmt.Println("req body", c.Request.Body)
	fmt.Println("req form", c.Request.Form)
	fmt.Println("req url", c.Request.URL)
	fmt.Println("req postForm", c.Request.PostForm)
	var req alipay.TradeNotification

	if err := c.ShouldBindWith(&req, binding.FormPost); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	fmt.Println("req1:", req)

	response := payapi.AlipayResponse{
		Code: "10000",
		Msg:  "Success",
	}
	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}
