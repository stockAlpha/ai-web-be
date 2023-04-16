package alipayapi

type AlipayResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type TenantInfoRequest struct {
	QrCodeId string `json:"qr_code_id"` // 聚合收钱码的码值
	Ua       string `json:"ua"`         // 扫码客户端userAgent
}

type TenantInfoResponse struct {
	MerchantName string `json:"merchant_name"` // 商户名称
	MerchantLogo string `json:"merchant_logo"` // 商家logo图片资源地址
	MerchantId   string `json:"merchant_id"`   // 商户在支付宝的标识
	AlipayAppId  string `json:"alipay_app_id"` // 服务商主体在支付宝的标识
	AuthRedirect string `json:"auth_redirect"` // 支付宝认证回调地址，为服务商入驻支付宝时与appid绑定的地址
}

type CreateOrderRequest struct {
	QrCodeId   string `json:"qr_code_id"`   // 需创单的聚合收钱码码值
	Amount     string `json:"amount"`       // 用户支付的金额，精确到小数点后两位
	AuthCode   string `json:"auth_code"`    // 用户鉴权token，标识用户该次支付身份
	OutTradeNo string `json:"out_trade_no"` // 外部交易号，用于重复发起创单时的幂等
	Remark     string `json:"remark"`       // 用户支付时的备注信息
	Ua         string `json:"ua"`           // 扫码客户端userAgent
}

type CreateOrderResponse struct {
	MerchantName  string `json:"merchant_name"`   // 聚合收钱码所属商户名称
	AlipayTradeNo string `json:"alipay_trade_no"` // 该笔交易对应的支付宝交易号
	OutTradeNo    string `json:"out_trade_no"`    // 创单时传入的外部幂等交易号
	OrderId       string `json:"order_id"`        // 服务商商户订单号，与支付宝交易号相对应
	MerchantId    string `json:"merchant_id"`     // 商户在支付宝的标识
	Amount        string `json:"amount"`          // 实际交易金额，以元为单位
}
