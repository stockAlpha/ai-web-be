package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/logic/userapi"
	"stock-web-be/utils"
	"strconv"
)

// @Tags	代理OpenAI相关接口
// @Summary	对话
// @param		req	body	openai.ChatCompletionRequest	true	"openai请求参数"
// @Router		/api/v1/openai/v1/chat/completions [post]
func Completions(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	apiKey := conf.Handler.GetString(`openai.key`)
	client := openai.NewClient(apiKey)
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	var req openai.ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	ctx := context.Background()
	tx := db.DbIns.Begin()

	// 计费，对话次数目前都按照1来计费
	err := userapi.SubUserIntegral(userId, 1, tx)
	if err != nil {
		c.Header("content-type", "application/json")
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "record user integral error: %s", err.Error())
		if err.Error() == "余额不足" {
			cg.Res(http.StatusBadRequest, controller.ErrIntegralNotEnough)
		} else {
			cg.Res(http.StatusBadRequest, controller.ErrServer)
		}
		tx.Rollback()
		return
	}
	if req.Stream {
		c.Header("Transfer-Encoding", "chunked")
		stream, err := client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			c.Header("content-type", "application/json")
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "ChatCompletionStream error: %s", err.Error())
			cg.Res(http.StatusBadRequest, controller.ErrServer)
			tx.Rollback()
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				c.Header("content-type", "application/json")
				tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "Stream error: %s", err.Error())
				cg.Res(http.StatusBadRequest, controller.ErrServer)
				tx.Rollback()
				return
			}

			content := response.Choices[0].Delta.Content

			fmt.Printf("openai return content:%s, response:%s", content, response)
			content = utils.ReplaceSensitiveWord(content, consts.SensitiveWordReplaceMap)
			fmt.Printf("openai replace return content:%s", content)

			if _, err := c.Writer.Write([]byte(content)); err != nil {
				// 发送失败，退出协程
				return
			}
			// 强制刷新响应缓冲区，将数据发送给客户端
			c.Writer.(http.Flusher).Flush()
		}
	} else {
		resp, err := client.CreateChatCompletion(
			ctx,
			req,
		)
		if err != nil {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request openai error: %s", err.Error())
			cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
			tx.Rollback()
			return
		}
		cg.Resp(http.StatusOK, controller.ErrnoSuccess, resp.Choices[0].Message.Content)
	}
	tx.Commit()
}
