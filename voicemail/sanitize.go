package voicemail

import "strings"

func sanitizeOption(value string) string {
	value = strings.ReplaceAll(value, "|", "")
	value = strings.ReplaceAll(value, ",", "")
	return value
}
