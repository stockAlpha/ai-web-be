package notify

import (
	"errors"
	"fmt"
	"net/smtp"
	"sync"
	"time"

	"stock-web-be/gocommon/conf"

	"github.com/stockAlpha/gopkg/common/safego"
)

var (
	SendMail     string
	SendPassword string
	SmtpServer   string
	SmtpPort     string
)

type MailClient struct {
	SendMail     string
	SendPassword string
	SendAuth     string
	SmtpServer   string
	SmtpPort     string
	AppendMsg    string //新浪要求
	Life         int    //标记失败次数
}

// var Mails sync.Map
var Mails MailMap

type MailMap struct {
	mu       sync.Mutex
	Maps     map[MailType]MailClient
	FailMaps map[MailType]MailClient
}
type MailType string

const (
	MailTypeOutLook MailType = "outlook"
	MailType163     MailType = "163"
	MailTypeGmail   MailType = "gmail"
	MailTypeQQ      MailType = "qq"
	MailTypeSOHU    MailType = "sohu"
	MailTypeSINA    MailType = "sina"
)

func Init() {
	Mails.Maps = make(map[MailType]MailClient)
	Mails.FailMaps = make(map[MailType]MailClient)
	// CIFlow
	//outlook
	Mails.Maps[MailTypeOutLook] = MailClient{
		SendMail:     conf.Handler.GetString("mail_outlook.from"),
		SendPassword: conf.Handler.GetString("mail_outlook.password"),
		SmtpServer:   conf.Handler.GetString("mail_outlook.smtpServer"),
		SmtpPort:     conf.Handler.GetString("mail_outlook.smtpPort"),
		AppendMsg:    "",
		Life:         3,
	}
	Mails.Maps[MailType163] = MailClient{
		SendMail:     conf.Handler.GetString("mail_163.from"),
		SendPassword: conf.Handler.GetString("mail_163.password"),
		SmtpServer:   conf.Handler.GetString("mail_163.smtpServer"),
		SmtpPort:     conf.Handler.GetString("mail_163.smtpPort"),
		AppendMsg:    "",
		Life:         3,
	}
	Mails.Maps[MailTypeGmail] = MailClient{
		SendMail:     conf.Handler.GetString("mail_gmail.from"),
		SendPassword: conf.Handler.GetString("mail_gmail.password"),
		SmtpServer:   conf.Handler.GetString("mail_gmail.smtpServer"),
		SmtpPort:     conf.Handler.GetString("mail_gmail.smtpPort"),
		AppendMsg:    "",
		Life:         3,
	}
	// todo 等qq授权码
	// todo 退信了
	// todo 疑似官方bug

	RetryFailMail()
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}

func SendEmail(to, subject, body string) error {
	//随机选择一个mail
	// todo 现在只拿第一个
	if len(Mails.Maps) <= 0 {
		return errors.New("not enough alive mail")
	}
	var mailType MailType
	var err error
	// 只重试一次
	for i := 0; i < 2; i++ {
		for k := range Mails.Maps {
			mailType = k
			break
		}
		err = sendEmails(mailType, to, subject, body)
		if err == nil {
			return nil
		}
	}
	return err
}
func sendEmails(mailType MailType, to, subject, body string) error {
	mail := Mails.Maps[mailType]
	// outlook need use custom
	auth := LoginAuth(mail.SendMail, mail.SendPassword)
	msgStr := "To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n"
	if mail.AppendMsg != "" {
		msgStr = msgStr + mail.AppendMsg + " \r\n"
	}
	msg := []byte(msgStr)
	err := smtp.SendMail(mail.SmtpServer+":"+mail.SmtpPort, auth, mail.SendMail, []string{to}, msg)
	fmt.Println("send email success to: ", to, " subject: ", subject, " body: ", body, " message: ", msgStr)
	if err != nil {
		Mails.mu.Lock()
		mail.Life = 0
		Mails.FailMaps[mailType] = mail
		delete(Mails.Maps, mailType)
		Mails.mu.Unlock()
	}
	return err
}
func RetryFailMail() {
	// 纯后台任务 无需wg和管控
	safego.SafeGo(func() {
		for {
			// todo 暂定5min刷新一次
			time.Sleep(time.Minute * 5)
			for k, v := range Mails.FailMaps {
				err := sendEmails(k, "260721735@qq.com", "retry mail", "retry mail")
				if err == nil {
					Mails.mu.Lock()
					Mails.Maps[k] = v
					delete(Mails.FailMaps, k)
					Mails.mu.Unlock()
				}
				time.Sleep(time.Second * 5)
			}
		}
	})
}
