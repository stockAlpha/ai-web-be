package alipayapi

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"stock-web-be/client/alipayclient"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/payapi"
	"stock-web-be/idl/userapi/order"
	"strconv"
)

// @Tags	alipay支付相关接口
// @Summary	异步通知
// @Router		/api/v1/alipay/notify [post]
func Notify(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req payapi.NotifyRequest

	if err := c.ShouldBind(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	// 校验订单
	orderId := req.OutTradeNo
	amount := req.TotalAmount
	appId := req.AppId
	if appId != alipayclient.APPID {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "appId not match req: %v", req)
		c.String(http.StatusOK, "failed")
		return
	}
	parseOrderId, err := strconv.ParseUint(orderId, 10, 64)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "orderId：%s parse error: %s", orderId, err.Error())
		c.String(http.StatusOK, "failed")
		return
	}
	existOrder, err := order.GetOrderById(parseOrderId)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "get order error, error: %s", err.Error())
		c.String(http.StatusOK, "failed")
		return
	}
	if existOrder == nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "order not found, orderId: %s", orderId)
		c.String(http.StatusOK, "failed")
		return
	}
	decimalAmount, err := decimal.NewFromString(amount)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "parse order amount error, error: %s", err.Error())
		c.String(http.StatusOK, "failed")
		return
	}
	if decimalAmount != existOrder.Amount {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "order amount not match req: %v, error: %s", req, err.Error())
		c.String(http.StatusOK, "failed")
		return
	}

	// 修改订单状态
	status := req.TradeStatus
	if status == "TRADE_SUCCESS" || status == "TRADE_FINISHED" {
		err = order.UpdateOrderStatus(parseOrderId, 1, nil)
		if err != nil {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "update order status error, error: %s", err.Error())
			c.String(http.StatusOK, "failed")
			return
		}
	}

	c.String(http.StatusOK, "success")
}
