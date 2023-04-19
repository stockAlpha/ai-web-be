package utils

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func Md5(str string) (string, error) {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil)), nil
}
