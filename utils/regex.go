package utils

import "regexp"

// 验证邮箱格式是否有效
func IsEmailValid(email string) bool {
	// 定义电子邮件地址的正则表达式
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

// 检查密码格式是否合法
func IsValidPasswordFormat(password string) bool {
	// 密码长度必须大于等于8
	if len(password) < 8 {
		return false
	}

	// 密码必须包含数字、小写字母、大写字母、特殊字符中的至少两种
	digitRegExp := `[0-9]`
	lowercaseRegExp := `[a-z]`
	uppercaseRegExp := `[A-Z]`
	specialCharacterRegExp := `[!@#\$%\^\&*\(\)_\+{}:"\|;',\.\?\/\-\=\\]`
	var matchCount int
	if match, _ := regexp.MatchString(digitRegExp, password); match {
		matchCount++
	}
	if match, _ := regexp.MatchString(lowercaseRegExp, password); match {
		matchCount++
	}
	if match, _ := regexp.MatchString(uppercaseRegExp, password); match {
		matchCount++
	}
	if match, _ := regexp.MatchString(specialCharacterRegExp, password); match {
		matchCount++
	}
	if matchCount < 2 {
		return false
	}

	return true
}
