package integral

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/integral"
	"stock-web-be/logic/userapi/notify"
	"time"
)

// @Tags	积分相关接口
// @Summary	生成积分充值密钥并发送到指定邮箱
// @param		req	body	integral.BatchGenerateKeyRequest	true	"生成key的请求参数"
// @Router		/api/v1/integral/generate_key [post]
func GenerateKey(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	email := c.GetString("email")
	if consts.CanGenerateRechargeKeyUserMap[email] == "" {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "user: %s not allow generate key", email)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not authorized",
		})
		return
	}
	var req integral.BatchGenerateKeyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	count := req.Count
	ret := "Generate key as follow: \n"
	for i := 0; i < count; i++ {
		key := uuid.New().String()
		rechargeKey := &db.RechargeKey{
			RechargeKey: key,
			Status:      0,
			Type:        req.Type,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		}
		// 插入生成的key
		err := rechargeKey.InsertRechargeKey()
		if err != nil {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "insert recharge key error: %s", err.Error())
			cg.Res(http.StatusBadRequest, controller.ErrRechargeKey)
			return
		}
		// 使用换行符拼接key
		ret += key + "\n"
	}
	// 通过邮件发送key
	err := notify.SendEmail("stalary@163.com", "generate recharge key", ret)
	if err != nil {
		cg.Res(http.StatusBadRequest, controller.ErrGenerateRechargeKey)
	}
	cg.Res(http.StatusOK, controller.ErrnoSuccess)
}
