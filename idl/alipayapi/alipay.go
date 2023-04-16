package alipayapi

type TenantInfoResponse struct {
	MerchantName string `json:"merchant_name"`
	MerchantLogo string `json:"merchant_logo"`
	MerchantId   string `json:"merchant_id"`
	AlipayAppId  string `json:"alipay_app_id"`
	AuthRedirect string `json:"auth_redirect"`
	Code         string `json:"code"`
	Msg          string `json:"msg"`
}
