package controller

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	ErrNo  int         `json:"errNo"`
	ErrMsg string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

// 业务错误码(对外)，非http状态码
const (
	ErrnoSuccess = 0
	ErrnoError   = 1

	ErrnoNotFoundPipeline = 20
	ErrnoMultiPipeline    = 21
)

// 业务错误码标识的信息(对外)
const (
	ErrmsgSuccess = "success"
	ErrmsgError   = "error"
)

// EchoJSON json格式输出
func EchoJSON(ctx *gin.Context, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		ctx.Data(200, "application/json", []byte(`{"errno":1, "errmsg":"`+err.Error()+`","data":null}`))
	} else {
		ctx.Data(200, "application/json", b)
	}
}

// EchoError json格式输出，只有非0错误码，和错误信息
func EchoError(ctx *gin.Context, errno int, errmsg string) {
	ctx.Data(200, "application/json", []byte(`{"errno":`+strconv.Itoa(errno)+`, "errmsg":"`+strings.Trim(strconv.Quote(errmsg), "\"")+`","data":null}`))
}

// EchoTextPlain
func EchoTextPlain(ctx *gin.Context, text string) {
	ctx.Data(200, "text/plain;charset=utf-8", []byte(text))
}

// Download Text Plain
func DownloadTextPlain(ctx *gin.Context, text string, filename string) {
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Data(200, "text/plain;charset=utf-8", []byte(text))
}

func EchoData(ctx *gin.Context, data interface{}) {
	b, err := json.Marshal(Resp{ErrnoSuccess, ErrmsgSuccess, data})
	if err != nil {
		ctx.Data(200, "application/json", []byte(`{"errno":1, "errmsg":"`+err.Error()+`","data":null}`))
	} else {
		ctx.Data(200, "application/json", b)
	}
}
