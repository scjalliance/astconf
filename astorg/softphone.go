package astorg

// Softphone represents a software phone within an organization.
type Softphone struct {
	Username  string
	Location  string
	Secret    string
	Templates []string
}

// Equal reports whether p and q are equal.
func (p *Softphone) Equal(q *Softphone) bool {
	// Compare nil-ness
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}

	// Compare simple fields
	if p.Username != q.Username {
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
