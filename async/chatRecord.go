package async

import (
	"sync"
	"time"

	"stock-web-be/dao/db"
	"stock-web-be/gocommon/consts"
	"stock-web-be/gocommon/tlog"

	"github.com/stockAlpha/gopkg/common/safego"
)

type ChatRecordChanType struct {
	Record []db.ChatRecord
}
type AsyncToChatRecord struct {
	mu      sync.Mutex
	Records []db.ChatRecord
}

// ChatRecordChan 异步接收统一用chan
var ChatRecordChan chan ChatRecordChanType
var asyncToChatRecord AsyncToChatRecord

// ToMysqlTime 异步落库时间间隔
var ToMysqlTime = time.Minute

func init() {
	ChatRecordChan = make(chan ChatRecordChanType)
	var chatRecord ChatRecordChanType
	var ok bool
	ticker := time.NewTicker(ToMysqlTime)
	safego.SafeGoWithWG(func() {
		for {
			select {
			case chatRecord, ok = <-ChatRecordChan:
				if !ok {
					tlog.Handler.Infof(nil, consts.SyncStop, "ChatRecord ready stop")
					asyncToChatRecord.mu.Lock()
					_ = db.InsertRecord(asyncToChatRecord.Records)
					asyncToChatRecord.Records = []db.ChatRecord{}
					asyncToChatRecord.mu.Unlock()
					// 可能会有panic  ticker还没启动的话会有问题，不过会被safego兜住，同时落库也已经完成了
					ticker.Stop()
					return
				}
				var insertChatRecord AsyncToChatRecord
				asyncToChatRecord.mu.Lock()
				if len(chatRecord.Record) > 0 {
					asyncToChatRecord.Records = append(asyncToChatRecord.Records, chatRecord.Record...)
				}

				// 如果记录已经大于100条，直接准备同步落库
				if len(chatRecord.Record) >= 100 {
					insertChatRecord = asyncToChatRecord
					asyncToChatRecord.Records = []db.ChatRecord{}
				}
				asyncToChatRecord.mu.Unlock()
				if len(insertChatRecord.Records) != 0 {
					_ = db.InsertRecord(insertChatRecord.Records)
				}
			case <-ticker.C:
				var insertChatRecord AsyncToChatRecord
				// 加个锁防止有记录miss
				asyncToChatRecord.mu.Lock()
				insertChatRecord = asyncToChatRecord
				asyncToChatRecord.Records = []db.ChatRecord{}
				asyncToChatRecord.mu.Unlock()
				if len(insertChatRecord.Records) != 0 {
					_ = db.InsertRecord(insertChatRecord.Records)
				}
			}

		}
	})
}
