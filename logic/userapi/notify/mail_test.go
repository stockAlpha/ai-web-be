package notify

import (
	"testing"

	"stock-web-be/gocommon/conf"
)

func TestOutlook(t *testing.T) {
	conf.Init("../../../conf/app.prod.toml")
	Init()
	err := sendEmails(MailTypeOutLook, "260721735@qq.com", "testOutLook", "testOutLook")
	if err != nil {
		t.Fatal(err)
	}
}
func TestNetease(t *testing.T) {
	conf.Init("../../../conf/app.prod.toml")
	Init()
	err := sendEmails(MailType163, "260721735@qq.com", "testOutLook", "testOutLook")
	if err != nil {
		t.Fatal(err)
	}
}
func TestGmail(t *testing.T) {
	conf.Init("../../../conf/app.prod.toml")
	Init()
	err := sendEmails(MailTypeGmail, "stalary@163.com", "testOutLook", "testOutLook")
	if err != nil {
		t.Fatal(err)
	}
}
func TestMail(t *testing.T) {
	conf.Init("../../../conf/app.prod.toml")
	Init()
	err := SendEmail("260721735@qq.com", "testOutLook", "testOutLook")
	if err != nil {
		t.Fatal(err)
	}
}
