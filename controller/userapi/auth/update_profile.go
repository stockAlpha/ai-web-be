package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/user"
	"stock-web-be/logic/userapi"
	"strconv"
)

// @Tags	用户相关接口
// @Summary	修改用户信息
// @Router		/api/v1/user/profile [post]
// @param		req	body		user.ProfileRequest	true	"用户信息"
func UpdateProfile(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	req := user.ProfileRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	userapi.UpdateUser(userId, req.NickName, req.Avatar)
	cg.Res(http.StatusOK, controller.ErrnoSuccess)
}
