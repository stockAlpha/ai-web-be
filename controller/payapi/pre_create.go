package payapi

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/smartwalle/alipay/v3"
	"net/http"
	"os"
	"stock-web-be/client/alipayclient"
	"stock-web-be/controller"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/payapi"
	"stock-web-be/idl/userapi/order"
	"stock-web-be/utils"
	"strconv"
)

// @Tags	支付相关接口
// @Summary	预创建交易订单
// @Router		/api/v1/pay/pre_create [post]
// @param		req	body		payapi.PreCreateRequest	true	"订单请求参数"
// @Success 200 {object} payapi.PreCreateResponse "创建订单返回参数"
func PreCreate(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req payapi.PreCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	client := alipayclient.GetAlipayClient()
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	amount := 0
	switch req.ProductType {
	case 1:
		amount = 10
	case 2:
		amount = 30
	case 3:
		amount = 100
	default:
		amount = 10
	}
	orderId, err := order.AddOrder(userId, decimal.NewFromInt(int64(amount)), strconv.Itoa(utils.GetAmount(req.ProductType)), nil)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "add order error", err.Error())
		cg.Resp(http.StatusBadRequest, controller.ErrCreateOrder, "创建订单失败，请重试或者联系客服")
		return
	}
	totalAmount := "0.01"
	if os.Getenv(consts.Env) == "prod" {
		// 线上环境走真实金额
		totalAmount = strconv.Itoa(amount)
	}
	res, err := client.TradePreCreate(alipay.TradePreCreate{
		Trade: alipay.Trade{
			Subject:     "ChatAlpha积分充值",
			NotifyURL:   conf.Handler.GetString("alipay.notify_url"),
			OutTradeNo:  strconv.FormatUint(orderId, 10),
			TotalAmount: totalAmount,
			ProductCode: "FACE_TO_FACE_PAYMENT",
		},
	})
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "alipay pre create", err.Error())
		cg.Resp(http.StatusBadRequest, controller.ErrServer, res)
		return
	}
	ret := &payapi.PreCreateResponse{
		OrderId: res.Content.OutTradeNo,
		QRCode:  res.Content.QRCode,
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, ret)
	return
}
