package xfapi

import (
	"fmt"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/tlog"
	"testing"
)

func TestRun(t *testing.T) {
	conf.Init("../../conf/app.local.toml")
	tlog.Init()
	fmt.Println(Run("你好，我要一份鱼香肉丝", "cn", "en"))
}
