package contactbuilder

import "strings"

func sanitize(name, sep string) string {
	return strings.Join(strings.Fields(strings.ToLower(name)), sep)
}
