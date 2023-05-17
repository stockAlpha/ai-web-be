package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GetRandomAvatar() string {
	rand.Seed(time.Now().UTC().UnixNano())
	randomNumber := rand.Intn(12) + 1
	return fmt.Sprintf("https://chatalpha.oss-cn-beijing.aliyuncs.com/avatar/%d.jpeg", randomNumber)
}
