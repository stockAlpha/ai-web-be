package payapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/order"
	"strconv"
)

// @Tags	支付相关接口
// @Summary	获取支付状态
// @Router		/api/v1/pay/status [get]
// @param		 orderId query	string			true	"订单id"
// @Success 200 {int} int "支付状态：1-待支付,2-已支付,3-已取消"
func Status(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	orderId, _ := strconv.ParseUint(c.GetString("orderId"), 10, 64)
	orderStatus, err := order.GetOrderById(orderId)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "find order error orderId=%s", orderId, err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrNotExistOrder)
		return
	}
	if orderStatus == nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "orderId=%d not found", orderId)
		cg.Res(http.StatusBadRequest, controller.ErrNotExistOrder)
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, orderStatus.Status)
}
