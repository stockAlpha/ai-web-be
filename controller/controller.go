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
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	LogID string      `json:"logid"`
	Data  interface{} `json:"data"`
}

func (g *Gin) Resp(httpCode, errCode int, data interface{}) {
	res := GinResp{
		Code:  errCode,
		Msg:   ErrMsg[errCode],
		LogID: g.Ctx.GetString(tlog.LOGID),
		Data:  data,
	}

	g.Ctx.JSON(httpCode, res)
}

func (g *Gin) RespWithMsg(httpCode, errCode int, msg string, data interface{}) {
	resp := GinResp{
		Code:  errCode,
		Msg:   msg,
		LogID: g.Ctx.GetString(tlog.LOGID),
		Data:  data,
	}

	g.Ctx.JSON(httpCode, resp)
}

func (g *Gin) Res(httpCode, errCode int) {
	res := GinResp{
		Code:  errCode,
		Msg:   ErrMsg[errCode],
		LogID: g.Ctx.GetString(tlog.LOGID),
		Data:  nil,
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
	ErrComputeHashPassword      = 40108
	ErrPasswordNotMatch         = 40109
	ErrEmailNotFound            = 40110
	ErrAddUser                  = 40111
	ErrRechargeKey              = 40112
	ErrGenerateRechargeKey      = 40113
	ErrRechargeKeyUsed          = 40114
	ErrAddIntegral              = 40115
	ErrRegister                 = 40116
	ErrIntegralNotEnough        = 40117
	ErrServer                   = 40118
	ErrUserSubjectIdNotFound    = 40119
	ErrPasswordNotChange        = 40120
	ErrChangePassword           = 40121
	ErrNotExistToken            = 40122
	ErrTokenNotExistTime        = 40123
	ErrTokenAddBlackList        = 40124
)

var ErrMsg = map[int]string{
	ErrnoSuccess: "success",
	ErrnoError:   "error",

	ErrnoInvalidPrm:             "非法的参数",
	ErrNotFormatEmail:           "邮件格式非法",
	ErrEmailAlreadyExists:       "此邮箱已注册",
	ErrSendMailFail:             "发送验证码失败",
	ErrStoreEmailCode:           "存储验证码失败",
	ErrQueryVerificationCode:    "查询验证码失败",
	ErrVerificationCodeNotFound: "验证码未找到",
	ErrNotFormatPassword:        "密码格式非法",
	ErrGenerateJwtToken:         "生成 jwt token 错误",
	ErrComputeHashPassword:      "计算密码 hash 错误",
	ErrPasswordNotMatch:         "密码输入错误",
	ErrEmailNotFound:            "邮箱未找到",
	ErrAddUser:                  "创建用户错误",
	ErrRechargeKey:              "充值错误",
	ErrGenerateRechargeKey:      "生成充值密钥错误",
	ErrRechargeKeyUsed:          "充值密钥已使用",
	ErrAddIntegral:              "添加积分错误",
	ErrRegister:                 "注册失败",
	ErrIntegralNotEnough:        "积分不足，请点击左下角设置进行充值",
	ErrServer:                   "服务器开小车了～，请重试一次",
	ErrUserSubjectIdNotFound:    "用户不存在",
	ErrPasswordNotChange:        "旧密码与新密码相同",
	ErrChangePassword:           "修改用户密码失败",
	ErrNotExistToken:            "token不存在",
	ErrTokenNotExistTime:        "token解析后不存在失效时间",
	ErrTokenAddBlackList:        "token加入黑名单失败",
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
