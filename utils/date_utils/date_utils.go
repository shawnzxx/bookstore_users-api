package date_utils

import "time"

const isoFormat = "2006-01-02T15:04:05Z"
const dateTimeFormat = "2006-01-02 15:04:05"

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowISOString() string {
	return GetNow().Format(isoFormat)
}

func GetNowDateTimeString() string {
	return GetNow().Format(dateTimeFormat)
}
