package chat

import (
	"context"
	"errors"
	"io"
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
// @Summary	对话
// @param		req	body	openai.ChatCompletionRequest	true	"openai请求参数"
// @Router		/api/v1/openai/v1/chat/completions [post]
func Completions(c *gin.Context) {
	prostartTime := time.Now()
	cg := controller.Gin{Ctx: c}
	apiKey := conf.Handler.GetString(`openai.key`)
	client := public_openai.NewClient(apiKey)
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)

	var req openai.ChatCompletionRequest
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
	ctx := context.Background()
	// 计费，避免长事务，先扣减积分，再对话
	amount := 1
	err := userapi.SubUserIntegral(userId, amount, nil)
	if err != nil {
		c.Header("content-type", "application/json")
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "record user integral error: %s", err.Error())
		if err.Error() == "余额不足" {
			cg.Res(http.StatusBadRequest, controller.ErrIntegralNotEnough)
		} else {
			cg.Res(http.StatusBadRequest, controller.ErrServer)
		}
		// 补回积分
		userapi.AddUserIntegral(userId, amount, nil)
		return
	}
	recordContent := ""
	if req.Stream {
		c.Header("Transfer-Encoding", "chunked")
		stream, err := client.CreateChatCompletionStream(ctx, req.ChatCompletionRequest)
		if err != nil {
			c.Header("content-type", "application/json")
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "ChatCompletionStream error: %s", err.Error())
			cg.Res(http.StatusBadRequest, controller.ErrServer)
			// 补回积分
			_ = userapi.AddUserIntegral(userId, amount, nil)
			return
		}
		defer stream.Close()
		defer func() {
			if recordContent != "" && len(req.Messages) > 0 {
				prompt := ""
				for i := range req.Messages {
					if req.Messages[i].Role == "user" {
						prompt = req.Messages[i].Content
						break
					}
				}
				if prompt != "" {
					// prompt冗余存储，到时候查记录不用额外处理
					chatRecordPrompt := db.ChatRecord{UserID: userId, Data: prompt, UUID: uuID, MessageID: messageID, DataType: 0, CreatedAt: prostartTime}
					chatRecordContent := db.ChatRecord{UserID: userId, Prompt: prompt, Data: recordContent, UUID: uuID, MessageID: messageID, DataType: 1, CreatedAt: time.Now()}
					async.ChatRecordChan <- async.ChatRecordChanType{Record: []db.ChatRecord{chatRecordPrompt, chatRecordContent}}
				}
			}
		}()
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				c.Header("content-type", "application/json")
				tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "Stream error: %s", err.Error())
				cg.Res(http.StatusBadRequest, controller.ErrServer)
				// 补回积分
				_ = userapi.AddUserIntegral(userId, amount, nil)
				return
			}

			if _, err := c.Writer.Write([]byte(response.Choices[0].Delta.Content)); err != nil {
				c.Header("content-type", "application/json")
				tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "Write content error: %s", err.Error())
				cg.Res(http.StatusBadRequest, controller.ErrServer)
				// 补回积分
				_ = userapi.AddUserIntegral(userId, amount, nil)
				return
			}
			recordContent = recordContent + response.Choices[0].Delta.Content
			// 强制刷新响应缓冲区，将数据发送给客户端
			c.Writer.(http.Flusher).Flush()
		}
	} else {
		resp, err := client.CreateChatCompletion(
			ctx,
			req.ChatCompletionRequest,
		)
		if err != nil {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request openai error: %s", err.Error())
			cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
			// 补回积分
			_ = userapi.AddUserIntegral(userId, amount, nil)
			return
		}
		cg.Resp(http.StatusOK, controller.ErrnoSuccess, resp.Choices[0].Message.Content)
	}
}
