package utils

import (
	"math/rand"
	"time"
)

// 生成随机验证码
func GenerateCode() string {
	// 定义可用于生成验证码的字符集
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// 设置随机数生成器的种子
	rand.Seed(time.Now().UTC().UnixNano())

	// 生成一个包含 6 个字符的随机字符串
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}

	return string(code)
}
