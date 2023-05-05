package openaiapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/logic/userapi"
	"strconv"
)

// @Tags	OpenAI相关接口
// @Summary	获取可用模型
// @Router		/api/v1/openai/v1/model [get]
// @Response 200 {object} map[string][]string 模型列表
func Model(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	userProfile, _ := userapi.GetUserById(userId)
	imageModel := []string{"dall-e2"}
	// vip用户可以优先体验sd
	if userProfile.VipUser {
		imageModel = []string{"dall-e2", "stable-diffusion"}
	}
	model := make(map[string][]string)
	model["chat"] = []string{"gpt-3.5-turbo"}
	model["image"] = imageModel
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, model)
}
