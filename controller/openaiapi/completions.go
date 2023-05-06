package openaiapi

import (
	"context"
	"errors"
	"github.com/pkoukk/tiktoken-go"
	"io"
	"math"
	"net/http"
	"stock-web-be/idl/aiapi"
	"strconv"

	"stock-web-be/controller"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/logic/userapi"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// @Tags	OpenAI相关接口
// @Summary	对话
// @param		req	body	aiapi.ChatCompletionRequest	true	"openai请求参数"
// @Router		/api/v1/openai/v1/chat/completions [post]
func Completions(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	apiKey := conf.Handler.GetString(`openai.key`)
	client := openai.NewClient(apiKey)
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	var req aiapi.ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	// 根据用户是否位vip来控制max_tokens
	user, _ := userapi.GetUserById(userId)
	// 普通用户只支持2k tokens
	maxModelTokens := 2000

	// vip用户可以支持更多的数量
	if user.VipUser {
		maxModelTokens = 4000
	}
	// 模型最大tokens
	userTokens := 0
	var messages []openai.ChatCompletionMessage
	encoding, _ := tiktoken.EncodingForModel(req.Model)
	for i := len(req.Messages) - 1; i >= 0; i-- {
		curMessage := req.Messages[i]
		// 跳过system的数据
		if curMessage.Role == "system" {
			break
		}
		content := curMessage.Content
		curLen := len(encoding.Encode(content, nil, nil))
		if userTokens+curLen < maxModelTokens/2 {
			userTokens += curLen
			messages = append([]openai.ChatCompletionMessage{req.Messages[i]}, messages...)
		} else {
			// 最后一个总会要保留
			if userTokens == 0 {
				userTokens += curLen
				messages = append([]openai.ChatCompletionMessage{req.Messages[i]}, messages...)
			}
			break
		}
	}
	if req.Role == "" {
		messages = append([]openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "你现在是一个专业的AI对话助手,你的名字叫ChatAlpha",
			},
		}, messages...)
	}
	// todo：支持按角色的prompt

	maxTokens := int(math.Max(1, math.Min(float64(maxModelTokens-userTokens), float64(maxModelTokens/2))))

	openaiReq := openai.ChatCompletionRequest{
		Model:            req.Model,
		Messages:         messages,
		MaxTokens:        maxTokens,
		Temperature:      req.Temperature,
		Stream:           req.Stream,
		FrequencyPenalty: req.FrequencyPenalty,
	}
	req.Messages = messages
	req.MaxTokens = maxTokens
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
		_ = userapi.AddUserIntegral(userId, amount, nil)
		return
	}
	if req.Stream {
		c.Header("Transfer-Encoding", "chunked")
		stream, err := client.CreateChatCompletionStream(ctx, openaiReq)
		if err != nil {
			c.Header("content-type", "application/json")
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "ChatCompletionStream error: %s", err.Error())
			cg.Res(http.StatusBadRequest, controller.ErrServer)
			// 补回积分
			_ = userapi.AddUserIntegral(userId, amount, nil)
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

			// 强制刷新响应缓冲区，将数据发送给客户端
			c.Writer.(http.Flusher).Flush()
		}
	} else {
		resp, err := client.CreateChatCompletion(
			ctx,
			openaiReq,
		)
		if err != nil {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request openai error: %s", err.Error())
			cg.Res(http.StatusBadRequest, controller.ErrServer)
			// 补回积分
			_ = userapi.AddUserIntegral(userId, amount, nil)
			return
		}
		cg.Resp(http.StatusOK, controller.ErrnoSuccess, resp.Choices[0].Message.Content)
	}
}
