package alipayclient

import (
	"github.com/smartwalle/alipay/v3"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
)

var client *alipay.Client

func GetAlipayClient() *alipay.Client {
	return client
}

func Init() {
	// 创建支付宝客户端实例
	client, _ = alipay.New(conf.Handler.GetString("alipay.app_id"),
		conf.Handler.GetString("alipay.app_private_key"),
		true)

	client.LoadAppPublicCert(conf.Handler.GetString("alipay.app_public_cert"))
	client.LoadAliPayPublicKey(conf.Handler.GetString("alipay.alipay_public_key"))

	// 构造交易参数
	client.SetEncryptKey(conf.Handler.GetString("alipay.encrypt_key"))
	tlog.Handler.Infof(nil, consts.SLTagAlipaySuccess, "alipay connect success")
}
