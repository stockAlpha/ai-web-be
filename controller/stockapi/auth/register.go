package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/stockapi/user"
	"stock-web-be/logic/stockapi"
	"stock-web-be/utils"
	"strconv"
)

func Register(c *gin.Context) {
	cg := controller.Gin{Ctx: c}

	var req user.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	fmt.Println(req)

	//验证邮箱格式
	if req.Email == "" || !utils.IsEmailValid(req.Email) {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "email is out of specification")
		cg.Res(http.StatusBadRequest, controller.ErrNotFormatEmail)
		return
	}

	//校验密码格式
	if !utils.IsValidPasswordFormat(req.Password) {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "password is out of specification")
		cg.Res(http.StatusBadRequest, controller.ErrNotFormatEmail)
		return
	}

	//验证当前邮箱是否已注册
	existUser, err := stockapi.GetUserInfoByEmail(req.Email)
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

	//验证code是否存在
	existCode, err := stockapi.ExistCode(req.Code, req.Email)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query code by email is fatal")
		cg.Res(http.StatusBadRequest, controller.ErrQueryVerificationCode)
		return
	}

	if !existCode {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "the code is not exist or expire")
		cg.Res(http.StatusBadRequest, controller.ErrVerificationCodeNotFound)
		return
	}

	//对密码进行加密,并添加用户
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "compute hash password err")
		cg.Res(http.StatusBadRequest, controller.ErrComputeHashPassword)
		return
	}

	userId, err := stockapi.AddUser(req.Email, hashPassword, req.TenantId)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "add user error")
		cg.Res(http.StatusBadRequest, controller.ErrAddUser)
		return
	}

	//对userId, email加入jwt信息中
	token, err := stockapi.GenerateToken(strconv.FormatUint(userId, 10), req.Email)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "generate token error")
		cg.Res(http.StatusBadRequest, controller.ErrGenerateJwtToken)
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, token)
}
