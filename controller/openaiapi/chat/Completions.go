package chat

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/openai"
)

// @Tags	代理OpenAI相关接口
// @Summary	对话
// @param		req	body	openai.CompletionsRequest	true	"openai请求参数"
// @Router		/api/v1/openai/v1/chat/completions [post]
func Completions(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req openai.CompletionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	apiKey := req.OpenAIKey
	if apiKey == "" {
		apiKey = conf.Handler.GetString(`openai.key`)
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		fmt.Println("Error sending request:", err)
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	if resp.StatusCode() != 200 {
		fmt.Println("res error", resp)
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, resp)
	return
}
