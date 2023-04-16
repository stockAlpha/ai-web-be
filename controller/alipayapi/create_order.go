package alipayapi

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/payapi"
)

// @Tags	alipay支付相关接口
// @Summary	创建聚合收钱单
// @Router		/api/v1/alipay/create_order [post]
// @param		req	body		payapi.CreateOrderRequest	true	"请求参数"
// @Response 200 {object} payapi.CreateOrderResponse 返回信息
func CreateOrder(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req payapi.CreateOrderRequest

	if err := c.ShouldBindWith(&req, binding.Form); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	response := payapi.AlipayResponse{Data: payapi.CreateOrderResponse{
		MerchantName: "ChatAlpha",
		MerchantId:   "2021003189689338",
		Amount:       req.Amount,
		OutTradeNo:   req.OutTradeNo,
	},
		Code: "10000",
		Msg:  "Success",
	}
	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}
