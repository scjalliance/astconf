package astorg

// A Tag holds a set of phone settings shared by people and phone roles
// that have been assigned the tag.
type Tag struct {
	Name         string
	Firmware     string
	PagingGroups []string
	Alerts       []string
	Ringtones    []string
	Apps         []string
}

// Equal reports whether p and q are equal.
func (p *Tag) Equal(q *Tag) bool {
	// Compare nil-ness
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}

	// Compare simple fields
	if p.Name != q.Name {
		return false
	}
	if p.Firmware != q.Firmware {
		return false
	}

	// Compare slices
	if !equalStrings(p.PagingGroups, q.PagingGroups) {
		return false
	}
	if !equalStrings(p.Alerts, q.Alerts) {
		return false
	}
	if !equalStrings(p.Ringtones, q.Ringtones) {
		return false
	}
	if !equalStrings(p.Apps, q.Apps) {
		return false
	}
	return true
}
