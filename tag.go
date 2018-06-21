package astconf

import (
	"strings"
)

type fieldTag struct {
	name string
	omit bool
	opts tagOptions
}

func parseTag(s string) (tag fieldTag) {
	if s == "" {
		return
	}
	if s == "-" {
		tag.omit = true
		return
	}
	if idx := strings.Index(s, ","); idx != -1 {
		tag.name = s[:idx]
		tag.opts = tagOptions(strings.Split(s[idx+1:], ","))
	} else {
		tag.name = s
	}
	return
}

func (ft *fieldTag) Contains(option string) bool {
	return ft.opts.Contains(option)
}

type tagOptions []string

func (o tagOptions) Contains(option string) bool {
	if o == nil {
		return false
	}
	for _, s := range o {
		if s == option {
			return true
		}
	}
	return false
}

/*
func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}
*/
