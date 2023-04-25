package router

import (
	"net/http"
	"stock-web-be/controller/alipayapi"
	"stock-web-be/controller/openaiapi/chat"
	"stock-web-be/controller/payapi"
	"stock-web-be/controller/userapi/integral"
	"stock-web-be/controller/userapi/user"
	"stock-web-be/docs"
	"stock-web-be/gocommon/consts"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}
func Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	//c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
}

func Register(r *gin.Engine) *gin.Engine {
	// 将 public 目录下的所有静态文件提供给客户端进行访问
	r.Static("/public", "./disk/public")
	r.Use(Options)
	r.Use(Secure)
	swagger(r)
	registerUser(r.Group(consts.UserPrefix))
	registerIntegral(r.Group(consts.IntegralPrefix))
	registerOpenAI(r.Group(consts.OpenaiPrefix))
	registerAlipay(r.Group(consts.AlipayPrefix))
	registerPay(r.Group(consts.PayPrefix))
	return r
}

func swagger(r *gin.Engine) {
	docs.SwaggerInfo.Title = "Stock Web API"
	docs.SwaggerInfo.Description = "This is stock web server api."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func registerUser(group *gin.RouterGroup) {
	group.POST(consts.SendVerificationCodeApi, user.SendVerificationCode)
	group.POST(consts.RegisterApi, user.Register)
	group.POST(consts.LoginApi, user.Login)
	group.POST(consts.LogoutApi, user.Logout)

	group.GET(consts.ProfileApi, user.Profile)
	group.POST(consts.ProfileApi, user.UpdateProfile)
	group.POST(consts.FeedbackApi, user.Feedback)
	group.POST(consts.SendPasswordVerificationCodeApi, user.SendPasswordVerificationCode)
	group.POST(consts.ChangePasswordApi, user.ChangePassword)
}

func registerIntegral(group *gin.RouterGroup) {
	group.POST(consts.RechargeApi, integral.Recharge)
	group.POST(consts.ManualRechargeApi, integral.ManualRecharge)
	group.POST(consts.GenerateRechargeKeyApi, integral.GenerateKey)
}

func registerOpenAI(group *gin.RouterGroup) {
	group.POST(consts.OpenaiCompletionsApi, chat.Completions)
	group.POST(consts.ImageApi, chat.Image)
	group.POST(consts.AudioApi, chat.Audio)
}

func registerAlipay(group *gin.RouterGroup) {
	group.POST(consts.NotifyApi, alipayapi.Notify)
}

func registerPay(group *gin.RouterGroup) {
	group.POST(consts.PreCreateApi, payapi.PreCreate)
}
