package async

import (
	"github.com/robfig/cron/v3"
	"log"
	"stock-web-be/controller/public"
)

var cronjob *cron.Cron

func init() {
	cronjob = cron.New()
}
func StartCron() {
	//五分钟一次
	cronjob.AddFunc("*/5 * * * *", func() {
		log.Println("cronjob start", public.CheckAll())
	})
	cronjob.Start()
}
