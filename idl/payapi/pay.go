package payapi

type PreCreateRequest struct {
	PayType     string `json:"payType" default:"alipay"`       // 支付类型，目前只支持alipay
	ProductType int    `json:"productType" binding:"required"` // 商品类型，1-10元,2-30元,3-100元
}
