package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/user"
	"stock-web-be/logic/userapi"
	"stock-web-be/utils"
	"strconv"
)

// @Tags	用户相关接口
// @Summary	用户登录
// @Accept		json
// @Produce	json
// @param		req	body		user.LoginRequest	true	"登录请求参数"
// @Success	200	{string}	string					"返回token"
// @Router		/api/v1/user/login [post]
func Login(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req user.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	//验证邮箱格式
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

	if existUser == nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "not found user by email")
		cg.Res(http.StatusBadRequest, controller.ErrEmailNotFound)
		return
	}

	//验证密码md5值
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(req.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "compare hash password not match")
			cg.Res(http.StatusBadRequest, controller.ErrPasswordNotMatch)
			return
		}
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "compute hash password error")
		cg.Res(http.StatusBadRequest, controller.ErrComputeHashPassword)
		return
	}

	//生成jwt token
	//对userId, email加入jwt信息中
	token, err := userapi.GenerateToken(strconv.FormatUint(existUser.ID, 10), req.Email)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "generate token error")
		cg.Res(http.StatusBadRequest, controller.ErrGenerateJwtToken)
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, token)
}
