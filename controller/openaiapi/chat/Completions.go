package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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

	key := conf.Handler.GetString(`openai.key`)
	requestJSON, _ := json.Marshal(req)

	client := &http.Client{}
	openAIReq, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestJSON))
	openAIReq.Header.Add("Authorization", "Bearer "+key)
	openAIReq.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(openAIReq)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request OpenAI, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var response openai.CompletionsResponse
	_ = json.NewDecoder(resp.Body).Decode(&response)
	fmt.Println("response:", response)
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, response)
	return
}
