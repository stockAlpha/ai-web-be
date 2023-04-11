package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/user"
	"stock-web-be/logic/userapi"
)

// @Tags	用户相关接口
// @Summary	获取用户信息
// @Router		/api/v1/user/profile [get]
// @Response 200 {object} user.ProfileResponse 用户信息
func Profile(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	res := user.ProfileResponse{}
	email := c.GetString("email")
	userProfile, _ := userapi.GetUserProfileByEmail(email)
	userId := userProfile.ID
	userIntegral, _ := userapi.GetUserIntegralByUserId(userId)
	if userIntegral == nil {
		integral, err := userapi.CreateUserIntegral(userId)
		if err != nil {
			tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "add integral error")
			cg.Res(http.StatusBadRequest, controller.ErrAddIntegral)
			return
		}
		res.Integral = integral.Amount
	} else {
		res.Integral = userIntegral.Amount
	}
	res.NickName = userProfile.NickName
	res.Email = email
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, res)
}
