package user

import (
	"fmt"
	"net/http"
	"stock-web-be/async"
	"strconv"

	"stock-web-be/controller"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/user"
	"stock-web-be/logic/userapi"
	"stock-web-be/utils"

	"github.com/gin-gonic/gin"
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

	inviteCode := req.InviteCode
	userId, err := transactionRegister(c, req.Email, hashPassword, inviteCode)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "register error", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrServer)
		return
	}
	// 对userId, email加入jwt信息中
	token, err := userapi.GenerateToken(strconv.FormatUint(userId, 10), req.Email)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "generate token error")
		cg.Res(http.StatusBadRequest, controller.ErrGenerateJwtToken)
		return
	}
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, token)
}

// 开启事务注册
func transactionRegister(c *gin.Context, email, hashPassword, inviteCode string) (uint64, error) {
	// 声明个db，做事务回滚
	tx := db.DbIns.Begin()
	// 注册新用户
	userId, err := userapi.AddUser(email, hashPassword, tx)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "add user error")
		tx.Rollback()
		return 0, err
	}

	// 新注册用户赠送20个积分
	// 判断是否为被邀请用户，如果是则邀请人增加积分
	addAmount := 20
	inviteUser, err := userapi.GetUserByInviteCode(inviteCode, tx)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query user by invite code error")
		tx.Rollback()
		return 0, err
	}
	if inviteUser != nil {
		// 邀请人增加积分
		fromAddAmount := 10
		fromUserId := inviteUser.ID
		fromEmail := inviteUser.Email
		err := userapi.AddUserIntegral(fromUserId, fromAddAmount, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		// 插入邀请关系
		err = userapi.AddInviteRelation(fromUserId, userId, inviteCode, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		async.MailChan <- async.MailChanType{To: fromEmail, Subject: consts.InviteSubject, Body: fmt.Sprintf(consts.InviteContent, email, fromAddAmount)}
	}
	_, err = userapi.CreateUserIntegral(userId, addAmount, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return userId, nil
}
