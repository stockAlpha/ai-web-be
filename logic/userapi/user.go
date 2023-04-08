package userapi

import (
	"math/rand"
	"stock-web-be/dao/db"
	"strconv"
	"time"
)

func GetUserInfoByEmail(email string) (*db.User, error) {
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

func AddUser(email string, hashPassword string, tenantId uint64) (uint64, error) {
	//生成随机nickName
	// todo: nickname使用租户名称前缀
	nickName := "chat-" + strconv.Itoa(rand.Intn(10000))

	user := &db.User{
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
