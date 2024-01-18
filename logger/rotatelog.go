package logger

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func RotateDailyLog(filePath string, days uint) (*rotatelogs.RotateLogs, error) {
	writer, err := rotatelogs.New(
		filePath+".%Y%m%d",
		//rotatelogs.WithLinkName(filePath),           // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(days)*time.Hour*24), // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24),               // 日志切割时间间隔
		//rotatelogs.ForceNewFile(),
	)
	return writer, err
}
