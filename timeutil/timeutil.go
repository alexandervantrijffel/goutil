package timeutil

import (
	"time"

	"github.com/alexandervantrijffel/goutil/errorcheck"
)

func GetLastDayOfMonth(t time.Time) time.Time {
	currentYear, currentMonth, _ := t.Date()
	currentLocation := t.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return lastOfMonth
}

func TimeToShortDateTimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func TimeFromShortDateTimeString(t string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", t)
}

func Ago(t time.Time) time.Duration {
	return time.Now().UTC().Sub(t.UTC())
}

// ParseRfc3339Time: Parse time in the format 2006-01-02T15:04:05.999999999Z
func ParseRfc3339Time(timeSz string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, timeSz)
	if err != nil {
		_ = errorcheck.CheckLogf(err, "Failed to parse time %s", timeSz)
	}
	return t
}
