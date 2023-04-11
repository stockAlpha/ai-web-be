package userapi

import (
	"math/rand"
	"stock-web-be/dao/db"
	"strconv"
	"time"
)

func GetUserProfileByEmail(email string) (*db.User, error) {
	user := &db.User{}
	err := user.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return user, nil
}

func GetUserIntegralByUserId(userId uint64) (*db.UserIntegral, error) {
	u := &db.UserIntegral{}
	err := u.GetUserIntegralByUserId(userId)
	if err != nil {
		return nil, err
	}
	if u.ID == 0 {
		return nil, nil
	}
	return u, nil
}

func CreateUserIntegral(userId uint64) (db.UserIntegral, error) {
	integral := db.UserIntegral{
		UserId: userId,
		// 初始化10积分
		Amount:      10,
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}
	err := integral.InsertUserIntegral()
	if err != nil {
		return integral, err
	}
	return integral, nil
}

func AddUserIntegral(userId uint64, amount int) (db.UserIntegral, error) {
	integral := db.UserIntegral{}
	err := integral.AddAmount(userId, amount)
	if err != nil {
		return integral, err
	}
	return integral, nil
}

func SubUserIntegral(userId uint64, amount int) (db.UserIntegral, error) {
	integral := db.UserIntegral{}
	err := integral.SubAmount(userId, amount)
	if err != nil {
		return integral, err
	}
	return integral, nil
}

func AddUser(email string, hashPassword string) (uint64, error) {
	//生成随机nickName
	nickName := "chat-" + strconv.Itoa(rand.Intn(10000))

	user := &db.User{
		NickName:    nickName,
		Email:       email,
		Password:    hashPassword,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	err := user.InsertUser()
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
