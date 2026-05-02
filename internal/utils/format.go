package utils

import "time"

// FormatDateOrEmpty devuelve la fecha en YYYY-MM-DD o cadena vacía si es nil
func FormatDateOrEmpty(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// FormatTime devuelve el tiempo en formato RFC3339
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
