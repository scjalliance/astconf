package dpma

func mergeLogoFileSlice(from, to *[]LogoFile) {
	if *from != nil {
		*to = *from
	}
}
