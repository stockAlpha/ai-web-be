package integral

import (
	"net/http"
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
	// todo:暂时先简单计
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
	err := userapi.SubUserIntegral(userId, amount)
	if err != nil {
		// 链接错误的为了不影响用户使用，先返回成功
		if err.Error() == "invalid connection" {
			cg.Res(http.StatusOK, controller.ErrnoSuccess)
			return
		}
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "record user integral error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrIntegralNotEnough)
		return
	}
	cg.Res(http.StatusOK, controller.ErrnoSuccess)
}
