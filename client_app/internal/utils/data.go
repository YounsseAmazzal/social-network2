package utils

import (
	"fmt"
	"strings"
	"time"
)

func RegleThedate(str string) string {
	str = strings.TrimSuffix(strings.Replace(str, "T", " ", 1), "Z")
	return TimeNewFormat(str)
}

// TimeSinceCreation conveeert the date 
func TimeNewFormat(creationDate string) string {
	const layout = "2006-01-02 15:04:05"
	t, err := time.Parse(layout, creationDate)
	if err != nil {
		return "Just now"
	}

	elapsed := time.Since(t)
	switch {
	case elapsed.Seconds() < 60:
		return formatTime(int(elapsed.Seconds()), "second")
	case elapsed.Minutes() < 60:
		return formatTime(int(elapsed.Minutes()), "minute")
	case elapsed.Hours() < 24:
		return formatTime(int(elapsed.Hours()), "hour")
	case elapsed.Hours() < 24*30:
		return formatTime(int(elapsed.Hours()/24), "day")
	case elapsed.Hours() < 24*30*12:
		return formatTime(int(elapsed.Hours()/(24*30)), "month")
	default:
		return formatTime(int(elapsed.Hours()/(24*365)), "year")
	}
}

// formatTime heelpeer function to colect the date 
func formatTime(value int, unit string) string {
	if value == 1 {
		return fmt.Sprintf("%d %s ago", value, unit)
	}
	return fmt.Sprintf("%d %ss ago", value, unit)
}
