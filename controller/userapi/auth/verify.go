package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/user"
	"stock-web-be/logic/userapi"
	"stock-web-be/logic/userapi/notify"
	"stock-web-be/utils"
)

func SendVerificationCode(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req user.SendVerificationCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	//验证邮箱验证码格式
	if req.Email == "" || !utils.IsEmailValid(req.Email) {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "email is out of specification")
		cg.Res(http.StatusBadRequest, controller.ErrNotFormatEmail)
		return
	}

	//验证当前邮箱是否已注册
	existUser, err := userapi.GetUserInfoByEmail(req.Email)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query existUser by email is fatal")
		cg.Res(http.StatusBadRequest, controller.ErrnoError)
		return
	}

	if existUser != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query existUser by email is fatal")
		cg.Res(http.StatusBadRequest, controller.ErrEmailAlreadyExists)
		return
	}

	//生成随机验证码(是否需要控频?)
	code := utils.GenerateCode()

	//发送验证码
	subject := "验证码"
	body := "您的验证码为：" + code
	err = notify.SendEmail(req.Email, subject, body)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "send verification code occur err %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrSendMailFail)
	}

	//验证码存入db
	err = userapi.InsertEmailVerificationCode(code, req.Email)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "store verification code occur err")
		cg.Res(http.StatusBadRequest, controller.ErrStoreEmailCode)
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, true)
}
