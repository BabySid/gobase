package gobase

import (
	"fmt"
	"time"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	TimeFormat     = "15:04:05"
	DateFormat     = "2006-01-02"
)

func FormatTime() string {
	return FormatTimeStampWithFormat(time.Now().Unix(), TimeFormat)
}

func FormatDate() string {
	return FormatTimeStampWithFormat(time.Now().Unix(), DateFormat)
}

func FormatDateTime() string {
	return FormatTimeStampWithFormat(time.Now().Unix(), DateTimeFormat)
}

func FormatTimeStamp(ts int64) string {
	return FormatTimeStampWithFormat(ts, DateTimeFormat)
}

func FormatTimeStampWithFormat(ts int64, format string) string {
	return time.Unix(ts, 0).Format(format)
}

func FormatTimeStampMilliWithFormat(ts int64, format string) string {
	return time.UnixMilli(ts).Format(format)
}

func ParseTimestamp(timestamp string) (time.Time, error) {
	// 尝试多种 RFC3339 变体格式
	formats := []string{
		time.RFC3339Nano,                // 2025-06-07T14:15:32.307+08:00
		time.RFC3339,                    // 2025-06-07T14:15:32+08:00
		"2006-01-02T15:04:05.000Z07:00", // 其他毫秒格式
		"2006-01-02T15:04:05Z07:00",     // 无毫秒格式
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestamp); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间戳: %s", timestamp)
}
