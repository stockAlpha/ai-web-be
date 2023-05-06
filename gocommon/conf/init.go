package conf

import (
	"fmt"
	"os"
	"strings"

	"stock-web-be/gocommon/consts"

	"github.com/spf13/viper"
)

// Root current dir
var Root string

// config handler
var Handler *viper.Viper

func Init(confPath string) {
	var err error
	//_, fn, _, _ := runtime.Caller(0)
	//Root = filepath.Dir(filepath.Dir(fn))
	Root, err = os.Getwd()
	if err != nil {
		panic(fmt.Errorf("Initialize Root error: %s", err))
	}

	Handler = LoadConfig(confPath)
}

func LoadConfig(confPath string) *viper.Viper {
	handler := viper.New()

	if confPath == "" {
		// local or prod
		confEnv := os.Getenv(consts.Env)
		if confEnv != "" {
			handler.SetConfigName("app." + confEnv)
		} else {
			handler.SetConfigName("app.local") // 默认文件配置文件为app.toml
			// 配置本地代理
			//proxyUrl := "http://127.0.0.1:7890"
			//proxyURL, err := url.Parse(proxyUrl)
			//if err != nil {
			//	panic(err)
			//}
			//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			//os.Setenv("HTTP_PROXY", proxyUrl)
			//os.Setenv("HTTPS_PROXY", proxyUrl)
		}
		handler.AddConfigPath(Root + "/conf")
	} else {
		handler.SetConfigFile(confPath)
	}
	err := handler.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	handler.AutomaticEnv()
	for _, key := range handler.AllKeys() {
		if flag := os.Getenv(key); flag != "" {
			handler.Set(key, flag)
			continue
		}
		upperKey := strings.ToUpper(key)
		if flag := os.Getenv(upperKey); flag != "" {
			handler.Set(key, flag)
			continue
		}
		newKey := strings.Replace(key, ".", "_", -1)
		if flag := os.Getenv(newKey); flag != "" {
			handler.Set(key, flag)
			continue
		}
		newKey = strings.Replace(upperKey, ".", "_", -1)
		if flag := os.Getenv(newKey); flag != "" {
			handler.Set(key, flag)
			continue
		}
	}
	return handler
}

// FreshHandler 从磁盘读取配置文件，方便本地修改配置的测试，无需重启服务
func FreshHandler() {
	Handler.ReadInConfig()
}
