package ossclient

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
)

var bucket *oss.Bucket

func GetOssBucket() *oss.Bucket {
	return bucket
}

func Init() {
	client, err := oss.New("oss-cn-beijing.aliyuncs.com", "LTAI5tQ7pCSV7aQcJjaXrdMG", "4kIQWUll9JKZ3F3kf8IEvXeF5X5DgR")
	if err != nil {
		tlog.Handler.Errorf(nil, consts.SLTagAliOssSuccess, "ali oss get client error", err.Error())
	}
	bucket, err = client.Bucket("chatalpha")
	if err != nil {
		tlog.Handler.Errorf(nil, consts.SLTagAliOssSuccess, "ali oss get bucket error", err.Error())
	}
	tlog.Handler.Infof(nil, consts.SLTagAliOssSuccess, "ali oss get bucket success")
}
