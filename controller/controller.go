package controller

import (
	"encoding/json"
	"fmt"
	"stock-web-be/gocommon/tlog"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	Ctx *gin.Context
}

type GinResp struct {
	ErrNo  int         `json:"errNo"`
	ErrMsg string      `json:"errMsg"`
	LogID  string      `json:"logid"`
	Data   interface{} `json:"data"`
}

func (g *Gin) Resp(httpCode, errCode int, data interface{}) {
	res := GinResp{
		ErrNo:  errCode,
		ErrMsg: ErrMsg[errCode],
		LogID:  g.Ctx.GetString(tlog.LOGID),
		Data:   data,
	}

	g.Ctx.JSON(httpCode, res)
}

func (g *Gin) RespWithMsg(httpCode, errCode int, msg string, data interface{}) {
	resp := GinResp{
		ErrNo:  errCode,
		ErrMsg: msg,
		LogID:  g.Ctx.GetString(tlog.LOGID),
		Data:   data,
	}

	g.Ctx.JSON(httpCode, resp)
}

func (g *Gin) Res(httpCode, errCode int) {
	res := GinResp{
		ErrNo:  errCode,
		ErrMsg: ErrMsg[errCode],
		LogID:  g.Ctx.GetString(tlog.LOGID),
		Data:   nil,
	}

	g.Ctx.JSON(httpCode, res)
}

// 业务错误码(对外)，非http状态码
const (
	ErrnoSuccess = 0
	ErrnoError   = 1

	ErrnoInvalidPrm = 40000

	ErrNotFormatEmail           = 40100
	ErrEmailAlreadyExists       = 40101
	ErrSendMailFail             = 40102
	ErrStoreEmailCode           = 40103
	ErrQueryVerificationCode    = 40104
	ErrVerificationCodeNotFound = 40105
	ErrNotFormatPassword        = 40106
	ErrGenerateJwtToken         = 40107
)

var ErrMsg = map[int]string{
	ErrnoSuccess: "success",
	ErrnoError:   "error",

	ErrnoInvalidPrm:             "invalid parms",
	ErrNotFormatEmail:           "not format email",
	ErrEmailAlreadyExists:       "email already exists",
	ErrSendMailFail:             "send verification code error",
	ErrStoreEmailCode:           "store email verification code error",
	ErrQueryVerificationCode:    "query verification code error",
	ErrVerificationCodeNotFound: "verification code not found",
	ErrNotFormatPassword:        "not format password",
	ErrGenerateJwtToken:         "generate jwt token error",
}

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
	b, err := json.Marshal(GinResp{0, "", "", data})
	if err != nil {
		ctx.Data(200, "application/json", []byte(`{"errno":1, "errmsg":"`+err.Error()+`","data":null}`))
	} else {
		ctx.Data(200, "application/json", b)
	}
}
