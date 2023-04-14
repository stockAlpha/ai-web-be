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
				tag := consts.SLTagPanic
				if strings.Contains(fmt.Sprintf("%v", err), brokePipeErr) {
					tag = consts.SLTagBrokePipe
				}
				if strings.Contains(fmt.Sprintf("%v", err), connRSTErr) {
					tag = consts.SLTagConnRST
				}
				tlog.Handler.Panicf(c, tag, "Method=%s||Host=%s||Url=%s||Err=%v||Stack=\r\n%v ",
					r.Method, r.Host, r.URL, err, string(debug.Stack()))
				c.String(http.StatusOK, `{"code":-100,"msg":"panic","data":{}}`)
				c.Abort()
			}
		}()
		c.Next()
	}
}
