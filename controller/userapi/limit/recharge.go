package limit

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/limit"
)

// @Summary	充值
// @Accept		json
// @Produce	json
// @param		req	body	limit.RechargeRequest	true	"充值请求参数"
// @Router		/api/v1/user/recharge [post]
func Recharge(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req limit.RechargeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	key := req.Key
	if key == "" {
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKey)
	}

	rechargeKey := &db.RechargeKey{}
	err := rechargeKey.GetRechargeKey(key)
	if err != nil {
		// todo：判断状态,更新次数
		cg.Res(http.StatusBadRequest, controller.ErrRechargeKey)
	}
}
