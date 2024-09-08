package utils

import "regexp"

func ValidateRange(value int, min int, max int) bool {
	return value >= min && value <= max
}
func ValidateRange64(value int64, min int64, max int64) bool {
	return value >= min && value <= max
}

func ValidateLength(value string, min int, max int) bool {
	return len(value) >= min && len(value) <= max
}

func ValidateRegex(value string, pattern string) bool {
	return regexp.MustCompile(pattern).MatchString(value)
}

func ValidateName(value string) bool {
	return ValidateRegex(value, `^[a-zA-Z0-9_]+$`)
}
func ValidateSHA256(value string) bool {
	return ValidateRegex(value, `^[a-f0-9]{64}$`)
}
func ValidateTime(value string) bool {
	return ValidateRegex(value, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)
}
func ValidateFileName(value string) bool {
	return ValidateRegex(value, `^[^<>:;,?"*|/\\]+$`)
}
