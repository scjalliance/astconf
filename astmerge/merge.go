package astmerge

// StringSlice merges two string slices, combining and deduplicating
// their values.
func StringSlice(from, to *[]string) {
	*to = dedupStringSlice(append(*from, (*to)...))
}
