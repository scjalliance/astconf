package astorg

// equalStrings returns true of two slices of strings have the same
// set of values.
func equalStrings(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// equalNumbers returns true of two slices of numbers have the same
// set of values.
func equalNumbers(a []Number, b []Number) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// equalEmails returns true of two slices of emails have the same
// set of values.
func equalEmails(a []Email, b []Email) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
