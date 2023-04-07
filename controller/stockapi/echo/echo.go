package echo

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/controller"
	"stock-web-be/idl/stockapi"
)

func Echo(ctx *gin.Context) {
	cg := controller.Gin{Ctx: ctx}

	var err error
	req := &stockapi.EchoReq{}

	if err = ctx.ShouldBindJSON(req); err != nil {
		controller.EchoError(ctx, 1, "Params illegal! "+err.Error())
		return
	}
	echoResponse := &stockapi.EchoResponse{}
	echoResponse.Text = "hello world," + req.Msg
	cg.Resp(http.StatusOK, controller.ErrnoSuccess, echoResponse)
}
