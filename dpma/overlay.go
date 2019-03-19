package dpma

func overlayLogoFileSlice(from, to *[]LogoFile) {
	if *from != nil {
		*to = *from
	}
}
