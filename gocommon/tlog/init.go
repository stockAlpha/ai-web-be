package tlog

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"stock-web-be/gocommon/conf"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ILog log interface
type ILog interface {
	Debugf(ctx context.Context, tag string, format string, args ...interface{})
	Infof(ctx context.Context, tag string, format string, args ...interface{})
	Warnf(ctx context.Context, tag string, format string, args ...interface{})
	Errorf(ctx context.Context, tag string, format string, args ...interface{})
	Fatalf(ctx context.Context, tag string, format string, args ...interface{})
	Panicf(ctx context.Context, tag string, format string, args ...interface{})
	Metricf(ctx context.Context, tag string, format string, args ...interface{})
	Accessf(ctx context.Context, tag string, format string, args ...interface{})
}

// Handler zapLog
var Handler ILog

const (
	LOGID    = "logid"
	NOTICES  = "notices"
	RESPONSE = "response"

	WfSuffix      = ".wf"
	MetricSuffix  = ".metric"
	AccessSuffix  = ".access"
	LogFileFormat = ".%Y%m%d%H"
)

// Init Init
func Init() {
	var err error

	// 读取配置
	//filename := conf.Handler.GetString("log.filename") // "log/info.log.%Y%m%d%H"
	filePrefix := conf.Handler.GetString("log.filePrefix") // log/info
	fileSuffix := conf.Handler.GetString("log.fileSuffix") // .log
	maxHourAge := time.Duration(conf.Handler.GetInt("log.maxHourAge"))
	maxHourRotate := time.Duration(conf.Handler.GetInt("log.maxHourRotate"))

	// 文件拆分规则
	// info日志
	rlogs, err := rotatelogs.New(
		//conf.Root+"/"+filename, // 日志采集对软连有要求，不能是绝对路径
		filePrefix+fileSuffix+LogFileFormat, // 日志采集对软连有要求，不能是绝对路径
		rotatelogs.WithLinkName(filePrefix+fileSuffix),
		rotatelogs.WithMaxAge(maxHourAge*time.Hour),
		rotatelogs.WithRotationTime(maxHourRotate*time.Hour),
	)
	if err != nil {
		log.Fatal("Init zap log error: ", err)
	}

	// wf日志
	wfRlogs, err := rotatelogs.New(
		filePrefix+WfSuffix+fileSuffix+LogFileFormat, // 日志采集对软连有要求，不能是绝对路径
		rotatelogs.WithLinkName(filePrefix+WfSuffix+fileSuffix),
		rotatelogs.WithMaxAge(maxHourAge*time.Hour),
		rotatelogs.WithRotationTime(maxHourRotate*time.Hour),
	)
	if err != nil {
		log.Fatal("Init zap log error: ", err)
	}

	// metric日志
	metricRlogs, err := rotatelogs.New(
		filePrefix+MetricSuffix+fileSuffix+LogFileFormat, // 日志采集对软连有要求，不能是绝对路径
		rotatelogs.WithLinkName(filePrefix+MetricSuffix+fileSuffix),
		rotatelogs.WithMaxAge(maxHourAge*time.Hour),
		rotatelogs.WithRotationTime(maxHourRotate*time.Hour),
	)
	if err != nil {
		log.Fatal("Init zap log error: ", err)
	}

	// access日志
	accessRlogs, err := rotatelogs.New(
		filePrefix+AccessSuffix+fileSuffix+LogFileFormat, // 日志采集对软连有要求，不能是绝对路径
		rotatelogs.WithLinkName(filePrefix+AccessSuffix+fileSuffix),
		rotatelogs.WithMaxAge(maxHourAge*time.Hour),
		rotatelogs.WithRotationTime(maxHourRotate*time.Hour),
	)
	if err != nil {
		log.Fatal("Init zap log error: ", err)
	}

	// 文件内容格式和之前保持一致
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = ""
	encoderConfig.TimeKey = ""

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(rlogs),
		getLogLevel(conf.Handler.GetString("log.level")),
	)
	wfCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(wfRlogs),
		getLogLevel(conf.Handler.GetString("log.level")),
	)
	metricCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(metricRlogs),
		getLogLevel(conf.Handler.GetString("log.level")),
	)
	accessCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(accessRlogs),
		getLogLevel(conf.Handler.GetString("log.level")),
	)

	Handler = NewTLog(zap.New(core), zap.New(wfCore), zap.New(metricCore), zap.New(accessCore))
}

func getLogLevel(level string) zapcore.Level {
	logLevel := zap.DebugLevel

	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		logLevel = zap.DebugLevel
	case "INFO":
		logLevel = zap.InfoLevel
	case "WARN":
		logLevel = zap.WarnLevel
	case "ERROR":
		logLevel = zap.ErrorLevel
	case "FATAL":
		logLevel = zap.FatalLevel
	}

	return logLevel
}

// TLog ...
type TLog struct {
	log       *zap.Logger
	wflog     *zap.Logger
	metriclog *zap.Logger
	accesslog *zap.Logger
}

// NewTLog ...
func NewTLog(log, wflog, metriclog, accesslog *zap.Logger) *TLog {
	return &TLog{log: log, wflog: wflog, metriclog: metriclog, accesslog: accesslog}
}

// Debugf ...
func (tl *TLog) Debugf(ctx context.Context, tag string, format string, args ...interface{}) {
	tl.format(ctx, "DEBUG", tag, format, args...)
}

// Infof ...
func (tl *TLog) Infof(ctx context.Context, tag string, format string, args ...interface{}) {
	tl.format(ctx, "INFO", tag, format, args...)
}

// Warnf ...
func (tl *TLog) Warnf(ctx context.Context, tag string, format string, args ...interface{}) {
	tl.format(ctx, "WARN", tag, format, args...)
}

// Errorf ...
func (tl *TLog) Errorf(ctx context.Context, tag string, format string, args ...interface{}) {
	tl.format(ctx, "ERROR", tag, format, args...)
}

// Fatalf ...
func (tl *TLog) Fatalf(ctx context.Context, tag string, format string, args ...interface{}) {
	tl.format(ctx, "FATAL", tag, format, args...)
}

// Panicf ...
func (tl *TLog) Panicf(ctx context.Context, tag string, format string, args ...interface{}) {
	tl.format(ctx, "PANIC", tag, format, args...)
}

// Metric ...
func (tl *TLog) Metricf(ctx context.Context, tag string, format string, args ...interface{}) {
	tl.format(ctx, "METRIC", tag, format, args...)
}

// Access ...
func (tl *TLog) Accessf(ctx context.Context, tag string, format string, args ...interface{}) {
	tl.format(ctx, "ACCESS", tag, format, args...)
}

func (tl *TLog) format(ctx context.Context, level, tag, format string, args ...interface{}) string {
	// time
	//ts := time.Now().Format("2006-01-02 15:04:05.000000")
	ts := time.Now().Format(time.RFC3339Nano)

	// file, line
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "file"
		line = -1
	}

	// igonre dir
	file = strings.TrimPrefix(file, conf.Root+"/")
	var logid string
	if ctx != nil {
		if !reflect.ValueOf(ctx).IsNil() {
			logid, _ = ctx.Value(LOGID).(string)
		}
	}
	// return fmt.Sprintf("[%s][%s][%s:%d] %s", level, ts, file, line, fmt.Sprintf(format, args...))
	fmt.Printf("%s %s %s:%d %s||logid=%s||%s\n", ts, level, file, line, tag, logid, fmt.Sprintf(format, args...))
	return fmt.Sprintf("%s %s %s:%d %s||logid=%s||%s", ts, level, file, line, tag, logid, fmt.Sprintf(format, args...))
}
