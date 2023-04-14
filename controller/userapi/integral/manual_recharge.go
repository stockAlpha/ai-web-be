package integral

import (
	"github.com/gin-gonic/gin"
	"github.com/stockAlpha/gopkg/common/safego"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/integral"
	"stock-web-be/logic/userapi"
	"stock-web-be/logic/userapi/notify"
	"strconv"
)

// @Tags	积分相关接口
// @Summary	手动充值
// @param		req	body	integral.ManualRechargeRequest	true	"手动充值请求参数"
// @Router		/api/v1/integral/manual/recharge [post]
func ManualRecharge(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req integral.ManualRechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	if req.AuthCode != "chat_alpha_0414" {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request authCode invalid, %s", req.AuthCode)
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	user, err := userapi.GetUserByEmail(req.ToEmail)
	if err != nil {
		cg.Res(http.StatusBadRequest, controller.ErrEmailNotFound)
		return
	}
	userId := user.ID
	email := user.Email

	if consts.CanGenerateRechargeKeyUserMap[email] == "" {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "user: %s not allow generate key", email)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Not authorized",
		})
		return
	}

	key := req.Key
	if key == "" {
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKey)
	}

	// 判断key是否有效
	rechargeKey := &db.RechargeKey{}
	err = rechargeKey.GetRechargeKey(key)
	if err != nil || rechargeKey.ID == 0 {
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKey)
		return
	}
	if rechargeKey.Status != 0 {
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKeyUsed)
		return
	}

	amount := 0
	switch rechargeKey.Type {
	case 1:
		amount = 30
	case 2:
		amount = 100
	case 3:
		amount = 500
	default:
		amount = 30
	}
	// 添加积分
	// todo 事务
	userapi.AddUserIntegral(userId, amount)

	// 修改状态
	rechargeKey.Status = 1
	rechargeKey.UseAccount = userId
	err = rechargeKey.UpdateRechargeKey()
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "use recharge key error, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKeyUsed)
		return
	}
	safego.SafeGoWithWG(func() {
		notify.SendEmail(email, "充值成功", "您已成功充值"+strconv.Itoa(amount)+"积分，快去看看吧！")
	})
	cg.Res(http.StatusOK, controller.ErrnoSuccess)
}
