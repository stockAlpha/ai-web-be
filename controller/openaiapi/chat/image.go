package chat

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"stock-web-be/async"
	"stock-web-be/controller"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/openai"
	"stock-web-be/logic/userapi"

	"github.com/gin-gonic/gin"
	public_openai "github.com/sashabaranov/go-openai"
)

// @Tags	代理OpenAI相关接口
// @Summary	生成图片
// @param		req	body	openai.ImageRequest	true	"openai请求参数"
// @Router		/api/v1/openai/v1/image [post]
func Image(c *gin.Context) {
	prostartTime := time.Now()
	cg := controller.Gin{Ctx: c}
	ctx := context.Background()
	apiKey := conf.Handler.GetString(`openai.key`)
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	client := public_openai.NewClient(apiKey)
	var req openai.ImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	uuID := req.UUID
	messageID := req.MessageID
	if uuID == 0 {
		uuID, _ = strconv.Atoi(c.GetString("uuid"))
	}

	if messageID == "" {
		messageID = c.GetString("message_id")
	}
	respUrl, err := client.CreateImage(ctx, req.ImageRequest)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	// 计费
	amount := req.N
	switch req.Size {
	case "256x256":
		amount = amount * 2
	case "512x512":
		amount = amount * 3
	case "1024x1024":
		amount = amount * 4
	default:
		amount = amount * 2
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
	// prompt冗余存储，到时候查记录不用额外处理
	chatRecordPrompt := db.ChatRecord{UserID: userId, Data: req.Prompt, UUID: uuID, MessageID: messageID, DataType: 0, CreatedAt: prostartTime}
	chatRecord := async.ChatRecordChanType{Record: []db.ChatRecord{chatRecordPrompt}}
	for i := range respUrl.Data {
		chatRecord.Record = append(chatRecord.Record, db.ChatRecord{UserID: userId, Prompt: req.Prompt, Data: respUrl.Data[i].URL, UUID: uuID, MessageID: messageID, DataType: 2, CreatedAt: time.Now()})
	}
	async.ChatRecordChan <- chatRecord
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, respUrl.Data)
}
