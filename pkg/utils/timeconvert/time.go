// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package timeconvert

import (
	"time"
)

const (
	LayoutDateTime   = "2006-01-02 15:04:05"
	LayoutDateOnly   = "2006-01-02"
	LayoutTimeOnly   = "15:04:05"
	LayoutCompact    = "20060102150405"
	LayoutISO8601    = "2006-01-02T15:04:05Z07:00"
)

var (
	DefaultTimeZone = "Local"
)

// SetDefaultTimeZone sets the default timezone for the application
func SetDefaultTimeZone(tz string) {
	DefaultTimeZone = tz
	loc, err := time.LoadLocation(tz)
	if err == nil {
		time.Local = loc
	}
}

// CurrentTime returns current time in standard layout
func CurrentTime() string {
	return time.Now().Format(LayoutDateTime)
}

// CurrentDate returns current date in standard layout
func CurrentDate() string {
	return time.Now().Format(LayoutDateOnly)
}

// FormatToDateTime formats time to "YYYY-MM-DD HH:mm:ss"
func FormatToDateTime(t time.Time) string {
	return t.Format(LayoutDateTime)
}

// FormatToDate formats time to "YYYY-MM-DD"
func FormatToDate(t time.Time) string {
	return t.Format(LayoutDateOnly)
}

// ParseDateTime parses string "YYYY-MM-DD HH:mm:ss" to time
func ParseDateTime(str string) (time.Time, error) {
	return time.ParseInLocation(LayoutDateTime, str, time.Local)
}

// ParseDate parses string "YYYY-MM-DD" to time
func ParseDate(str string) (time.Time, error) {
	return time.ParseInLocation(LayoutDateOnly, str, time.Local)
}

// StringToTime attempts to parse a string in various layouts
func StringToTime(str string) (time.Time, error) {
	layouts := []string{
		LayoutDateTime,
		LayoutDateOnly,
		LayoutISO8601,
		time.RFC3339,
		LayoutCompact,
	}

	var lastErr error
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, str, time.Local)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

// CalculateDuration returns the duration between two time strings
func CalculateDuration(start, end string) (time.Duration, error) {
	startTime, err := ParseDateTime(start)
	if err != nil {
		return 0, err
	}
	endTime, err := ParseDateTime(end)
	if err != nil {
		return 0, err
	}
	return endTime.Sub(startTime), nil
}

// IsExpired checks if the target time is before now
func IsExpired(targetTime time.Time) bool {
	return targetTime.Before(time.Now())
}

// AddDays adds n days to the time
func AddDays(t time.Time, n int) time.Time {
	return t.AddDate(0, 0, n)
}
