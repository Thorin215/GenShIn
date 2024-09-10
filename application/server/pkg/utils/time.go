package utils

import "time"

// Return the current time in RFC3339 format
func GetTimeString() string {
	return time.Now().Format(time.RFC3339)
}
