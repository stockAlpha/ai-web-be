package userapi

import (
	"math/rand"
	"strconv"
	"time"

	"stock-web-be/dao/db"
	"stock-web-be/utils"

	"gorm.io/gorm"
)

var (
	EMail = 1
	Phone = 2
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

func GetUserById(id uint64) (*db.User, error) {
	user := &db.User{}
	err := user.GetUserById(id)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return user, nil
}

func GetUserByInviteCode(inviteCode string, transaction *gorm.DB) (*db.User, error) {
	user := &db.User{}
	err := user.GetUserByInviteCode(inviteCode, transaction)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return user, nil
}

func GetUserIntegralByUserId(userId uint64, transaction *gorm.DB) (*db.UserIntegral, error) {
	u := &db.UserIntegral{}
	err := u.GetUserIntegralByUserId(userId, transaction)
	if err != nil {
		return nil, err
	}
	if u.ID == 0 {
		return nil, nil
	}
	return u, nil
}

func CreateUserIntegral(userId uint64, amount int, transaction *gorm.DB) (db.UserIntegral, error) {
	integral := db.UserIntegral{
		UserId:      userId,
		Amount:      amount,
		TotalAmount: amount,
		UpdateTime:  time.Now(),
		CreateTime:  time.Now(),
	}
	err := integral.InsertUserIntegral(transaction)
	if err != nil {
		return integral, err
	}
	return integral, nil
}

func AddUserIntegral(userId uint64, amount int, transaction *gorm.DB) error {
	integral, err := GetUserIntegralByUserId(userId, transaction)
	if err != nil {
		return err
	}
	err = integral.AddAmount(amount, transaction)
	if err != nil {
		return err
	}
	return nil
}

func SubUserIntegral(userId uint64, amount int) error {
	integral, err := GetUserIntegralByUserId(userId, nil)
	if err != nil {
		return err
	}
	err = integral.SubAmount(amount)
	if err != nil {
		return err
	}
	return nil
}

func AddUser(email, hashPassword string, transaction *gorm.DB) (uint64, error) {
	// 生成随机nickName
	nickName := "chat-" + strconv.Itoa(rand.Intn(10000))
	// 生成邀请码
	inviteCode := utils.GenerateCode()
	// 生成头像
	avatar := utils.GetRandomAvatar()

	user := &db.User{
		NickName:   nickName,
		Email:      email,
		Password:   hashPassword,
		InviteCode: inviteCode,
		Avatar:     avatar,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	err := user.InsertUser(transaction)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func UpdateUser(userId uint64, nickName, avatar string) {
	user := &db.User{
		ID:       userId,
		NickName: nickName,
		Avatar:   avatar,
	}
	user.UpdateUser()
}

func UpdateUserPassword(userId uint64, password string, transaction *gorm.DB) error {
	user := &db.User{
		ID:       userId,
		Password: password,
	}
	return user.UpdateUserPassword(transaction)
}

func AddInviteRelation(fromUserId uint64, toUserId uint64, inviteCode string, transaction *gorm.DB) error {
	relation := &db.InviteRelation{
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		InviteCode: inviteCode,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	err := relation.InsertRelation(transaction)
	if err != nil {
		return err
	}
	return nil
}

func AddFeedback(fromUserId uint64, feedbackType int, content string) error {
	feedback := &db.Feedback{
		FromUserId:   fromUserId,
		FeedbackType: feedbackType,
		Content:      content,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}

	err := feedback.InsertFeedback()
	if err != nil {
		return err
	}
	return nil
}

func GetUserByAuthType(subjectId string, subjectType int) (*db.User, error) {
	var user *db.User
	var err error
	if subjectType == EMail {
		user, err = GetUserByEmail(subjectId)
	}

	if err != nil {
		return nil, err
	}
	return user, err
}
