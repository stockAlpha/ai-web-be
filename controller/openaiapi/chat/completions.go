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
	ctx := context.Background()
	resp, err := client.CreateChatCompletion(
		ctx,
		req,
	)

	if req.Stream {
		c.Header("Transfer-Encoding", "chunked")
		stream, err := client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			fmt.Printf("ChatCompletionStream error: %v\n", err)
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				return
			}

			if _, err := c.Writer.Write([]byte(response.Choices[0].Delta.Content)); err != nil {
				// 发送失败，退出协程
				return
			}

			// 强制刷新响应缓冲区，将数据发送给客户端
			c.Writer.(http.Flusher).Flush()
		}

	} else {
		if err != nil {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request openai error: %s", err.Error())
			cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
			return
		}
		cg.Resp(http.StatusOK, controller.ErrnoSuccess, resp.Choices[0].Message.Content)
	}
}
