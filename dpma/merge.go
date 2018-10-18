package dpma

import (
	"time"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
)

func mergeSectionName(from, to *astconf.SectionName) {
	if *from != "" {
		*to = *from
	}
}

func mergeAstYesNoNone(from, to *astval.YesNoNone) {
	if (*from).Specified() {
		*to = *from
	}
}

func mergeInt(from, to *int) {
	if *from != 0 {
		*to = *from
	}
}

func mergeAstInt(from, to *astval.Int) {
	if (*from).Specified() {
		*to = *from
	}
}

func mergeString(from, to *string) {
	if *from != "" {
		*to = *from
	}
}

func mergeAstString(from, to *astval.String) {
	if (*from).Specified() {
		*to = *from
	}
}

func mergeStringSlice(from, to *[]string) {
	if *from != nil {
		*to = *from
	}
}

func mergeDuration(from, to *time.Duration) {
	if *from != 0 {
		*to = *from
	}
}

func mergeAstSeconds(from, to *astval.Seconds) {
	if (*from).Specified() {
		*to = *from
	}
}

func mergeLogoFileSlice(from, to *[]LogoFile) {
	if *from != nil {
		*to = *from
	}
}
