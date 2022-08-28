package util

import "time"

func CurrentMs() int64 {
	return time.Now().UnixNano() / 1000000
}
