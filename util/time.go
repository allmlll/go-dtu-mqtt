package util

import "time"

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func NowWithoutSecond() string {
	return time.Now().Format("2006-01-02 15:04") + ":00"
}
