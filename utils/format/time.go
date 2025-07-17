package format

import "time"

func GetCurrentTimeAsTimestamp() int64 {
	return time.Now().Unix()
}
