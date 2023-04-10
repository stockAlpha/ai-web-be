package chat

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
)

// @Tags	代理OpenAI相关接口
// @Summary	对话
// @param		req	body	openai.ImageRequest	true	"openai请求参数"
// @Router		/api/v1/openai/v1/image [post]
func Image(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	ctx := context.Background()
	apiKey := conf.Handler.GetString(`openai.key`)
	client := openai.NewClient(apiKey)
	var req openai.ImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	respUrl, err := client.CreateImage(ctx, req)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, respUrl.Data[0].URL)
}
