package chat

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"stock-web-be/controller"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/utils"
	"strconv"
)

// @Tags	代理OpenAI相关接口
// @Summary	音频转文字
// @Param audio formData file true "音频文件"
// @Param model query string false "model"
// @Param language query string false "language"
// @Param prompt query string false "prompt"
// @Param temperature query float32 false "temperature"
// @Router		/api/v1/openai/v1/audio [post]
func Audio(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	ctx := context.Background()
	apiKey := conf.Handler.GetString(`openai.key`)
	client := openai.NewClient(apiKey)
	model := c.DefaultQuery("model", openai.Whisper1)
	language := c.DefaultQuery("language", "zh")
	prompt, _ := c.GetQuery("prompt")
	temperature, _ := strconv.ParseFloat(c.DefaultQuery("temperature", "0"), 32)
	// temperature转化为float32
	fmt.Println(language)
	file, _, err := c.Request.FormFile("audio")
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "from file audio, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	filename := "recording.mp3"
	err = saveFile(file, filename)
	req := openai.AudioRequest{
		Model:       model,
		FilePath:    filename,
		Language:    language,
		Prompt:      prompt,
		Temperature: float32(temperature),
	}
	defer file.Close()
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "save audio file, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	resp, err := client.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	resp.Text = utils.ReplaceSensitiveWord(resp.Text, consts.SensitiveWordReplaceMap)
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, resp.Text)
}

func saveFile(file multipart.File, filename string) error {
	defer file.Close()

	// 创建要保存的文件
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// 将上传的文件复制到要保存的文件
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}
