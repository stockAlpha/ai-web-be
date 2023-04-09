package userapi

import (
	"github.com/dgrijalva/jwt-go"
	"stock-web-be/gocommon/conf"
	"time"
)

var JWTSign string

func Init() {
	JWTSign = conf.Handler.GetString("jwt.sign")
}

func GenerateToken(userId, email string) (string, error) {
	// 定义载荷（payload）
	claims := jwt.MapClaims{}
	claims["userId"] = userId
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	// 使用密钥对载荷进行签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(JWTSign))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
