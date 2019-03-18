package astoverlay

import (
	"time"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
)

// SectionName merges two asterisk section names.
func SectionName(from, to *astconf.SectionName) {
	if *from != "" {
		*to = *from
	}
}

// AstYesNoNone merges two asterisk YesNoNone values.
func AstYesNoNone(from, to *astval.YesNoNone) {
	if (*from).Specified() {
		*to = *from
	}
}

// Int merges two int values.
func Int(from, to *int) {
	if *from != 0 {
		*to = *from
	}
}

// AstInt merges two asterisk Int values.
func AstInt(from, to *astval.Int) {
	if (*from).Specified() {
		*to = *from
	}
}

// String merges two string values.
func String(from, to *string) {
	if *from != "" {
		*to = *from
	}
}

// AstString merges two asterisk String values.
func AstString(from, to *astval.String) {
	if (*from).Specified() {
		*to = *from
	}
}

// StringSlice merges two string slices.
func StringSlice(from, to *[]string) {
	if *from != nil {
		*to = *from
	}
}

// Duration merges two duration values.
func Duration(from, to *time.Duration) {
	if *from != 0 {
		*to = *from
	}
}

// AstSeconds merges two asterisk Seconds values.
func AstSeconds(from, to *astval.Seconds) {
	if (*from).Specified() {
		*to = *from
	}
}
