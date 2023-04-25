package integral

import (
	"fmt"
	"net/http"
	"stock-web-be/utils"
	"time"

	"stock-web-be/async"
	"stock-web-be/controller"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/integral"
	"stock-web-be/logic/userapi"

	"github.com/gin-gonic/gin"
)

// @Tags	积分相关接口
// @Summary	手动充值(管理员使用)
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

	amount := utils.GetAmount(rechargeKey.Type)
	// 添加积分
	tx := db.DbIns.Begin()
	err = userapi.AddUserIntegral(userId, amount, tx)
	if err != nil {
		tx.Rollback()
		cg.Res(http.StatusBadRequest, controller.ErrAddIntegral)
		return
	}

	// 修改状态
	rechargeKey.Status = 1
	rechargeKey.UseAccount = userId
	rechargeKey.UpdateTime = time.Now()
	err = rechargeKey.UpdateRechargeKey(tx)
	if err != nil {
		tx.Rollback()
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "use recharge key error, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKeyUsed)
		return
	}
	tx.Commit()
	async.MailChan <- async.MailChanType{To: email, Subject: consts.RechargeNotifySubject, Body: fmt.Sprintf(consts.RechargeNotifyContent, amount)}
	cg.Res(http.StatusOK, controller.ErrnoSuccess)
}
