package integral

import (
	"fmt"
	"net/http"
	"strconv"

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
// @Summary	充值
// @param		req	body	integral.RechargeRequest	true	"充值请求参数"
// @Router		/api/v1/integral/recharge [post]
func Recharge(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req integral.RechargeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	key := req.Key
	if key == "" {
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKey)
	}

	// 判断key是否有效
	rechargeKey := &db.RechargeKey{}
	err := rechargeKey.GetRechargeKey(key)
	if err != nil || rechargeKey.ID == 0 {
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKey)
		return
	}
	if rechargeKey.Status != 0 {
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKeyUsed)
		return
	}

	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	email := c.GetString("email")
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
