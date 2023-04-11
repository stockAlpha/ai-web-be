package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	brokePipeErr = "write: broken pipe"
	connRSTErr   = "write: connection reset by peer"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				r := c.Request
				h := r.Header
				tag := consts.SLTagPanic
				if strings.Contains(fmt.Sprintf("%v", err), brokePipeErr) {
					tag = consts.SLTagBrokePipe
				}
				if strings.Contains(fmt.Sprintf("%v", err), connRSTErr) {
					tag = consts.SLTagConnRST
				}
				tlog.Handler.Panicf(c, tag, "Method=%s||Host=%s||Url=%s||CallerUri=%s||Module=%s||Idc=%s"+
					"||ContentType=%s||ContentLength=%s||UserAgent=%s||Product=%s"+
					"||SpanId=%s||Subsys=%s||UniqId=%s||UserIp=%s"+
					"||Proto=%s||RemoteAddr=%s||Err=%v||Stack=\r\n%v ",
					r.Method, r.Host, r.URL, h.Get("X_bd_caller_uri"), h.Get("X_bd_module"), h.Get("X_bd_idc"),
					h.Get("Content-Type"), h.Get("Content-Length"), h.Get("User-Agent"), h.Get("X_bd_product"),
					h.Get("X_bd_spanid"), h.Get("X_bd_subsys"), h.Get("X_bd_uniqid"), h.Get("X_bd_userip"),
					r.Proto, r.RemoteAddr, err, string(debug.Stack()))
				c.String(http.StatusOK, `{"code":-100,"msg":"panic","data":{}}`)
				c.Abort()
			}
		}()
		c.Next()
	}
}
