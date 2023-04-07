package notify

import (
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

func SendEmail(to, subject, body string) error {
	fmt.Println(SendMail + "," + SendPassword + "," + SmtpServer)
	auth := smtp.PlainAuth("", SendMail, SendPassword, SmtpServer)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	err := smtp.SendMail(SmtpServer+":"+SmtpPort, auth, SendMail, []string{to}, msg)
	return err
}
