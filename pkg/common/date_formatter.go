package common

import "time"

// FormatDateTime formats a time.Time object into a string in the format "DD-Mon-YYYY HH:mm:ss"
func FormatDateTime(t time.Time) string {
	return t.Format("02-Jan-2006 15:04:05")
}

// FormatDate formats time to DD-Mon-YYYY
func FormatDate(t time.Time) string {
	return t.Format("02-Jan-2006")
}

// FormatTime formats time to HH:mm:ss
func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}
