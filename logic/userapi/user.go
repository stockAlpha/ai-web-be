package userapi

import (
	"fmt"
	"math/rand"
	"stock-web-be/dao/db"
	"stock-web-be/utils"
	"strconv"
	"time"
)

func GetUserByEmail(email string) (*db.User, error) {
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

func GetUserByInviteCode(inviteCode string) (*db.User, error) {
	user := &db.User{}
	err := user.GetUserByInviteCode(inviteCode)
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

func CreateUserIntegral(userId uint64, amount int) (db.UserIntegral, error) {
	integral := db.UserIntegral{
		UserId: userId,
		// 初始化10积分
		Amount:     amount,
		UpdateTime: time.Now(),
		CreateTime: time.Now(),
	}
	err := integral.InsertUserIntegral()
	if err != nil {
		return integral, err
	}
	return integral, nil
}

func AddUserIntegral(userId uint64, amount int) error {
	integral, err := GetUserIntegralByUserId(userId)
	if err != nil {
		return err
	}
	err = integral.AddAmount(amount)
	if err != nil {
		return err
	}
	return nil
}

func SubUserIntegral(userId uint64, amount int) error {
	integral, err := GetUserIntegralByUserId(userId)
	if err != nil {
		return err
	}
	err = integral.SubAmount(amount)
	if err != nil {
		return err
	}
	return nil
}

func AddUser(email string, hashPassword string) (uint64, error) {
	// 生成随机nickName
	nickName := "chat-" + strconv.Itoa(rand.Intn(10000))
	// 生成邀请码
	inviteCode := utils.GenerateCode()
	// 生成头像
	randomNumber := rand.Intn(12) + 1
	avatar := fmt.Sprintf("/public/avatar/%d.jpeg", randomNumber)

	user := &db.User{
		NickName:   nickName,
		Email:      email,
		Password:   hashPassword,
		InviteCode: inviteCode,
		Avatar:     avatar,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	err := user.InsertUser()
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func AddInviteRelation(fromUserId uint64, toUserId uint64, inviteCode string) error {
	relation := &db.InviteRelation{
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		InviteCode: inviteCode,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	err := relation.InsertRelation()
	if err != nil {
		return err
	}
	return nil
}
