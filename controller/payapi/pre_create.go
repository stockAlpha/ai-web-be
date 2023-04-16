package payapi

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/smartwalle/alipay/v3"
	"net/http"
	"stock-web-be/client/alipayclient"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/payapi"
	"stock-web-be/logic/userapi"
	"stock-web-be/utils"
	"strconv"
)

// @Tags	支付相关接口
// @Summary	预创建交易订单
// @Router		/api/v1/pay/pre_create [post]
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
	//amount := utils.GetAmount(req.ProductType)
	// 测试阶段先都用0.01
	amount := 0.01
	orderId, err := userapi.AddOrder(userId, decimal.NewFromInt(int64(amount)), strconv.Itoa(amount)+"元积分套餐", nil)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "add order error", err.Error())
		cg.Resp(http.StatusBadRequest, controller.ErrCreateOrder, "创建订单失败，请重试或者联系客服")
		return
	}
	res, err := client.TradePreCreate(alipay.TradePreCreate{
		Trade: alipay.Trade{
			Subject:     "ChatAlpha积分充值",
			NotifyURL:   "https://web-be-test.stockalpha.top/api/v1/alipay/notify",
			OutTradeNo:  strconv.FormatUint(orderId, 10),
			TotalAmount: strconv.Itoa(amount),
			ProductCode: "FACE_TO_FACE_PAYMENT",
		},
	})
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "alipay pre create", err.Error())
		cg.Resp(http.StatusBadRequest, controller.ErrServer, res)
		return
	}
	c.JSON(http.StatusOK, res.Content)
	return
}
