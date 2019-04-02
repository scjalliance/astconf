package astorg

import "github.com/scjalliance/astconf/astorg/astorgvm"

// Person represents a person in an organization.
type Person struct {
	//Server             string // Identifies which server manages presence
	Extension          string // Primary Extension
	Username           string
	FirstName          string
	LastName           string
	FullName           string
	Organization       string
	Title              string
	Location           string
	Hidden             bool // Hidden from organization contact groups
	CallerID           string
	AreaCode           string
	Firmware           string
	CalendarURL        string
	VoicemailExtension string
	VoicemailCode      string
	VoicemailAccess    astorgvm.Access
	Alerts             []string
	Ringtones          []string
	Apps               []string
	PagingGroups       []string
	ContactNumbers     []Number // Additional phone numbers
	Phones             []string // MACs of phones to be assigned this role
	EmailAddresses     []Email
}

// Equal reports whether p and q are equal.
func (p *Person) Equal(q *Person) bool {
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
	if p.FirstName != q.FirstName {
		return false
	}
	if p.LastName != q.LastName {
		return false
	}
	if p.FullName != q.FullName {
		return false
	}
	if p.Organization != q.Organization {
		return false
	}
	if p.Title != q.Title {
		return false
	}
	if p.Location != q.Location {
		return false
	}
	if p.Hidden != q.Hidden {
		return false
	}
	if p.CallerID != q.CallerID {
		return false
	}
	if p.AreaCode != q.AreaCode {
		return false
	}
	if p.Firmware != q.Firmware {
		return false
	}
	if p.CalendarURL != q.CalendarURL {
		return false
	}
	if p.VoicemailExtension != q.VoicemailExtension {
		return false
	}
	if p.VoicemailCode != q.VoicemailCode {
		return false
	}
	if p.VoicemailAccess != q.VoicemailAccess {
		return false
	}

	// Compare slices
	if !equalStrings(p.Alerts, q.Alerts) {
		return false
	}
	if !equalStrings(p.Ringtones, q.Ringtones) {
		return false
	}
	if !equalStrings(p.Apps, q.Apps) {
		return false
	}
	if !equalStrings(p.PagingGroups, q.PagingGroups) {
		return false
	}
	if !equalNumbers(p.ContactNumbers, q.ContactNumbers) {
		return false
	}
	if !equalStrings(p.Phones, q.Phones) {
		return false
	}
	if !equalEmails(p.EmailAddresses, q.EmailAddresses) {
		return false
	}

	return true
}
