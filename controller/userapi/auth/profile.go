package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/idl/userapi/user"
	"stock-web-be/logic/userapi"
	"strconv"
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
	userProfile, _ := userapi.GetUserByEmail(email)
	userIntegral, _ := userapi.GetUserIntegralByUserId(userId)
	res.Integral = userIntegral.Amount
	res.NickName = userProfile.NickName
	res.Avatar = userProfile.Avatar
	res.InviteCode = userProfile.InviteCode
	res.Email = email
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, res)
}
