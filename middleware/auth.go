package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"stock-web-be/gocommon/conf"
	"stock-web-be/gocommon/consts"
	"stock-web-be/logic/userapi"
)

var JWTSign string

func Init() {
	JWTSign = conf.Handler.GetString("jwt.sign")
}

func ValidUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		notAuthApis := consts.NotAuthApisMap
		// 如果请求路径为 不需要鉴权的api
		if notAuthApis[path] != "" {
			c.Next() // 跳过鉴权，继续处理请求
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Not authorized",
			})
			return
		}
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Invalid token format",
			})
			return
		}
		tokenString := authHeader[7:]

		//判断token是否在黑名单中
		exist, err := userapi.IsOnBlackList(tokenString)
		if err != nil {
			//不稳定的情况下先不判断，避免用户重复登录
		}
		if exist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token is in blacklist",
			})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWTSign), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  err.Error(),
			})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := claims["userId"].(string)
			email := claims["email"].(string)
			exp := claims["exp"].(float64)
			c.Set("user_id", userId)
			c.Set("email", email)
			c.Set("token", tokenString)
			c.Set("exp", exp)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Invalid token",
			})
		}
	}
}
