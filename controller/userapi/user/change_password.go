package user

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/user"
	"stock-web-be/logic/userapi"
	"stock-web-be/utils"
)

// @Tags	用户相关接口
// @Summary	在忘记密码时发送验证码
// @param		req	body	user.SendPasswordVerificationCodeRequest	true	"发送验证码请求参数(默认为email)"
// @Router		/api/v1/user/change_password/verify/code [post]
func SendPasswordVerificationCode(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req user.SendPasswordVerificationCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	//判断传入的subject type是否合法
	if req.SubjectType != userapi.ChangePasswordMailCode {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "subject type is invalid")
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	authType := userapi.VerificationCodeTypeToAuthType(req.SubjectType)

	user, err := userapi.GetUserByAuthType(req.SubjectName, authType)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query user is fail, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoError)
		return
	}
	//判断最终的用户是否为空
	if user == nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query user is fail, error: %s")
		cg.Res(http.StatusBadRequest, controller.ErrUserSubjectIdNotFound)
		return
	}

	//向用户发送验证码
	err = userapi.SendGeneralVerificationCode(req.SubjectName, req.SubjectType)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query user is fail, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoError)
		return
	}

	cg.Resp(http.StatusOK, controller.ErrnoSuccess, true)
}

// @Tags	用户相关接口
// @Summary	用户修改密码
// @param		req	body	user.ChangePasswordRequest	true	"发送验证码请求参数(默认为email)"
// @Router		/api/v1/user/change_password [post]
func ChangePassword(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req *user.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	//判断传入的subject type是否合法
	if req.SubjectType != userapi.ChangePasswordMailCode {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "subject type is invalid")
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}

	authType := userapi.VerificationCodeTypeToAuthType(req.SubjectType)

	//判断新密码的格式
	if !utils.IsValidPasswordFormat(req.NewPassword) {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "password is out of specification")
		cg.Res(http.StatusBadRequest, controller.ErrNotFormatPassword)
		return
	}

	oldUser, err := userapi.GetUserByAuthType(req.SubjectName, authType)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query oldUser is fail, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoError)
		return
	}
	//判断最终的用户是否为空
	if oldUser == nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query oldUser is fail, error: %s")
		cg.Res(http.StatusBadRequest, controller.ErrUserSubjectIdNotFound)
		return
	}

	//判断用户旧密码是否和新密码一致
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(oldUser.Password), []byte(req.NewPassword))
	//代表计算hash值出现错误
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "compute hash password error", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrComputeHashPassword)
		return
	}
	//代表密码比对是一致的
	if err == nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "new password is same to old password")
		cg.Res(http.StatusBadRequest, controller.ErrPasswordNotChange)
		return
	}

	// 验证code是否存在
	existCode, err := userapi.ExistCodeByCodeType(req.VerificationCode, req.SubjectName, req.SubjectType)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query code by email is fatal", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrQueryVerificationCode)
		return
	}

	if !existCode {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "the code is not exist or expire")
		cg.Res(http.StatusBadRequest, controller.ErrVerificationCodeNotFound)
		return
	}

	//更新密码
	//计算密码hash值
	hashPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "compute hash password err", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrComputeHashPassword)
		return
	}

	//更新密码同时失效code
	err = transactionChangePassword(c, oldUser.ID, hashPassword, req)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "change password transaction err", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrChangePassword)
		return
	}

	cg.Resp(http.StatusOK, controller.ErrnoSuccess, true)
}

func transactionChangePassword(c *gin.Context, userId uint64, hashPassword string, req *user.ChangePasswordRequest) error {
	// 声明个db，做事务回滚
	tx := db.DbIns.Begin()

	//更新用户密码
	err := userapi.UpdateUserPassword(userId, hashPassword, tx)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "update user password err", err.Error())
		tx.Rollback()
		return err
	}

	//将code失效
	err = userapi.ExpireCode(req.VerificationCode, req.SubjectName, req.SubjectType, tx)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "expire code err", err.Error())
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
