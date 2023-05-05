package utils

import (
	"fmt"
	"testing"
)

func TestIsEmailValid(t *testing.T) {
	fmt.Println(IsEmailValid("stalary@163.com"))
	fmt.Println(IsEmailValid("stalary@162.com"))
	fmt.Println(IsEmailValid("stalary@gmail.com"))
	fmt.Println(IsEmailValid("stalary@qq.com"))
	fmt.Println(IsEmailValid("stalary@chacuo.net"))
}

func TestContainsChinese(t *testing.T) {
	fmt.Println(ContainsChinese("stalary"))
	fmt.Println(ContainsChinese("stalary的博客"))
	fmt.Println(ContainsChinese("你好，世界"))
}
