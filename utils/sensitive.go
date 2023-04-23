package utils

import "github.com/zeromicro/go-zero/core/stringx"

func ReplaceSensitiveWord(content string, replaceMap map[string]string) string {
	replacer := stringx.NewReplacer(replaceMap)
	return replacer.Replace(content)
}
