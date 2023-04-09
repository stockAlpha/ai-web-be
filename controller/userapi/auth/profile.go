package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/idl/userapi/user"
)

// @Tags	用户相关接口
// @Summary	获取用户信息
// @Router		/api/v1/user/profile [get]
// @Response 200 {object} user.ProfileResponse 用户信息
func Profile(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	res := user.ProfileResponse{}
	res.Email = c.GetString("email")
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, res)
}
