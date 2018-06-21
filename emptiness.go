package astconf

import "strings"

var emptiness = []byte(strings.Repeat(" ", 256))

func makeIndent(length int) []byte {
	if length < len(emptiness) {
		return emptiness[0:length]
	}
	return []byte(strings.Repeat(" ", length))
}
