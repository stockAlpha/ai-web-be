package client

import (
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/pkg/xlog"
	"os"
)

func Init() {
	// 初始化支付宝客户端
	//    appid：应用ID
	//    privateKey：应用私钥，支持PKCS1和PKCS8
	//    isProd：是否是正式环境
	isProd := false
	if os.Getenv("ENV") == "prod" {
		isProd = true
	}
	client, err := alipay.NewClient("2016091200494382", "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDNnHMn3qgthmg1aUeXM9/KP8+cOF44lgZL5YC3295Ur900VPpN3EmNH9xZGb5XKDtLA027IEfpSCIAGzEoYrXeeguD28eR29eF6H1WwT3ys/CgtnGPf5O6bNX9IvEOB3bg2djperRJhXll5lJFElnp24E25hWg/PtTVFJADUNXd6MMx0b465tH0ys4TXslq4CajYH5nZ0JmU+sDDcnDfz+PMw/yQhNbCihQf/Uclw5IMA4p05EIE8SKipkWigEKUNQJNSIuCLSw30+y4KAdi44pe7e0GSBsA4RnUH9bzQ5vQc80RcS9jrs+ilGtqqyo8Yx6HYiA2JYXCzJvspCYM5BAgMBAAECggEACRn0yv4DKci6Uirv3VSRgm1irkKcgRq6+p8VHI5iABOs0gU08z9eDh4r7dHC6C7CuSZTSYY7SKtFvgV2HpiH/moeml6fLeiyWJ4a9j0lldm7PBH3Yue9zVHGAyeJzhose/WIsYUZ9+stnYIdgCs4ya5KwPhcWrz4Dw87eNRdd8CipFYQAuGIPLLXuTHSNRdDAQ99Ww+q9vgFoerMbcik1AbWF9JeZ9ViCrp8I7gnixmSvNuNdNuAFAguGdGHsHv4jJC3NgQfHFL6g6b2K9MeQXpOzM8DKB4NE0YL7P2M7Wz3BnRdWiRDJ8d1//pD8byQvjjYrAChXRkDFulN5s948QKBgQDwstbXjkNqHGuRak968Cs86YwaNwb+K676NvXG0UcsMJTnLBikM1x3TlkvC+EQ4np4hfR9mfuuhHz5z2V7+0hN3LKlDeRWKeKFDAz+xK68EhvQEjjCcOKnQwR0CVk1QMO/NtX1TyCOHSENk5qfJZMPXL14cvnRX1E8VI1NnbR+bQKBgQDarpe3yGn391U1YMjXcoaFfnxLYnysk7O+vVKXdvNTWO/N5gZ/ndSUSAI97M978HN3DD3zTxPdv7Ba2f8KR5AffyjI9ldzUOPlWKXSSnJILm6QLy3UcTx8T2yqCMgVH8+tjozuL1oEWQPtPmz5G4lkQPOroPrMCHPYuf1wR11apQKBgQDEduU/2qISGZJ/hgvL5/8S3/p4Z1Pw4M0Y9QVVu/phCmJv8qFGXZnq0+udqA+UDZgziftPDgHNp9yutuc59JhG3Y5/hMBMyDFZscVlqjqJziofgtALfcKzDdOztvG2st9T0zl+2pBTD1msUD+UCUJo9qS8jPR2Plv7Z3RS1xe9XQKBgAjqdWUY1rk1bFPwzj96e+GVdpvcOBkRLJLqRSPHxcPwLFbIuhsZ5EDnbq/3p731379K0HvLDZRM7HPHz44rvMSL+q223XnmImSHaLtaLa6jtf5K7iNrOnwXAOct1HqIAX+iADz10UW8G1zg3rCJXuCnhUfKGG+ZKJ/9dUfOoQ/BAoGBAKS+7NjoNOqWoVqic2GASBZGYo/omASYtk1liHDYJLX+aSPVCH5E/7WEGTsGeZ5mAM0un5yxgLmTZniHLr0QsFUcFUCQCsk56GwNBBKIt22TXqesRF8L5UNy7h2TZ++xpvQ0Zl7nVs3RTSrvaHLYdinxI9Cphj+ggZn3/76MBIbr", isProd)
	if err != nil {
		xlog.Error(err)
		return
	}

	// 打开Debug开关，输出日志，默认关闭
	client.DebugSwitch = gopay.DebugOn

	// 设置支付宝请求 公共参数
	//    注意：具体设置哪些参数，根据不同的方法而不同，此处列举出所有设置参数
	client.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
							SetCharset(alipay.UTF8).                                            // 设置字符编码，不设置默认 utf-8
							SetSignType(alipay.RSA2).                                           // 设置签名类型，不设置默认 RSA2
							SetReturnUrl("https://chatalpha.top").                              // 设置返回URL
							SetNotifyUrl("https://web-be-test.stockalpha.top/api/pay/callback") // 设置异步通知URL

	// 自动同步验签（只支持证书模式）
	// 传入 alipayCertPublicKey_RSA2.crt 内容
	client.AutoVerifySign([]byte("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoCc4rE2TC+wvqkHRn11WQmRnTJRgiP+WR7NmCH1vSuMlG1jyzVd3/bMBrGM1Q+gkJx0/1Qo2L5SpaVO9LHbKR2yOpUevhwRXzgTskMZAp4itedRmIVlJt5DoERmRYQeiI4gAiPfG8LowR1TuAWxoUcguW422zmlpuZbGEqhmzcPNuJ+ImtziefkEK5X4mOwkDfOslU/HGiPTn0EaRfZZOx0CJbHTFftcvJ/ibeg9jHV2x0Q+RVuxJGeOe8VYkoV24Jp7EvVQm4RQq922Ij0Y0JkXPSYMjavBU+Uc5fhvRndvbS+k70Pn0Xrz9cAtlarv/p23QeLNQeZ8gZbHyj23LwIDAQAB"))
}
