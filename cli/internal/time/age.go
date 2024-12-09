package time

import (
	"fmt"
	"time"
)

// ToAge formats a time into an human-redable age string.
func ToAge(t time.Time) string {
	return toAge(t, time.Now())
}
func toAge(t, now time.Time) string {
	d := now.Sub(t)
	if sec := int(d.Seconds()); sec < 60 {
		return fmt.Sprintf("%ds", sec)
	} else if min := int(d.Minutes()); min < 60 {
		return fmt.Sprintf("%dm", min)
	} else if d.Hours() < 6 {
		return fmt.Sprintf("%.0fh%dm", d.Hours(), min%60)
	} else if d.Hours() < 24 {
		return fmt.Sprintf("%.0fh", d.Hours())
	} else if d.Hours() < 24*365 {
		return fmt.Sprintf("%.0fd", d.Hours()/24)
	}
	return fmt.Sprintf("%.0fy", d.Hours()/(24*365))
}
