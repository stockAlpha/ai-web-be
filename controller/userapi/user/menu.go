package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/idl/userapi/menu"
)

// @Tags	用户相关接口
// @Summary	获取主菜单信息
// @Router		/api/v1/user/menu [get]
func Menu(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var toolItems []menu.Item
	toolItems = append(toolItems, menu.Item{Type: "chat", Name: "AI助手"})
	toolItems = append(toolItems, menu.Item{Type: "image", Name: "文生图"})
	toolItems = append(toolItems, menu.Item{Type: "chat", Name: "写作"})
	toolItems = append(toolItems, menu.Item{Type: "chat", Name: "计算"})

	var roleItems []menu.Item
	roleItems = append(roleItems, menu.Item{Type: "chat", Name: "程序员"})
	roleItems = append(roleItems, menu.Item{Type: "chat", Name: "小说家"})
	roleItems = append(roleItems, menu.Item{Type: "chat", Name: "数学家"})
	roleItems = append(roleItems, menu.Item{Type: "chat", Name: "导游"})
	cg.Res(http.StatusOK, controller.ErrnoSuccess)
}
