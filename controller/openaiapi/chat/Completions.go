package chat

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/openai"
)

// @Tags	OpenAI相关接口
// @Summary	对话
// @param		req	body	openai.CompletionsRequest	true	"openai请求参数"
// @Router		/api/v1/openai/v1/completions [post]
func Completions(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req openai.CompletionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
}
