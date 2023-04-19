package alipayapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"os"
	"stock-web-be/async"
	"stock-web-be/controller"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/payapi"
	"stock-web-be/idl/userapi/order"
	"stock-web-be/logic/userapi"
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
	if appId != conf.Handler.GetString("alipay.app_id") {
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
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "orderId=%s not found order", orderId)
		c.String(http.StatusOK, "failed")
		return
	}
	decimalAmount, err := decimal.NewFromString(amount)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "parse order amount error, error: %s", err.Error())
		c.String(http.StatusOK, "failed")
		return
	}
	// 线上环境校验金额是否匹配
	if os.Getenv(consts.Env) == "prod" {
		if !decimalAmount.Equal(existOrder.Amount) {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "order amount not match req: %v, error: %s", req, err.Error())
			c.String(http.StatusOK, "failed")
			return
		}
	}

	tx := db.DbIns.Begin()

	// 修改订单状态,充值积分
	status := req.TradeStatus
	if status == "TRADE_SUCCESS" || status == "TRADE_FINISHED" {
		err = order.UpdateOrderStatus(parseOrderId, 1, tx)
		if err != nil {
			tx.Rollback()
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "update order status error, error: %s", err.Error())
			c.String(http.StatusOK, "failed")
			return
		}
		integralAmount, err := strconv.Atoi(existOrder.ProductInfo)
		if err != nil {
			tx.Rollback()
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "parse integral amount=%s error, error: %s", existOrder.ProductInfo, err.Error())
			c.String(http.StatusOK, "failed")
			return
		}
		err = userapi.AddUserIntegral(existOrder.FromUserId, integralAmount, tx)
		if err != nil {
			tx.Rollback()
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "add user integral error, error: %s", err.Error())
			c.String(http.StatusOK, "failed")
			return
		}
		userId := existOrder.FromUserId
		user, err := userapi.GetUserById(userId)
		if err != nil {
			tx.Rollback()
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "userId=%d not found user, error: %s", userId, err.Error())
			c.String(http.StatusOK, "failed")
			return
		}
		if user == nil {
			tx.Rollback()
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "userId=%d not found user", userId)
			c.String(http.StatusOK, "failed")
			return
		}
		tx.Commit()
		async.MailChan <- async.MailChanType{To: user.Email, Subject: consts.RechargeNotifySubject, Body: fmt.Sprintf(consts.RechargeNotifyContent, integralAmount)}
	}

	c.String(http.StatusOK, "success")
}