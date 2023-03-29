package echo

import (
	"github.com/gin-gonic/gin"
	"stock-web-be/controller"
	"stock-web-be/idl/stockapi"
)

func Echo(ctx *gin.Context) {
	var err error
	req := &stockapi.EchoReq{}
	resp := controller.Resp{
		0,
		"",
		nil,
	}

	if err = ctx.ShouldBindJSON(req); err != nil {
		controller.EchoError(ctx, 1, "Params illegal! "+err.Error())
		return
	}
	echoResponse := &stockapi.EchoResponse{}
	echoResponse.Text = "hello world," + req.Msg
	resp.Data = echoResponse
	controller.EchoJSON(ctx, resp)
}
