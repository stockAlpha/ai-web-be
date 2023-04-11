package userapi

import (
	"math/rand"
	"stock-web-be/dao/db"
	"strconv"
	"time"
)

func GetUserInfoByEmail(email string) (*db.UserProfile, error) {
	user := &db.UserProfile{}
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

func AddUserIntegral(userId uint64) (db.UserIntegral, error) {
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

func AddUser(email string, hashPassword string, tenantId uint64) (uint64, error) {
	//生成随机nickName
	// todo: nickname使用租户名称前缀
	nickName := "chat-" + strconv.Itoa(rand.Intn(10000))

	user := &db.UserProfile{
		NickName:    nickName,
		Email:       email,
		Password:    hashPassword,
		TenantId:    tenantId,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	err := user.InsertUser()
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
