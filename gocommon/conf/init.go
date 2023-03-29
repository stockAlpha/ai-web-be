package conf

import (
	"fmt"
	"os"

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
		handler.SetConfigName("app") // 默认文件配置文件为app.toml
		handler.AddConfigPath(Root + "/conf")
	} else {
		handler.SetConfigFile(confPath)
	}
	err := handler.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	return handler
}

// FreshHandler 从磁盘读取配置文件，方便本地修改配置的测试，无需重启服务
func FreshHandler() {
	Handler.ReadInConfig()
}
