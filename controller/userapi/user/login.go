package user

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"math"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/user"
	"stock-web-be/logic/userapi"
	"stock-web-be/utils"
	"strconv"
	"time"
)

// @Tags	用户相关接口
// @Summary	登录
// @param		req	body		user.LoginRequest	true	"登录请求参数"
// @Success	200	{string}	string					"返回token"
// @Router		/api/v1/user/login [post]
func Login(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req user.LoginRequest
	time.Sleep(10 * time.Second)
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
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query existUser by email error", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrServer)
		return
	}

	if existUser == nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "not found user by email")
		cg.Res(http.StatusBadRequest, controller.ErrEmailNotFound)
		return
	}

	// 验证密码md5值
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

// @Tags	用户相关接口
// @Summary	登出
// @Success	200	{string}	string					"返回token"
// @Router		/api/v1/user/logout [post]
func Logout(c *gin.Context) {
	cg := controller.Gin{Ctx: c}

	//1.获取当前token
	token, exist := c.Get("token")
	if !exist {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "token is not exist")
		cg.Res(http.StatusBadRequest, controller.ErrNotExistToken)
		return
	}

	//2.计算exp时间
	tokenExp, exist := c.Get("exp")
	if !exist {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "token is not exp time")
		cg.Res(http.StatusBadRequest, controller.ErrTokenNotExistTime)
		return
	}
	tokenBlackExp := time.Unix(int64(math.Round(tokenExp.(float64))), 0).Sub(time.Now())

	//3.将token放入redis中
	err := userapi.AddTokenToBlackList(token.(string), tokenBlackExp)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "add token to black list error", err)
		cg.Res(http.StatusBadRequest, controller.ErrTokenAddBlackList)
		return
	}
	cg.Res(http.StatusOK, controller.ErrnoSuccess)
}
