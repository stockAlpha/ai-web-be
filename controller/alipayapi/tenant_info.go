package alipayapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/idl/alipayapi"
)

// @Tags	支付相关接口
// @Summary	聚合收钱码商户信息查询
// @Router		/api/v1/alipay/tenant_info [post]
// @Response 200 {object} alipayapi.TenantInfoResponse 商户信息
func TenantInfo(c *gin.Context) {
	response := alipayapi.AlipayResponse{Data: alipayapi.TenantInfoResponse{
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
