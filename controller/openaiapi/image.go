package openaiapi

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
	"stock-web-be/logic/userapi"
	"strconv"
)

// @Tags	代理OpenAI相关接口
// @Summary	生成图片
// @param		req	body	openai.ImageRequest	true	"openai请求参数"
// @Router		/api/v1/openai/v1/image [post]
func Image(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	ctx := context.Background()
	apiKey := conf.Handler.GetString(`openai.key`)
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
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
	// 计费
	amount := req.N
	switch req.Size {
	case "256x256":
		amount = amount * 7
	case "512x512":
		amount = amount * 8
	case "1024x1024":
		amount = amount * 9
	default:
		amount = amount * 8
	}
	// 先扣减积分，后面失败了再补回来
	err = userapi.SubUserIntegral(userId, amount, nil)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "record user integral error: %s", err.Error())
		if err.Error() == "余额不足" {
			cg.Res(http.StatusBadRequest, controller.ErrIntegralNotEnough)
		} else {
			cg.Res(http.StatusBadRequest, controller.ErrServer)
		}
		// 补回积分
		_ = userapi.AddUserIntegral(userId, amount, nil)
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, respUrl.Data)
}
