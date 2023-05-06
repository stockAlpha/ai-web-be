package xfapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"strings"
	"time"
)

var (
	host      = "itrans.xfyun.cn"
	Algorithm = "hmac-sha256"
	HttpProto = "HTTP/1.1"
	uri       = "/v2/its"
	url       = "http://itrans.xfyun.cn/v2/its"
)

func Run(text, from, to string) string { //生成client 参数为默认
	client := &http.Client{}

	//提交请求
	var data1 = []byte(text)
	param := map[string]interface{}{
		"common": map[string]interface{}{
			"app_id": conf.Handler.GetString("xf.app_id"), //appid 必须带上，只需第一帧发送
		},
		"business": map[string]interface{}{ //business 参数，只需一帧发送
			"from": from, //源语种
			"to":   to,   //目标语种
		},
		"data": map[string]interface{}{
			"text": base64.StdEncoding.EncodeToString(data1),
		},
	}
	tt, _ := json.Marshal(param)
	jsons := string(tt)
	jsoninfos := strings.NewReader(jsons)
	reqest, err := http.NewRequest("POST", url, jsoninfos)
	if err != nil {
		panic(err)
	}
	//增加header选项
	reqest.Header.Set("Content-Type", "application/json")
	reqest.Header.Set("Host", host)
	reqest.Header.Set("Accept", "application/json,version=1.0")
	currentTime := time.Now().UTC().Format(time.RFC1123)
	reqest.Header.Set("Date", currentTime)
	digest := "SHA-256=" + signBody(string(tt))
	reqest.Header.Set("Digest", digest)
	// 根据请求头部内容，生成签名
	sign := generateSignature(host, currentTime, "POST", uri, HttpProto, digest, conf.Handler.GetString("xf.api_secret"))
	// 组装Authorization头部
	authHeader := fmt.Sprintf(`api_key="%s", algorithm="%s", headers="host date request-line digest", signature="%s"`, conf.Handler.GetString("xf.api_key"), Algorithm, sign)
	reqest.Header.Set("Authorization", authHeader)
	//处理返回结果
	response, _ := client.Do(reqest)
	defer response.Body.Close()
	var result Result
	body, _ := io.ReadAll(response.Body)
	_ = json.Unmarshal(body, &result)
	if result.Code != 0 {
		tlog.Handler.Errorf(nil, consts.SLTagXFError, "科大讯飞翻译失败 text=%s, from=%s, to=%s, error=%s", text, from, to, result.Message)
		return text
	}
	tlog.Handler.Infof(nil, consts.SLTagXFSuccess, "科大讯飞翻译成功 text=%s, from=%s, to=%s, result=%s", text, from, to, result.Data.Result.TransResult.Dst)
	return result.Data.Result.TransResult.Dst
}

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Result struct {
			From        string `json:"from"`
			To          string `json:"to"`
			TransResult struct {
				Dst string `json:"dst"`
				Src string `json:"src"`
			} `json:"trans_result"`
		}
	} `json:"data"`
}

func generateSignature(host, date, httpMethod, requestUri, httpProto, digest string, secret string) string {
	// 不是request-line的话，则以header名称,后跟ASCII冒号:和ASCII空格，再附加header值
	var signatureStr string
	if len(host) != 0 {
		signatureStr = "host: " + host + "\n"
	}
	signatureStr += "date: " + date + "\n"
	// 如果是request-line的话，则以 http_method request_uri http_proto
	signatureStr += httpMethod + " " + requestUri + " " + httpProto + "\n"
	signatureStr += "digest: " + digest
	return hmacsign(signatureStr, secret)
}
func hmacsign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}
func signBody(data string) string {
	// 进行sha256签名
	//fmt.Println(data)
	sha := sha256.New()
	sha.Write([]byte(data))
	encodeData := sha.Sum(nil)
	// 经过base64转换
	return base64.StdEncoding.EncodeToString(encodeData)
}
