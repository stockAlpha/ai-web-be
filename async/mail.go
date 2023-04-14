package async

import (
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"
	"stock-web-be/logic/userapi/notify"

	"github.com/stockAlpha/gopkg/common/safego"
)

type MailChanType struct {
	To, Subject, Body string
}

// MailChan 异步发送mail统一用chan
var MailChan chan MailChanType

func init() {
	MailChan = make(chan MailChanType)
	var mail MailChanType
	var ok bool
	safego.SafeGoWithWG(func() {
		for {
			mail, ok = <-MailChan
			if !ok {
				tlog.Handler.Infof(nil, consts.SyncStop, "mail send ready stop")
				break
			}
			notify.SendEmail(mail.To, mail.Subject, mail.Body)
		}
	})
}
