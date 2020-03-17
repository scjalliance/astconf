package astorg

// PhoneRole is a role that is assigned to a particular phone.
type PhoneRole struct {
	Extension      string
	Username       string
	DisplayName    string
	Location       string
	AreaCode       string
	Firmware       string
	Hidden         bool // Hidden from organization contact groups
	ContactsSource string
	CallerID       string
	MailboxNumber  string
	Tags           []string
	PagingGroups   []string
	Apps           []string
	Phones         []string // MACs of phones to be assigned this role
	Softphones     []string // Usernames of software phones to be assigned this role
}

// Equal reports whether p and q are equal.
func (p *PhoneRole) Equal(q *PhoneRole) bool {
	// Compare nil-ness
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}

	// Compare simple fields
	if p.Extension != q.Extension {
		return false
	}
	if p.Username != q.Username {
		return false
	}
	if p.DisplayName != q.DisplayName {
		return false
	}
	if p.Location != q.Location {
		return false
	}
	if p.AreaCode != q.AreaCode {
		return false
	}
	if p.Firmware != q.Firmware {
		return false
	}
	if p.Hidden != q.Hidden {
		return false
	}
	if p.ContactsSource != q.ContactsSource {
		return false
	}
	if p.CallerID != q.CallerID {
		return false
	}
	if p.MailboxNumber != q.MailboxNumber {
		return false
	}

	// Compare slices
	if !equalStrings(p.Tags, q.Tags) {
		return false
	}
	if !equalStrings(p.PagingGroups, q.PagingGroups) {
		return false
	}
	if !equalStrings(p.Apps, q.Apps) {
		return false
	}
	if !equalStrings(p.Phones, q.Phones) {
		return false
	}
	if !equalStrings(p.Softphones, q.Softphones) {
		return false
	}

	return true
}
