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
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/utils"

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
			if len(string(body)) <= 10000 {
				args = string(body)
			}
		}
		//尝试脱敏
		args = makeSensitive(args)
		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		// 执行时间 单位:毫秒
		latency := end.Sub(start).Milliseconds()

		res := blw.body.String()
		if len(res) > 1000 {
			res = ""
		}
		// 没接nginx，主页/favicon.ico无需记录
		if !(c.Request.URL.Path == "/favicon.ico" && strings.ToLower(c.Request.Method) == "get") {
			email := c.GetString("email")
			userId, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
			tlog.Handler.Accessf(c, consts.SLTagRequest,
				"method=%s||uri=%s||userId=%d||email=%s||args=%s||errno=%d||response=%s||proc_time=%v",
				c.Request.Method,
				c.Request.URL.Path,
				userId,
				email,
				args,
				c.Writer.Status(),
				res,
				latency)
		}
	}
}

var (
	//PasswordReg = regexp.MustCompile(`^(.*\"password\")(\s*:\s*\")(.*)(\".*)$`)
	passwordReg = regexp.MustCompile(`"password"\s*:\s*"(.*?)"`)
)

// makeSensitive 日志脱敏
func makeSensitive(source string) (destination string) {
	// password脱敏
	destination = passwordReg.ReplaceAllStringFunc(source, func(s string) string {
		pass := passwordReg.FindStringSubmatch(s)[1]
		hash, err := utils.HashPassword(pass)
		if err != nil {
			log.Println("error in hash", err)
			return source
		}
		return fmt.Sprintf(`"password":"%s"`, "len:"+strconv.Itoa(len(pass))+",hash:"+hash)
	})
	//其他的脱敏也可以放这里
	return destination
}
