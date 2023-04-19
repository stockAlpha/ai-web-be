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
