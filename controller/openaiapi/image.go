package openaiapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/aiapi"
	"stock-web-be/logic/aliyunapi"
	"stock-web-be/logic/userapi"
	"stock-web-be/logic/xfapi"
	"stock-web-be/utils"
	"strconv"
	"time"
)

// @Tags	AI相关接口
// @Summary	生成图片
// @param		req	body	aiapi.ImageRequest	true	"生成图片请求参数"
// @Router		/api/v1/openai/v1/image [post]
func Image(c *gin.Context) {
	cg := controller.Gin{Ctx: c}

	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	var req aiapi.ImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	// 计费
	//amount := req.N
	//switch req.Size {
	//case "256x256":
	//	amount = amount * 7
	//case "512x512":
	//	amount = amount * 8
	//case "1024x1024":
	//	amount = amount * 9
	//default:
	//	// 默认为8积分
	//	amount = amount * 8
	//}
	// 目前统一为8积分
	amount := 8
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

	if req.Model == "stable-diffusion" {
		token := "Token " + conf.Handler.GetString("replicate.key")
		replicateRes, err := replicateRequest(token, prompt, req.Size, req.N, c)
		if err != nil {
			_ = userapi.AddUserIntegral(userId, amount, nil)
			cg.Res(http.StatusBadRequest, controller.ErrServer)
			return
		}
		res, err = replicateGet(replicateRes.Urls.Get, token, c)
		if err != nil {
			_ = userapi.AddUserIntegral(userId, amount, nil)
			cg.Res(http.StatusBadRequest, controller.ErrServer)
			return
		}
	} else if req.Model == "dall-e2" {
		res, err = dalle2Get(prompt, req.Size, req.N, c)
		if err != nil {
			_ = userapi.AddUserIntegral(userId, amount, nil)
			cg.Res(http.StatusBadRequest, controller.ErrServer)
			return
		}
	} else if req.Model == "midjourney" {
		res, err = midjourneyRequest(prompt, c)
		if err != nil {
			_ = userapi.AddUserIntegral(userId, amount, nil)
			cg.Res(http.StatusBadRequest, controller.ErrServer)
			return
		}
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, res)
}

// @Tags	AI相关接口
// @Summary	操作图片
// @param		req	body	aiapi.MjProxyOperate	true	"图片操作请求参数"
// @Router		/api/v1/openai/v1/image/operate [post]
func ImageOperate(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req aiapi.MjProxyOperate
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	res, err := midjourneyOperate(req, c)
	if err != nil {
		cg.Res(http.StatusBadRequest, controller.ErrServer)
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, res)
}

var httpClient = http.Client{}

func midjourneyGet(taskId string, c *gin.Context) ([]aiapi.ImageResponseDataInner, error) {
	var res []aiapi.ImageResponseDataInner
	host := conf.Handler.GetString("midjourney.host")
	req, err := http.NewRequest(http.MethodGet, host+"/mj/task/"+taskId+"/fetch", nil)
	if err != nil {
		return res, err
	}
	req.Header.Add("Content-Type", "x-www-form-urlencoded")
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
		response := aiapi.MjProxyGetRes{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			break
		}
		if response.Status == "SUCCESS" {
			res = append(res, aiapi.ImageResponseDataInner{URL: response.ImageUrl, TaskId: taskId})
			return res, nil
		}
		fmt.Println("image res", response)
		if response.Status == "FAILURE" {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "midjourney get error: %s", err.Error())
			return res, fmt.Errorf(response.FailReason)
		}
		time.Sleep(time.Second * 1)
	}
	return res, fmt.Errorf("midjourney get error")
}

func midjourneyOperate(operate aiapi.MjProxyOperate, c *gin.Context) ([]aiapi.ImageResponseDataInner, error) {
	reqValue := aiapi.MjProxySubmit{
		Action: operate.Action,
		Index:  operate.Index,
		TaskId: operate.TaskId,
	}
	j, _ := json.Marshal(reqValue)
	host := conf.Handler.GetString("midjourney.host")
	postReq, _ := http.NewRequest(http.MethodPost, host+"/mj/submit/change", bytes.NewReader(j))
	postReq.Header.Add("Content-Type", "application/json")
	postRes, err := httpClient.Do(postReq)
	var mjProxyRes aiapi.MjProxySubmitRes
	var res []aiapi.ImageResponseDataInner
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "midjourney operate error: %s", err.Error())
		return res, err
	}
	defer postRes.Body.Close()
	body, _ := io.ReadAll(postRes.Body)
	_ = json.Unmarshal(body, &mjProxyRes)
	return midjourneyGet(mjProxyRes.Result, c)
}

func midjourneyRequest(prompt string, c *gin.Context) ([]aiapi.ImageResponseDataInner, error) {
	reqValue := aiapi.MjProxySubmit{
		Prompt: prompt,
	}
	j, _ := json.Marshal(reqValue)
	// mj服务只部署一个，所以先用外网调用，
	host := conf.Handler.GetString("midjourney.host")
	postReq, _ := http.NewRequest(http.MethodPost, host+"/mj/submit/imagine", bytes.NewReader(j))
	postReq.Header.Add("Content-Type", "application/json")
	postRes, err := httpClient.Do(postReq)
	var mjProxyRes aiapi.MjProxySubmitRes
	var res []aiapi.ImageResponseDataInner
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "midjourney request error: %s", err.Error())
		return res, err
	}
	defer postRes.Body.Close()
	body, _ := io.ReadAll(postRes.Body)
	_ = json.Unmarshal(body, &mjProxyRes)
	return midjourneyGet(mjProxyRes.Result, c)
}

func dalle2Get(prompt, size string, n int, c *gin.Context) ([]aiapi.ImageResponseDataInner, error) {
	ctx := context.Background()
	apiKey := conf.Handler.GetString(`openai.key`)
	client := openai.NewClient(apiKey)
	var res []aiapi.ImageResponseDataInner
	respUrl, err := client.CreateImage(ctx, openai.ImageRequest{
		Prompt: prompt,
		N:      n,
		Size:   size,
	})
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "dalle2 get error: %s", err.Error())
		return res, err
	}

	for i := range respUrl.Data {
		res = append(res, aiapi.ImageResponseDataInner{URL: aliyunapi.UploadFileByUrl(respUrl.Data[i].URL, "image/jpg")})
	}
	return res, nil
}

func replicateRequest(token, prompt, size string, n int, c *gin.Context) (aiapi.ReplicateResponse, error) {
	reqValue := aiapi.ReplicateStableDiffusion{
		Version: conf.Handler.GetString("replicate.stable_diffusion_version"),
		Input: aiapi.ReplicateInput{
			Prompt:          prompt,
			NumOutputs:      n,
			ImageDimensions: size,
		},
	}
	j, _ := json.Marshal(reqValue)
	postReq, _ := http.NewRequest(http.MethodPost, "https://api.replicate.com/v1/predictions", bytes.NewReader(j))
	postReq.Header.Add("Authorization", token)
	postReq.Header.Add("Content-Type", "application/json")
	postRes, err := httpClient.Do(postReq)
	var replicateRes aiapi.ReplicateResponse
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "replicate predictions error: %s", err.Error())
		return replicateRes, err
	}
	defer postRes.Body.Close()
	body, _ := io.ReadAll(postRes.Body)
	_ = json.Unmarshal(body, &replicateRes)
	return replicateRes, err
}

func replicateGet(url, auth string, c *gin.Context) ([]aiapi.ImageResponseDataInner, error) {
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
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "replicate get error: %s", err.Error())
			return res, fmt.Errorf("failed")
		}
	}
	return res, err
}
