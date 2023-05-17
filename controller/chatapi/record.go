package chatapi

import (
	"net/http"

	"stock-web-be/controller"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/idl/userapi/record"

	"github.com/gin-gonic/gin"
)

// @Tags	用户记录
// @Summary	获取用户record
// @param			user_id		query		uint		true	"user_id"
//
//	@response		200		{object}	openai.ChatRecordResponse
//
// @Router		/api/v1/chat/record [get]
func GetChatRecord(c *gin.Context) {
	cg := controller.Gin{Ctx: c}
	var req record.ChatRecordRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "request params invalid, error: %s", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrnoInvalidPrm)
		return
	}
	chats, err := db.FindRecord(req)
	if err != nil {
		tlog.Handler.Errorf(c, consts.SLTagHTTPFailed, "query error", err.Error())
		cg.Res(http.StatusBadRequest, controller.ErrServer)
		return
	}
	resp := fmtChatRecord(chats)
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, resp)
}

func fmtChatRecord(records []db.ChatRecord) (resp record.ChatRecordResponse) {
	chatCacheMap := make(map[int][]record.ChatRecordChatData)
	for i := range records {
		if _, ok := chatCacheMap[records[i].UUID]; ok {
			//append 有序的
			chatCacheMap[records[i].UUID] = append(chatCacheMap[records[i].UUID], records[i].DbToOpenAIData())
		} else {
			chatCacheMap[records[i].UUID] = make([]record.ChatRecordChatData, 0)
		}
		//Active为最后一个
		if i == len(records)-1 {
			resp.Active = records[i].UUID
		}
	}
	for k, v := range chatCacheMap {
		resp.Chat = append(resp.Chat, record.ChatRecordChat{UUID: k, Data: v})
	}
	return
}
