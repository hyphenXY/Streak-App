package utils

import (
	"time"
)

// StartOfWeek returns the start of the week (Monday, 00:00:00) for the given time.
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	// Go's Weekday: Sunday=0, Monday=1, ..., Saturday=6
	// We want Monday as start of week
	if weekday == 0 {
		weekday = 7
	}
	start := t.AddDate(0, 0, -weekday+1)
	return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfWeek returns the end of the week (Sunday, 23:59:59) for the given time.
func EndOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	// Go's Weekday: Sunday=0, Monday=1, ..., Saturday=6
	// We want Sunday as end of week
	offset := 7 - weekday
	end := t.AddDate(0, 0, offset)
	return time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}