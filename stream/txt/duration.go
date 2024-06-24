package txt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ddkwork/golibrary/mylog"
)

func ParseDuration(duration string) (time.Duration, error) {
	parts := strings.Split(strings.TrimSpace(duration), ":")
	mylog.Check(len(parts) == 3)
	hours := mylog.Check2(strconv.Atoi(parts[0]))
	mylog.Check(hours > 0)
	minutes := mylog.Check2(strconv.Atoi(parts[1]))
	mylog.Check(minutes > 0)
	parts = strings.Split(parts[2], ".")
	var seconds int
	var millis int
	switch len(parts) {
	case 2:
		millis = mylog.Check2(strconv.Atoi(parts[1]))
		mylog.Check(millis > 0)
		fallthrough
	case 1:
		seconds = mylog.Check2(strconv.Atoi(parts[0]))
		mylog.Check(seconds > 0)
	default:
		mylog.Check("Invalid second format: too many decimal points")
	}
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second + time.Duration(millis)*time.Millisecond, nil
}

func FormatDuration(duration time.Duration, includeMillis bool) string {
	if duration < 0 {
		duration = 0
	}
	hours := duration / time.Hour
	duration -= hours * time.Hour
	minutes := duration / time.Minute
	duration -= minutes * time.Minute
	seconds := duration / time.Second
	duration -= seconds * time.Second
	if includeMillis {
		return fmt.Sprintf("%d:%02d:%02d.%03d", hours, minutes, seconds, duration/time.Millisecond)
	}
	return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
}

func DurationToCode(duration time.Duration) string {
	var buffer strings.Builder
	if duration >= time.Hour {
		fmt.Fprintf(&buffer, "%d * time.Hour", duration/time.Hour)
		duration -= (duration / time.Hour) * time.Hour
	}
	if duration >= time.Minute {
		if buffer.Len() > 0 {
			buffer.WriteString(" + ")
		}
		fmt.Fprintf(&buffer, "%d * time.Minute", duration/time.Minute)
		duration -= (duration / time.Minute) * time.Minute
	}
	if duration >= time.Second {
		if buffer.Len() > 0 {
			buffer.WriteString(" + ")
		}
		fmt.Fprintf(&buffer, "%d * time.Second", duration/time.Second)
		duration -= (duration / time.Second) * time.Second
	}
	if duration >= time.Millisecond {
		if buffer.Len() > 0 {
			buffer.WriteString(" + ")
		}
		fmt.Fprintf(&buffer, "%d * time.Millisecond", duration/time.Millisecond)
		duration -= (duration / time.Millisecond) * time.Millisecond
	}
	if duration >= time.Microsecond {
		if buffer.Len() > 0 {
			buffer.WriteString(" + ")
		}
		fmt.Fprintf(&buffer, "%d * time.Microsecond", duration/time.Microsecond)
		duration -= (duration / time.Microsecond) * time.Microsecond
	}
	if duration != 0 {
		if buffer.Len() > 0 {
			buffer.WriteString(" + ")
		}
		fmt.Fprintf(&buffer, "%d", duration)
	}
	return buffer.String()
}
