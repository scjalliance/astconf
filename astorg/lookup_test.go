package astorg_test

import (
	"testing"

	"github.com/scjalliance/astconf/astorg"
)

func TestDataSetLookup(t *testing.T) {
	data := astorg.DataSet{
		Locations: astorg.LocationList{{Name: "Quarry"}, {Name: ""}},
		People: astorg.PersonList{
			{
				Username:       "fred",
				Extension:      "100",
				Phones:         []string{"mac1"},
				EmailAddresses: []astorg.Email{{Address: "fred@example.com", Primary: true}, {Address: ""}},
			},
			{Username: "noext"},
		},
		PhoneRoles: astorg.PhoneRoleList{
			{Username: "lobby", Extension: "200", Phones: []string{"mac1", "mac2"}},
			{Username: "", Extension: ""},
		},
		Phones:     astorg.PhoneList{{MAC: "mac1"}, {MAC: "mac2"}, {MAC: "mac3"}},
		Softphones: astorg.SoftphoneList{{Username: "fred.soft"}},
	}
	lookup := data.Lookup()

	if _, ok := lookup.LocationByName["Quarry"]; !ok {
		t.Error("LocationByName missing Quarry")
	}
	if _, ok := lookup.LocationByName[""]; ok {
		t.Error("LocationByName contains an entry for the empty name")
	}

	if p, ok := lookup.PersonByEmail["fred@example.com"]; !ok || p.Username != "fred" {
		t.Errorf("PersonByEmail[fred@example.com] = %+v, %v", p, ok)
	}
	if _, ok := lookup.PersonByEmail[""]; ok {
		t.Error("PersonByEmail contains an entry for the empty address")
	}

	if p, ok := lookup.PersonByNumber["100"]; !ok || p.Username != "fred" {
		t.Errorf("PersonByNumber[100] = %+v, %v", p, ok)
	}
	if _, ok := lookup.PersonByNumber[""]; ok {
		t.Error("PersonByNumber contains an entry for the empty extension")
	}

	if r, ok := lookup.RoleByUsername["lobby"]; !ok || r.Extension != "200" {
		t.Errorf("RoleByUsername[lobby] = %+v, %v", r, ok)
	}
	if r, ok := lookup.RoleByNumber["200"]; !ok || r.Username != "lobby" {
		t.Errorf("RoleByNumber[200] = %+v, %v", r, ok)
	}

	// Phone assignments follow the person > role > unassigned priority.
	// Softphone assignments are not part of the phone assignment map.
	tests := []struct {
		mac      string
		username string
	}{
		{mac: "mac1", username: "fred"},
		{mac: "mac2", username: "lobby"},
		{mac: "mac3", username: "mac3"},
	}
	for _, tt := range tests {
		if got := lookup.PhoneAssignments[tt.mac].Username; got != tt.username {
			t.Errorf("PhoneAssignments[%s] = %q, want %q", tt.mac, got, tt.username)
		}
	}
	if _, ok := lookup.PhoneAssignments["fred.soft"]; ok {
		t.Error("PhoneAssignments contains a softphone")
	}
}

func TestDataSetSize(t *testing.T) {
	data := astorg.DataSet{
		Locations: astorg.LocationList{{}, {}},
		People:    astorg.PersonList{{}},
		Phones:    astorg.PhoneList{{}, {}, {}},
	}
	if got := data.Size(); got != 6 {
		t.Errorf("Size() = %d, want 6", got)
	}
}
