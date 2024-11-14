package helpers

import "time"

func FormatDateTime(t time.Time) string {
	return t.Format("02-Jan-2006 15:04:05")
}
