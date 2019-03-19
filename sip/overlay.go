package sip

func overlayType(from, to *Type) {
	if (*from).Specified() {
		*to = *from
	}
}
