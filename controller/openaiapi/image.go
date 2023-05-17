package openaiapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	"stock-web-be/idl/aiapi"
	"stock-web-be/logic/aliyunapi"
	"stock-web-be/logic/userapi"
	"stock-web-be/logic/xfapi"
	"stock-web-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// @Tags	OpenAI相关接口
// @Summary	生成图片
// @param		req	body	aiapi.ImageRequest	true	"openai请求参数"
// @Router		/api/v1/openai/v1/image [post]
func Image(c *gin.Context) {
	prostartTime := time.Now()
	cg := controller.Gin{Ctx: c}
	ctx := context.Background()
	apiKey := conf.Handler.GetString(`openai.key`)
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	client := openai.NewClient(apiKey)
	var req aiapi.ImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	uuID := req.UUID
	messageID := req.MessageID
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
	err := userapi.SubUserIntegral(userId, amount, nil)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "record user integral error: %s", err.Error())
		if err.Error() == "余额不足" {
			cg.Res(http.StatusBadRequest, controller.ErrIntegralNotEnough)
		} else {
			cg.Res(http.StatusBadRequest, controller.ErrServer)
		}
	}
	var res []aiapi.ImageResponseDataInner

	prompt := req.Prompt
	if err != nil {
		return
	}
	if utils.ContainsChinese(prompt) {
		prompt = xfapi.Run(prompt, "cn", "en")
	}
	// prompt冗余存储，到时候查记录不用额外处理
	//todo 还没加sd的
	chatRecordPrompt := db.ChatRecord{UserID: userId, Data: req.Prompt, UUID: uuID, MessageID: messageID, DataType: 0, CreatedAt: prostartTime}
	chatRecord := async.ChatRecordChanType{Record: []db.ChatRecord{chatRecordPrompt}}
	if req.Model == "stable-diffusion" {
		token := "Token " + conf.Handler.GetString("replicate.key")

		reqValue := aiapi.ReplicateStableDiffusion{
			Version: conf.Handler.GetString("replicate.stable_diffusion_version"),
			Input: aiapi.ReplicateInput{
				Prompt:          prompt,
				NumOutputs:      req.N,
				ImageDimensions: req.Size,
			},
		}
		j, _ := json.Marshal(reqValue)
		postReq, _ := http.NewRequest(http.MethodPost, "https://api.replicate.com/v1/predictions", bytes.NewReader(j))
		postReq.Header.Add("Authorization", token)
		postReq.Header.Add("Content-Type", "application/json")
		client := http.Client{}
		postRes, err := client.Do(postReq)
		if err != nil {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "replicate predictions error: %s", err.Error())
			// 补回积分
			_ = userapi.AddUserIntegral(userId, amount, nil)
			return
		}
		defer postRes.Body.Close()
		body, _ := io.ReadAll(postRes.Body)
		var replicateRes aiapi.ReplicateResponse
		_ = json.Unmarshal(body, &replicateRes)
		res, err = replicateGet(replicateRes.Urls.Get, token)
		if err != nil {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "replicate get error: %s", err.Error())
			// 补回积分
			_ = userapi.AddUserIntegral(userId, amount, nil)
			return
		}
	} else {
		respUrl, err := client.CreateImage(ctx, openai.ImageRequest{
			Prompt: prompt,
			N:      req.N,
			Size:   req.Size,
		})
		if err != nil {
			fmt.Printf("Image creation error: %v\n", err)
			cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
			return
		}

		for i := range respUrl.Data {
			res = append(res, aiapi.ImageResponseDataInner{URL: aliyunapi.UploadFileByUrl(respUrl.Data[i].URL, "image/jpg")})
			chatRecord.Record = append(chatRecord.Record, db.ChatRecord{UserID: userId, Prompt: req.Prompt, Data: respUrl.Data[i].URL, UUID: uuID, MessageID: messageID, DataType: 2, CreatedAt: time.Now()})
		}
	}
	async.ChatRecordChan <- chatRecord
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, res)

}

var httpClient = http.Client{}

func replicateGet(url, auth string) (output []aiapi.ImageResponseDataInner, err error) {
	// create a HTTP GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json")

	// create a HTTP client and use it to send the request
	var res []aiapi.ImageResponseDataInner

	for {
		resp, err := httpClient.Do(req)
		if err != nil {
			break
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			break
		}
		response := aiapi.ReplicateResponse{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			break
		}
		if response.Status == "succeeded" {
			for i := range response.Output {
				res = append(res, aiapi.ImageResponseDataInner{URL: aliyunapi.UploadFileByUrl(response.Output[i], "image/jpg")})
			}
			return res, nil
		}
		if response.Status == "failed" {
			return res, fmt.Errorf("failed")
		}
	}
	return res, err
}
