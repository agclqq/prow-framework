package times

import (
	"fmt"
	"time"
)

var MaxTime = time.Date(9999, 12, 31, 23, 59, 59, 0, time.Local)
var MaxTimeUnix = MaxTime.Unix()

const FormatDatetime = "2006-01-02 15:04:05"
const FormatDatetimeMicro = "2006-01-02 15:04:05.000"
const FormatDate = "2006-01-02"
const FormatTime = "15:04:05"
const FormatShortTime = "15:04"

func StringToTime(layout, str string) time.Time {
	location, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}
	}
	//inLocation, err := time.Parse(layout, str)
	inLocation, err := time.ParseInLocation(layout, str, location)
	if err != nil {
		return time.Time{}
	}
	return inLocation
}

func GetYesterdayDate() time.Time {
	t := time.Now()
	// 计算昨天的时间
	yesterday := t.AddDate(0, 0, -1)
	// 格式化时间，得到昨天的日期
	yesterdayStr := yesterday.Format("2006-01-02")
	return StringToTime(FormatDate, yesterdayStr)
}

func TimeToString(layout string, t time.Time) string {
	return t.Format(layout)
}

// Timestamp 获取当前时间戳: 20230227172847
func Timestamp() string {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	timestamp := fmt.Sprintf("%d%02d%02d%02d%02d%02d", year, month, day, hour, minute, second)
	return timestamp
}
