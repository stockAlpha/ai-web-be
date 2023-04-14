package user

import (
	"fmt"
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

// @Tags	用户相关接口
// @Summary	发送验证码
// @param		req	body	user.SendVerificationCodeRequest	true	"发送验证码请求参数(默认为email)"
// @Router		/api/v1/user/verify/send_code [post]
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
	existUser, err := userapi.GetUserByEmail(req.Email)
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
	subject := "ChatAlpha 验证码邮件"
	body := fmt.Sprintf("尊敬的 ChatAlpha 用户\n\n您好！感谢您使用 ChatAlpha 服务。\n\n为了确保您的账户安全，我们已向您发送了一封验证码邮件，请勿将验证码泄露给他人。验证码用于验证您的身份，并防止恶意攻击。\n\n验证码：【%s】\n\n如果您没有进行任何操作，或者不希望继续使用 ChatAlpha 服务，请忽略此邮件。\n\n如有任何问题或疑问，请随时联系我们的客服团队，我们将尽快为您解决问题。\n\n谢谢！\n\nChatAlpha 团队", code)
	err = notify.SendEmail(req.Email, subject, body)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "send verification code occur err %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrSendMailFail)
		return
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
