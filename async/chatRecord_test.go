package async

import (
	"testing"
	"time"
)

func TestChatRecord(t *testing.T) {
	cType := ChatRecordChanType{}
	ChatRecordChan <- cType
	time.Sleep(time.Second * 3)
	ChatRecordChan <- cType
	time.Sleep(time.Second * 3)
	close(ChatRecordChan)
	time.Sleep(time.Second * 3)
}
