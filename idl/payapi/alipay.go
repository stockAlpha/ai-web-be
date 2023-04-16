package payapi

type AlipayResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// auth_app_id：授权应用的AppId
// app_id：当前接口调用方AppId
// user_id：支付宝用户的唯一标识
// authorized_user_id：授权者（即认证机构）的支付宝用户id
// state：商户自定义参数
// scope：授权范围
// auth_code：授权码，可用于换取access_token
// result_code：认证结果，包括“SUCCESS”和“FAIL”
// result_message：认证结果描述信息
// gmt_create：认证时间
// gmt_expire：认证过期时间
type CallBackRequest struct {
	AuthAppId        string `json:"auth_app_id"`        // 授权应用的AppId
	AppId            string `json:"app_id"`             // 当前接口调用方AppId
	UserId           string `json:"user_id"`            // 支付宝用户的唯一标识
	AuthorizedUserId string `json:"authorized_user_id"` // 授权者（即认证机构）的支付宝用户id
	State            string `json:"state"`              // 商户自定义参数
	Scope            string `json:"scope"`              // 授权范围
	AuthCode         string `json:"auth_code"`          // 授权码，可用于换取access_token
	ResultCode       string `json:"result_code"`        // 认证结果，包括“SUCCESS”和“FAIL”
	ResultMessage    string `json:"result_message"`     // 认证结果描述信息
	GmtCreate        string `json:"gmt_create"`         // 认证时间
	GmtExpire        string `json:"gmt_expire"`         // 认证过期时间
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

type NotifyRequest struct {
	AuthAppId           string `form:"auth_app_id"`           // App Id
	NotifyTime          string `form:"notify_time"`           // 通知时间
	NotifyType          string `form:"notify_type"`           // 通知类型
	NotifyId            string `form:"notify_id"`             // 通知校验ID
	AppId               string `form:"app_id"`                // 开发者的app_id
	Charset             string `form:"charset"`               // 编码格式
	Version             string `form:"version"`               // 接口版本
	SignType            string `form:"sign_type"`             // 签名类型
	Sign                string `form:"sign"`                  // 签名
	TradeNo             string `form:"trade_no"`              // 支付宝交易号
	OutTradeNo          string `form:"out_trade_no"`          // 商户订单号
	OutBizNo            string `form:"out_biz_no"`            // 商户业务号
	BuyerId             string `form:"buyer_id"`              // 买家支付宝用户号
	BuyerLogonId        string `form:"buyer_logon_id"`        // 买家支付宝账号
	SellerId            string `form:"seller_id"`             // 卖家支付宝用户号
	SellerEmail         string `form:"seller_email"`          // 卖家支付宝账号
	TradeStatus         string `form:"trade_status"`          // 交易状态
	TotalAmount         string `form:"total_amount"`          // 订单金额
	ReceiptAmount       string `form:"receipt_amount"`        // 实收金额
	InvoiceAmount       string `form:"invoice_amount"`        // 开票金额
	BuyerPayAmount      string `form:"buyer_pay_amount"`      // 付款金额
	PointAmount         string `form:"point_amount"`          // 集分宝金额
	RefundFee           string `form:"refund_fee"`            // 总退款金额
	Subject             string `form:"subject"`               // 商品的标题/交易标题/订单标题/订单关键字等，是请求时对应的参数，原样通知回来。
	Body                string `form:"body"`                  // 商品描述
	GmtCreate           string `form:"gmt_create"`            // 交易创建时间
	GmtPayment          string `form:"gmt_payment"`           // 交易付款时间
	GmtRefund           string `form:"gmt_refund"`            // 交易退款时间
	GmtClose            string `form:"gmt_close"`             // 交易结束时间
	FundBillList        string `form:"fund_bill_list"`        // 支付金额信息
	PassbackParams      string `form:"passback_params"`       // 回传参数
	VoucherDetailList   string `form:"voucher_detail_list"`   // 优惠券信息
	AgreementNo         string `form:"agreement_no"`          //支付宝签约号
	ExternalAgreementNo string `form:"external_agreement_no"` // 商户自定义签约号
}
