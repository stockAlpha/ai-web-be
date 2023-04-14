package user

import (
	"net/http"
	"strconv"

	"stock-web-be/async"
	"stock-web-be/controller"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/user"
	"stock-web-be/logic/userapi"

	"github.com/gin-gonic/gin"
)

// @Tags	用户相关接口
// @Summary	意见反馈
// @Router		/api/v1/user/feedback [post]
// @param		req	body		user.FeedbackRequest	true	"反馈信息"
func Feedback(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	req := user.FeedbackRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	email := c.GetString("email")
	_ = userapi.AddFeedback(userId, req.FeedbackType, req.Content)
	// 1-问题反馈 2-功能建议 3-咨询 4-其他
	title := "收到新的反馈："
	switch req.FeedbackType {
	case 1:
		title += "问题反馈"
	case 2:
		title += "功能建议"
	case 3:
		title += "咨询"
	case 4:
		title += "其他"
	}
	ret := "用户ID：" + strconv.FormatUint(userId, 10) + "\n" + "邮箱：" + email + "\n" + "反馈内容：" + req.Content
	async.MailChan <- async.MailChanType{To: "stalary@163.com", Subject: title, Body: ret}
	cg.Res(http.StatusOK, controller.ErrnoSuccess)
}
