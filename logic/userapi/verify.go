package userapi

import (
	"stock-web-be/dao/db"
	"time"
)

type SendCodeType int

var (
	MailCode = 1
)

func InsertEmailVerificationCode(code string, email string) error {
	//todo 失效其余验证码
	//计算expireTime
	expireTime := time.Now().Add(time.Minute * 5)

	verificationCode := &db.VerificationCode{
		SendSubjectName: email,
		SendSubjectType: MailCode,
		Code:            code,
		ExpireTime:      expireTime,
		CreatedTime:     time.Now(),
		UpdatedTime:     time.Now(),
	}

	err := verificationCode.InsertCode()
	if err != nil {
		return err
	}
	return nil
}

func ExistCode(code string, email string) (bool, error) {
	verificationCode := &db.VerificationCode{}

	verificationCodes, err := verificationCode.GetSendCodeByEmailAndCode(email, code)
	if err != nil {
		return false, err
	}

	return len(verificationCodes) > 0, nil
}
