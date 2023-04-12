// logger-中间件，打印access日志
// Gin：
//
//	appEngine := gin.New()
//
// PrintResponseLen为响应打印长度，设置为0则不打印
//
//	appEngine.Use(log.GinLogger(log.LoggerConfig{PrintResponseLen: 0}))
package middleware

import (
	"bytes"
	"io/ioutil"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"time"

	"github.com/gin-gonic/gin"
)

//const (
//	REFERER   = "referer"
//	COOKIE    = "cookie"
//	CLIENT_IP = "client_ip"
//	LOCAL_IP  = "local_ip"
//	MODULE    = "module"
//	UA        = "ua"
//	HOST      = "host"
//	URI       = "uri"
//)

type LoggerConfig struct {
	// responsebody 打印长度
	PrintResponseLen int
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	if w.body != nil {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

// access日志打印
// logId串联请求链路
func GinLogger(config LoggerConfig) gin.HandlerFunc {
	// 当前模块名
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		c.Set(consts.TimeMilliRequestIn, start.UnixMilli())
		// 请求报文
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		//logid维护
		logid := c.GetHeader(tlog.LOGID)
		if len(logid) <= 0 {
			logid = tlog.GenLogId()
			c.Request.Header.Add(tlog.LOGID, logid)
		}
		c.Set(tlog.LOGID, logid)

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// args
		var args string
		switch c.Request.Method {
		case "GET":
			args = c.Request.URL.RawQuery
		case "POST":
			args = string(body)
		}

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		// 执行时间 单位:毫秒
		latency := end.Sub(start).Milliseconds()

		tlog.Handler.Accessf(c, consts.SLTagRequest,
			"method=%s||uri=%s||args=%s||errno=%d||response=%s||proc_time=%v",
			c.Request.Method,
			c.Request.URL.Path,
			args,
			c.Writer.Status(),
			blw.body.String(),
			latency)
	}
}
