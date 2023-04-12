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
	"encoding/json"
	"io/ioutil"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"strings"
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
	// 本地IP
	var localIP string = ""
	// 当前模块名
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		c.Set(consts.TimeMilliRequestIn, start.UnixMilli())
		//// 请求url
		//path := c.Request.URL.Path
		//raw := c.Request.URL.RawQuery
		//if raw != "" {
		//	path = path + "?" + raw
		//}
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

		// header
		header, _ := json.Marshal(c.Request.Header)

		// args
		var args string
		switch c.Request.Method {
		case "GET":
			args = c.Request.URL.RawQuery
		case "POST":
			args = string(body)
		}

		// 非OPTIONS请求打印日志
		if c.Request.Method != "OPTIONS" {
			tlog.Handler.Accessf(c, consts.SLTagRequestIn,
				"bff-x-request-id=%s||method=%s||uri=%s||proto=%s||refer=%s||cookie=%s||client_ip=%s||local_ip=%s||user_agent=%s||content_type=%s||header=%s||args=%s||errno=0||response=||proc_time=0",
				c.GetHeader("bff-x-request-id"),
				c.Request.Method,
				c.Request.URL.Path,
				c.Request.Proto,
				c.Request.Referer(),
				c.Request.Cookies(),
				c.ClientIP(),
				localIP,
				c.Request.UserAgent(),
				c.Request.Header.Get(consts.ContentType),
				string(header),
				args)
		}

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		// 执行时间 单位:毫秒
		latency := end.Sub(start).Milliseconds()

		// 统一access log格式便于解析
		// /pipeline/.*/log接口不用打印response
		if c.Request.URL.Path == `/metrics` || strings.HasSuffix(c.Request.URL.Path, "/log") {
			tlog.Handler.Accessf(c, consts.SLTagRequestOut,
				"bff-x-request-id=%s||method=%s||uri=%s||proto=||refer=||cookie=||client_ip=||local_ip=||user_agent=||content_type=||header=||args=||errno=%d||response=||proc_time=%v",
				c.GetHeader("bff-x-request-id"),
				c.Request.Method,
				c.Request.URL.Path,
				c.Writer.Status(),
				latency)
		} else {
			tlog.Handler.Accessf(c, consts.SLTagRequestOut,
				"bff-x-request-id=%s||method=%s||uri=%s||proto=||refer=||cookie=||client_ip=||local_ip=||user_agent=||content_type=||header=||args=||errno=%d||response=%s||proc_time=%v",
				c.GetHeader("bff-x-request-id"),
				c.Request.Method,
				c.Request.URL.Path,
				c.Writer.Status(),
				blw.body.String(),
				latency)
		}
	}
}
