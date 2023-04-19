package chat

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/dao/db"
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
	case "256×256":
		amount = amount * 2
	case "512×512":
		amount = amount * 3
	case "1024×1024":
		amount = amount * 4
	default:
		amount = amount * 2
	}
	tx := db.DbIns.Begin()
	err = userapi.SubUserIntegral(userId, amount, tx)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "record user integral error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrServer)
		tx.Rollback()
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, respUrl.Data)
	tx.Commit()
}
