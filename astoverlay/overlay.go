package astoverlay

import (
	"time"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
)

// SectionName overlays two asterisk section names.
//
// If from is non-empty, the value of to will be replaced with from.
func SectionName(from, to *astconf.SectionName) {
	if *from != "" {
		*to = *from
	}
}

// AstYesNoNone overlays two asterisk YesNoNone values.
//
// If from holds a specified value, the value of to will be replaced with from.
func AstYesNoNone(from, to *astval.YesNoNone) {
	if (*from).Specified() {
		*to = *from
	}
}

// Int overlays two int values.
//
// If from is non-empty, the value of to will be replaced with from.
func Int(from, to *int) {
	if *from != 0 {
		*to = *from
	}
}

// AstInt overlays two asterisk Int values.
//
// If from holds a specified value, the value of to will be replaced with from.
func AstInt(from, to *astval.Int) {
	if (*from).Specified() {
		*to = *from
	}
}

// String overlays two string values.
//
// If from is non-empty, the value of to will be replaced with from.
func String(from, to *string) {
	if *from != "" {
		*to = *from
	}
}

// AstString overlays two asterisk String values.
//
// If from holds a specified value, the value of to will be replaced with from.
func AstString(from, to *astval.String) {
	if (*from).Specified() {
		*to = *from
	}
}

// StringSlice overlays two string slices.
//
// If from is non-empty, the value of to will be replaced with from.
func StringSlice(from, to *[]string) {
	if *from != nil {
		*to = *from
	}
}

// Duration overlays two duration values.
//
// If from is non-empty, the value of to will be replaced with from.
func Duration(from, to *time.Duration) {
	if *from != 0 {
		*to = *from
	}
}

// AstSeconds overlays two asterisk Seconds values.
//
// If from holds a specified value, the value of to will be replaced with from.
func AstSeconds(from, to *astval.Seconds) {
	if (*from).Specified() {
		*to = *from
	}
}
