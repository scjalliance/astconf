package astgen

import "strings"

// callerIDName prepares a display name for use as the quoted name part of
// an asterisk callerid value. Double quotes would end the quoted part early.
func callerIDName(name string) string {
	return strings.ReplaceAll(name, "\"", "")
}
