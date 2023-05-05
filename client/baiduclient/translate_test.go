package baiduclient

import (
	"fmt"
	"testing"
)

func TestTranslate(t *testing.T) {
	baidu := &BaiDu_translate{
		appid:     "20230504001666315",
		secretKey: "KR5vjtqEm6w2C2foeG2A",
	}
	run, err := baidu.Run("一辆黑色的汽车", "auto", "en")
	if err != nil {
		return
	}
	fmt.Println("run ", run.TransResults[0].Dst)
}
