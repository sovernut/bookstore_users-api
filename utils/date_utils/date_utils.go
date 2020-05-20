package date_utils

import "time"

const (
	apiDateLayoout = "2006-01-02T15:04:05Z"
	apiDbLayout    = "2006-01-02 15:04:05"
)

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowString() string {
	now := GetNow()
	return now.Format(apiDateLayoout)
}

func GetNowDBFormat() string {
	now := GetNow()
	return now.Format(apiDbLayout)
}
