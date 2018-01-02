package time

import (
	"fmt"
	"time"
)

// base second
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

func TimeSince(then time.Time) string {
	now := time.Now()

	if then.After(now) {
		return "未来"
	}

	return computeTimeDiff(now.Unix() - then.Unix())
}

func computeTimeDiff(diff int64) string {
	switch {
	case diff <= 0:
		return "现在"
	case diff < 1*Minute:
		return fmt.Sprintf("%d秒前", diff)
	case diff < 1*Hour:
		return fmt.Sprintf("%d分钟前", diff/Minute)
	case diff < 1*Day:
		return fmt.Sprintf("%d小时前", diff/Hour)
	case diff < 1*Month:
		return fmt.Sprintf("%d天前", diff/Day)
	case diff < 1*Year:
		return fmt.Sprintf("%d月前", diff/Month)
	default:
		return fmt.Sprintf("%d年前", diff/Year)
	}
}
