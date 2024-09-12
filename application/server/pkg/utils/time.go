package utils

import "time"

// Return the current time in ISO8601 format
func GetTimeString() string {
	return time.Now(). /*.UTC()*/ Format("2006-01-02T15:04:05Z")
}
