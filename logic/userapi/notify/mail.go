package notify

import (
	"errors"
	"fmt"
	"net/smtp"
	"stock-web-be/gocommon/conf"
)

var (
	SendMail     string
	SendPassword string
	SmtpServer   string
	SmtpPort     string
)

func Init() {
	// CIFlow
	SendMail = conf.Handler.GetString("mail.from")
	SendPassword = conf.Handler.GetString("mail.password")
	SmtpServer = conf.Handler.GetString("mail.smtpServer")
	SmtpPort = conf.Handler.GetString("mail.smtpPort")
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
	fmt.Println(SendMail + "," + SendPassword + "," + SmtpServer)
	// outlook need use custom
	auth := LoginAuth(SendMail, SendPassword)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	err := smtp.SendMail(SmtpServer+":"+SmtpPort, auth, SendMail, []string{to}, msg)
	return err
}
