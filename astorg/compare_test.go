package astorg_test

import (
	"testing"

	"github.com/scjalliance/astconf/astorg"
)

func TestLocationEqual(t *testing.T) {
	a := astorg.Location{Name: "Quarry", Abbreviation: "QRY", PagingGroups: []string{"all"}}
	b := a
	if !a.Equal(&b) {
		t.Error("identical locations are not Equal")
	}
	b.PagingGroups = []string{"all", "quarry"}
	if a.Equal(&b) {
		t.Error("locations with different paging groups are Equal")
	}
	b = a
	b.AreaCode = "360"
	if a.Equal(&b) {
		t.Error("locations with different area codes are Equal")
	}
	var nilLoc *astorg.Location
	if !nilLoc.Equal(nil) {
		t.Error("two nil locations are not Equal")
	}
	if a.Equal(nil) || nilLoc.Equal(&a) {
		t.Error("nil and non-nil locations are Equal")
	}
}

func TestPersonEqual(t *testing.T) {
	a := astorg.Person{
		Username:       "fred",
		Extension:      "100",
		Phones:         []string{"mac1"},
		ContactNumbers: []astorg.Number{{Label: "Mobile", Dial: "3605551234"}},
		EmailAddresses: []astorg.Email{{Address: "fred@example.com", Primary: true}},
	}
	b := a
	if !a.Equal(&b) {
		t.Error("identical people are not Equal")
	}
	b.ContactNumbers = []astorg.Number{{Label: "Mobile", Dial: "3605554321"}}
	if a.Equal(&b) {
		t.Error("people with different contact numbers are Equal")
	}
	b = a
	b.EmailAddresses = []astorg.Email{{Address: "fred@example.com", Primary: false}}
	if a.Equal(&b) {
		t.Error("people with different email flags are Equal")
	}
	b = a
	b.Hidden = true
	if a.Equal(&b) {
		t.Error("people with different hidden flags are Equal")
	}
}

func TestPhoneRoleEqual(t *testing.T) {
	a := astorg.PhoneRole{Username: "lobby", Extension: "200", Phones: []string{"mac1"}}
	b := a
	if !a.Equal(&b) {
		t.Error("identical roles are not Equal")
	}
	b.Phones = []string{"mac2"}
	if a.Equal(&b) {
		t.Error("roles with different phones are Equal")
	}
}

func TestPhoneEqual(t *testing.T) {
	a := astorg.Phone{MAC: "mac1", Location: "Quarry", Secret: "s", Templates: []string{"t"}}
	b := a
	if !a.Equal(&b) {
		t.Error("identical phones are not Equal")
	}
	b.Templates = nil
	if a.Equal(&b) {
		t.Error("phones with different templates are Equal")
	}
}

func TestSoftphoneEqual(t *testing.T) {
	a := astorg.Softphone{Username: "fred.soft", Location: "Quarry", Secret: "s"}
	b := a
	if !a.Equal(&b) {
		t.Error("identical softphones are not Equal")
	}
	b.Secret = "x"
	if a.Equal(&b) {
		t.Error("softphones with different secrets are Equal")
	}
}
