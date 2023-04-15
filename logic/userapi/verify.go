package userapi

import (
	"fmt"
	"gorm.io/gorm"
	"stock-web-be/dao/db"
	"stock-web-be/gocommon/consts"
	"stock-web-be/logic/userapi/notify"
	"stock-web-be/utils"
	"time"
)

type SendCodeType int

var (
	MailCode               = 1
	ChangePasswordMailCode = 101
)

func SendGeneralVerificationCode(subjectName string, codeType int) error {
	//生成随机验证码(是否需要控频?)
	code := utils.GenerateCode()
	var err error
	//发送验证码
	if codeType == MailCode || codeType == ChangePasswordMailCode {
		err = notify.SendEmail(subjectName, consts.SendCodeSubject, fmt.Sprintf(consts.SendCodeContent, code))
	}

	//通用处理逻辑
	if err != nil {
		return err
	}

	//验证码存入db
	err = InsertVerificationCode(code, subjectName, codeType)
	if err != nil {
		return err
	}

	return nil
}

func InsertEmailVerificationCode(code string, email string) error {
	//todo 失效其余验证码
	//计算expireTime
	expireTime := time.Now().Add(time.Minute * 5)

	verificationCode := &db.VerificationCode{
		SendSubjectName: email,
		SendSubjectType: MailCode,
		Code:            code,
		ExpireTime:      expireTime,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}

	err := verificationCode.InsertCode()
	if err != nil {
		return err
	}
	return nil
}

func InsertVerificationCode(code string, subjectName string, codeType int) error {
	//计算expireTime
	expireTime := time.Now().Add(time.Minute * 5)

	verificationCode := &db.VerificationCode{
		SendSubjectName: subjectName,
		SendSubjectType: codeType,
		Code:            code,
		ExpireTime:      expireTime,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}

	return verificationCode.InsertCode()
}

func ExistCode(code string, email string) (bool, error) {
	verificationCode := &db.VerificationCode{}

	verificationCodes, err := verificationCode.GetSendCodeByEmailAndCode(email, code)
	if err != nil {
		return false, err
	}

	return len(verificationCodes) > 0, nil
}

func ExistCodeByCodeType(code string, subjectName string, subjectType int) (bool, error) {
	verificationCode := &db.VerificationCode{}

	verificationCodes, err := verificationCode.GetSendCodeByCodeType(subjectName, code, subjectType)
	if err != nil {
		return false, err
	}

	return len(verificationCodes) > 0, nil
}

func ExpireCode(code string, subjectName string, subjectType int, transaction *gorm.DB) error {
	verificationCode := &db.VerificationCode{
		Code:            code,
		SendSubjectName: subjectName,
		SendSubjectType: subjectType,
		ExpireTime:      time.Now(),
	}

	return verificationCode.UpdateByCode(transaction)
}

func VerificationCodeTypeToAuthType(verificationCodeType int) int {
	if verificationCodeType == MailCode ||
		verificationCodeType == ChangePasswordMailCode {
		return EMail
	}
	return EMail
}
