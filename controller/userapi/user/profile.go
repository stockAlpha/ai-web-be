package user

import (
	"encoding/json"
	"net/http"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"strconv"

	"stock-web-be/controller"
	"stock-web-be/idl/userapi/user"
	"stock-web-be/logic/userapi"

	"github.com/gin-gonic/gin"
)

// @Tags	用户相关接口
// @Summary	获取用户信息
// @Router		/api/v1/user/profile [get]
// @Response 200 {object} user.ProfileResponse 用户信息
func Profile(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	res := user.ProfileResponse{}
	email := c.GetString("email")
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	userProfile, err := userapi.GetUserById(userId)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "get user email error", err.Error())
		cg.Resp(http.StatusBadRequest, controller.ErrServer, res)
		return
	}
	userIntegral, err := userapi.GetUserIntegralByUserId(userId, nil)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "get user integral error", err.Error())
		cg.Resp(http.StatusBadRequest, controller.ErrServer, res)
		return
	}
	res.Integral = userIntegral.Amount
	res.NickName = userProfile.NickName
	res.Avatar = userProfile.Avatar
	res.InviteCode = userProfile.InviteCode
	res.Email = email
	res.VipUser = userProfile.VipUser
	json.Unmarshal(userProfile.CustomConfig, &res.CustomConfig)
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, res)
}
