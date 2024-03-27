package stream

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ddkwork/golibrary/mylog"
)

const TimeLayout = "2006-01-02 15:04:05"

func FormatTime(t time.Time) string { return t.Format(TimeLayout) }
func UnFormatTime(s string) time.Time {
	parse, err := time.Parse(TimeLayout, s)
	if !mylog.Error(err) {
		return time.Time{}
	}
	return parse
}

func FormatDuration(d time.Duration) string { return d.String() }
func UnFormatDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if !mylog.Error(err) {
		return 0
	}
	return duration
}
func GetTimeNowString() string { return time.Now().Format("2006-01-02 15:04:05 ") }

func GetTimeStamp13Bits() int64 { return time.Now().UnixNano() / 1000000 }

func GetTimeStamp() string { return strconv.FormatInt(time.Now().UnixNano()/1000000, 10) }

func GetDiffDays(dstTime string) string {
	a, _ := time.Parse("2006-01-02", dstTime)
	now := a.Sub(time.Now())
	days := int(now.Hours() / 24)
	years := days / 365
	months := (days % 365) / 30
	remainingDays := (days % 365) % 30
	hours := int(now.Hours()) % 24
	minutes := int(now.Minutes()) % 60
	seconds := int(now.Seconds()) % 60

	s := New("")
	s.WriteStringLn(fmt.Sprintf("相差天数 %d 天", days))
	s.WriteStringLn(fmt.Sprintf("相差年数 %d 年", years))
	s.WriteStringLn(fmt.Sprintf("相差月数 %d 月", months))
	s.WriteStringLn(fmt.Sprintf("相差时数 %d 时", hours))
	s.WriteStringLn(fmt.Sprintf("相差分数 %d 分", minutes))
	s.WriteStringLn(fmt.Sprintf("相差秒数 %d 秒", seconds))
	s.WriteStringLn(fmt.Sprintf("相差时间 %d 年 %d 月 %d 天 %d 时 %d 分 %d 秒",
		years, months, remainingDays, hours, minutes, seconds))
	return s.String()
}
