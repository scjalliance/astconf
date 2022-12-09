package astorg

// Phone represents a phone within an organization.
type Phone struct {
	MAC       string
	Location  string
	Secret    string
	Templates []string
}

// Equal reports whether p and q are equal.
func (p *Phone) Equal(q *Phone) bool {
	// Compare nil-ness
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}

	// Compare simple fields
	if p.MAC != q.MAC {
		return false
	}
	if p.Location != q.Location {
		return false
	}
	if p.Secret != q.Secret {
		return false
	}

	// Compare slices
	if !equalStrings(p.Templates, q.Templates) {
		return false
	}

	return true
}
