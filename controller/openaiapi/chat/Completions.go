package chat

import (
	"context"
	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
)

// @Tags	代理OpenAI相关接口
// @Summary	对话
// @param		req	body	openai.ChatCompletionRequest	true	"openai请求参数"
// @Router		/api/v1/openai/v1/chat/completions [post]
func Completions(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	apiKey := conf.Handler.GetString(`openai.key`)
	client := openai.NewClient(apiKey)
	var req openai.ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	resp, err := client.CreateChatCompletion(
		context.Background(),
		req,
	)

	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request openai error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, resp)
}
