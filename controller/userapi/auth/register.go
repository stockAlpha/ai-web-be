package auth

import (
	"github.com/gin-gonic/gin"
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
// @Summary	用户注册
// @Accept		json
// @Produce	json
// @param		req	body		user.RegisterRequest	true	"注册请求参数"
// @Success	200	{string}	string						"返回token"
// @Router		/api/v1/user/register [post]
func Register(c *gin.Context) {
	cg := controller.Gin{Ctx: c}

	var req user.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	// 验证邮箱格式
	if req.Email == "" || !utils.IsEmailValid(req.Email) {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "email is out of specification")
		cg.Res(http.StatusBadRequest, controller.ErrNotFormatEmail)
		return
	}

	// 校验密码格式
	if !utils.IsValidPasswordFormat(req.Password) {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "password is out of specification")
		cg.Res(http.StatusBadRequest, controller.ErrNotFormatPassword)
		return
	}

	// 验证当前邮箱是否已注册
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

	// 验证code是否存在
	existCode, err := userapi.ExistCode(req.Code, req.Email)
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

	// 对密码进行加密,并添加用户
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "compute hash password err")
		cg.Res(http.StatusBadRequest, controller.ErrComputeHashPassword)
		return
	}

	userId, err := userapi.AddUser(req.Email, hashPassword)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "add user error")
		cg.Res(http.StatusBadRequest, controller.ErrAddUser)
		return
	}

	// todo 事务
	// 新注册用户赠送10个积分
	// 判断是否为被邀请用户，如果是则赠送20个积分，并且给邀请人也赠送10个积分
	inviteCode := req.InviteCode
	addAmount := 10
	if inviteCode != "" {
		inviteUser, err := userapi.GetUserByInviteCode(inviteCode)
		if err != nil {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query user by invite code error")
		} else {
			if inviteUser != nil {
				// 邀请人增加10的积分
				fromUserId := inviteUser.ID
				userapi.AddUserIntegral(fromUserId, 10)
				// 插入邀请关系
				userapi.AddInviteRelation(fromUserId, userId, inviteCode)
			} else {
				addAmount += 10
			}
		}
	}
	userapi.CreateUserIntegral(userId, addAmount)
	// 对userId, email加入jwt信息中
	token, err := userapi.GenerateToken(strconv.FormatUint(userId, 10), req.Email)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "generate token error")
		cg.Res(http.StatusBadRequest, controller.ErrGenerateJwtToken)
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, token)
}
