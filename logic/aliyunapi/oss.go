package aliyunapi

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"io"
	"net/http"
	"stock-web-be/client/ossclient"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"time"
)

func UploadFileByUrl(url, mimeType string) string {
	bucket := ossclient.GetOssBucket()
	resp, err := http.Get(url)
	if err != nil {
		tlog.Handler.Errorf(nil, consts.SLTagHTTPFailed, "get url error url=%s", url)
		return ""
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	fileName := fmt.Sprintf("image/img-%s-%s", uuid.New().String(), time.Now().Format("2006-01-02"))
	err = bucket.PutObject(fileName, bytes.NewReader(content), oss.ContentType(mimeType))
	if err != nil {
		tlog.Handler.Errorf(nil, consts.SLTagHTTPFailed, "put object error url=%s", url, err.Error())
		return ""
	}
	return "https://chatalpha.oss-cn-beijing.aliyuncs.com/" + fileName + "?x-oss-process=image/resize,w_500"
}
