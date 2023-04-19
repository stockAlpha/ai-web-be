package integral

import (
	"net/http"
	"stock-web-be/dao/db"
	"strconv"

	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/integral"
	"stock-web-be/logic/userapi"

	"github.com/gin-gonic/gin"
)

// @Tags	积分相关接口
// @Summary	记录
// @param		req	body	integral.RecordRequest	true	"请求参数"
// @Router		/api/v1/integral/record [post]
func Record(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	var req integral.RecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	amount := 0
	switch req.Type {
	case "chat":
		amount = 1
	case "image":
		amount = 4
	case "audio":
		amount = 8
	default:
		amount = 1
	}
	tx := db.DbIns.Begin()
	err := userapi.SubUserIntegral(userId, amount, tx)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "record user integral error: %s", err.Error())
		if err.Error() == "余额不足" {
			cg.Res(http.StatusBadRequest, controller.ErrIntegralNotEnough)
		} else {
			cg.Res(http.StatusBadRequest, controller.ErrServer)
		}
		tx.Rollback()
		return
	}
	cg.Res(http.StatusOK, controller.ErrnoSuccess)
	tx.Commit()
}
