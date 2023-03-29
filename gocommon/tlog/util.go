package tlog

import (
	"strconv"
	"time"
)

// GenLogId /**
func GenLogId() string {
	now := time.Now()
	logId := ((now.UnixNano()*1e5 + now.UnixNano()/1e7) & 0x7FFFFFFF) | 0x80000000
	return strconv.FormatInt(logId, 10)
}
