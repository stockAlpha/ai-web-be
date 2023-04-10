package integral

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/integral"
)

// @Tags	积分相关接口
// @Summary	充值
// @param		req	body	integral.RechargeRequest	true	"充值请求参数"
// @Router		/api/v1/integral/recharge [post]
func Recharge(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	email := c.GetString("email")
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
	if err != nil {
		// todo：判断状态,更新次数
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKey)
		return
	}

	// 修改状态
	rechargeKey.Status = 1
	rechargeKey.UseAccount = email
	err = rechargeKey.UpdateRechargeKey()
	if err != nil {
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKeyUsed)
		return
	}
	cg.Res(http.StatusOK, controller.ErrnoSuccess)
}
