package integral

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/integral"
	"stock-web-be/logic/userapi"
)

// @Tags	积分相关接口
// @Summary	记录
// @param		req	body	integral.RecordRequest	true	"请求参数"
// @Router		/api/v1/integral/record [post]
func Record(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	userId := c.GetUint64("user_id")
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
	}
	userapi.SubUserIntegral(userId, amount)
	cg.Res(http.StatusOK, controller.ErrnoSuccess)
}
