package astorg_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/scjalliance/astconf/astorg"
)

func TestAssignmentMapAdd(t *testing.T) {
	m := make(astorg.AssignmentMap)

	if m.Add("", astorg.Assignment{Type: astorg.PersonAssigned, Username: "fred"}) {
		t.Error("Add() with empty key = true")
	}

	if !m.Add("mac1", astorg.Assignment{Type: astorg.Unassigned, Username: "mac1"}) {
		t.Error("first Add() = false")
	}

	// A role outranks unassigned
	if !m.Add("mac1", astorg.Assignment{Type: astorg.RoleAssigned, Username: "lobby"}) {
		t.Error("Add() of role over unassigned = false")
	}

	// A second role does not replace the first
	if m.Add("mac1", astorg.Assignment{Type: astorg.RoleAssigned, Username: "kitchen"}) {
		t.Error("Add() of second role = true")
	}
	if m["mac1"].Username != "lobby" {
		t.Errorf("mac1 assigned to %q, want lobby", m["mac1"].Username)
	}

	// A person outranks a role
	if !m.Add("mac1", astorg.Assignment{Type: astorg.PersonAssigned, Username: "fred"}) {
		t.Error("Add() of person over role = false")
	}

	// Nothing outranks a person
	if m.Add("mac1", astorg.Assignment{Type: astorg.RoleAssigned, Username: "lobby"}) {
		t.Error("Add() of role over person = true")
	}
	if m["mac1"].Username != "fred" {
		t.Errorf("mac1 assigned to %q, want fred", m["mac1"].Username)
	}
}

func TestAssignmentMapByUsername(t *testing.T) {
	m := astorg.AssignmentMap{
		"mac1": {Type: astorg.PersonAssigned, Username: "fred"},
		"mac2": {Type: astorg.PersonAssigned, Username: "fred"},
		"mac3": {Type: astorg.RoleAssigned, Username: "lobby"},
		"mac4": {Type: astorg.Unassigned},
	}
	got := m.ByUsername()
	for _, keys := range got {
		sort.Strings(keys)
	}
	want := map[string][]string{
		"fred":  {"mac1", "mac2"},
		"lobby": {"mac3"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ByUsername() = %v, want %v", got, want)
	}
}

func TestMergeAssignments(t *testing.T) {
	phones := astorg.PhoneList{{MAC: "mac1"}, {MAC: "mac2"}, {MAC: "mac3"}}.Assignments()
	roles := astorg.PhoneRoleList{{Username: "lobby", Phones: []string{"mac1", "mac2"}}}.Assignments()
	people := astorg.PersonList{{Username: "fred", Phones: []string{"mac1"}}}.Assignments()

	// Order of the maps does not matter; priority does
	got := astorg.MergeAssignments(people, phones, roles)
	want := astorg.AssignmentMap{
		"mac1": {Type: astorg.PersonAssigned, Username: "fred"},
		"mac2": {Type: astorg.RoleAssigned, Username: "lobby"},
		"mac3": {Type: astorg.Unassigned, Username: "mac3"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeAssignments() = %v, want %v", got, want)
	}
}
