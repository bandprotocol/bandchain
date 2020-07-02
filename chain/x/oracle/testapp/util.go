package testapp

import (
	"time"
)

// ParseTime is a helper function to parse from number to time.Time with UTC locale.
func ParseTime(t int64) time.Time {
	return time.Unix(t, 0).UTC()
}
