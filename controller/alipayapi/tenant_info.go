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
// @Summary	聚合收钱码商户信息查询
// @Router		/api/v1/alipay/tenant_info [post]
// @param		req	body		payapi.TenantInfoRequest	true	"请求参数"
// @Response 200 {object} payapi.TenantInfoResponse 返回信息
func TenantInfo(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req payapi.TenantInfoRequest

	if err := c.ShouldBindWith(&req, binding.Form); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	response := payapi.AlipayResponse{Data: payapi.TenantInfoResponse{
		MerchantName: "ChatAlpha",
		MerchantId:   "2021003189689338",
		MerchantLogo: "https://chatalpha.top/logo.svg",
		AlipayAppId:  "2021003189689338",
		AuthRedirect: "https://web-be.stockalpha.top/api/alipay/callback",
	},
		Code: "10000",
		Msg:  "Success",
	}
	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}
