package astorg_test

import (
	"testing"

	"github.com/scjalliance/astconf/astorg"
)

func TestServerEqual(t *testing.T) {
	base := astorg.Server{ID: "quarry", Name: "Quarry", Hub: false, Address: "quarry-phone.tailnet"}

	same := base
	if !base.Equal(&same) {
		t.Error("identical servers reported unequal")
	}

	tests := []struct {
		name   string
		change func(*astorg.Server)
	}{
		{name: "id", change: func(s *astorg.Server) { s.ID = "bedrock" }},
		{name: "name", change: func(s *astorg.Server) { s.Name = "Bedrock" }},
		{name: "hub", change: func(s *astorg.Server) { s.Hub = true }},
		{name: "address", change: func(s *astorg.Server) { s.Address = "other.tailnet" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			other := base
			tt.change(&other)
			if base.Equal(&other) {
				t.Errorf("servers differing in %s reported equal", tt.name)
			}
		})
	}

	var nilServer *astorg.Server
	if !nilServer.Equal(nil) {
		t.Error("two nil servers reported unequal")
	}
	if base.Equal(nil) {
		t.Error("a server reported equal to nil")
	}
}

func TestServerMapAdd(t *testing.T) {
	m := make(astorg.ServerMap)

	if m.Add(astorg.Server{Name: "No ID"}) {
		t.Error("Add() with empty id = true")
	}

	if !m.Add(astorg.Server{ID: "quarry", Name: "Quarry"}) {
		t.Error("first Add() = false")
	}

	// A second server with the same id does not replace the first
	if m.Add(astorg.Server{ID: "quarry", Name: "Bedrock"}) {
		t.Error("duplicate Add() = true")
	}
	if m["quarry"].Name != "Quarry" {
		t.Errorf("quarry is %q, want Quarry", m["quarry"].Name)
	}
}

func TestServerListByID(t *testing.T) {
	servers := astorg.ServerList{
		{ID: "quarry", Name: "Quarry"},
		{ID: "hub", Hub: true},
		{ID: ""},
		{ID: "quarry", Name: "Bedrock"},
	}
	byID := servers.ByID()

	if len(byID) != 2 {
		t.Errorf("ByID() has %d entries, want 2", len(byID))
	}
	if s, ok := byID["hub"]; !ok || !s.Hub {
		t.Errorf("ByID()[hub] = %+v, %v", s, ok)
	}
	if s := byID["quarry"]; s.Name != "Quarry" {
		t.Errorf("ByID()[quarry] = %+v, want the first quarry", s)
	}
	if _, ok := byID[""]; ok {
		t.Error("ByID() contains an entry for the empty id")
	}
}

func TestDataSetCountsAndComparesServers(t *testing.T) {
	a := astorg.DataSet{Servers: astorg.ServerList{{ID: "quarry"}, {ID: "hub", Hub: true}}}
	if got := a.Size(); got != 2 {
		t.Errorf("Size() = %d, want 2", got)
	}

	b := astorg.DataSet{Servers: astorg.ServerList{{ID: "quarry"}, {ID: "hub", Hub: true}}}
	if !a.Equal(&b) {
		t.Error("data sets with identical servers reported unequal")
	}

	c := astorg.DataSet{Servers: astorg.ServerList{{ID: "quarry"}, {ID: "hub"}}}
	if a.Equal(&c) {
		t.Error("data sets differing in a server reported equal")
	}

	d := astorg.DataSet{Servers: astorg.ServerList{{ID: "quarry"}}}
	if a.Equal(&d) {
		t.Error("data sets with different server counts reported equal")
	}
}

func TestLookupLocationServer(t *testing.T) {
	data := astorg.DataSet{
		Servers: astorg.ServerList{
			{ID: "quarry", Name: "Quarry", Address: "quarry-phone.tailnet"},
			{ID: "hub", Name: "Hub", Hub: true},
		},
		Locations: astorg.LocationList{
			{Name: "Quarry Office", Server: "quarry"},
			{Name: "Quarry Pit", Server: "quarry"},
			{Name: "Nowhere", Server: "missing"},
			{Name: "Unassigned"},
		},
	}
	lookup := data.Lookup()

	if s, ok := lookup.ServerByID["hub"]; !ok || !s.Hub {
		t.Errorf("ServerByID[hub] = %+v, %v", s, ok)
	}

	tests := []struct {
		location string
		wantID   string
		wantOK   bool
	}{
		{location: "Quarry Office", wantID: "quarry", wantOK: true},
		{location: "Quarry Pit", wantID: "quarry", wantOK: true},
		{location: "Nowhere", wantID: "missing", wantOK: false}, // names a server that does not exist
		{location: "Unassigned", wantOK: false},                 // names no server at all
		{location: "Not a place", wantOK: false},                // not a known location
		{location: "", wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.location, func(t *testing.T) {
			id, s, ok := lookup.LocationServer(tt.location)
			if id != tt.wantID {
				t.Errorf("LocationServer(%q) id = %q, want %q", tt.location, id, tt.wantID)
			}
			if ok != tt.wantOK {
				t.Fatalf("LocationServer(%q) ok = %v, want %v", tt.location, ok, tt.wantOK)
			}
			if ok && s.ID != tt.wantID {
				t.Errorf("LocationServer(%q) server = %q, want %q", tt.location, s.ID, tt.wantID)
			}
		})
	}
}

func TestLookupHub(t *testing.T) {
	tests := []struct {
		name    string
		servers astorg.ServerList
		wantID  string
		wantOK  bool
	}{
		{
			name:    "one hub",
			servers: astorg.ServerList{{ID: "quarry"}, {ID: "hub", Hub: true}},
			wantID:  "hub",
			wantOK:  true,
		},
		{
			name:    "no hub",
			servers: astorg.ServerList{{ID: "quarry"}, {ID: "bedrock"}},
		},
		{
			name:    "two hubs",
			servers: astorg.ServerList{{ID: "quarry", Hub: true}, {ID: "hub", Hub: true}},
		},
		{
			name: "no servers",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := astorg.DataSet{Servers: tt.servers}
			s, ok := data.Lookup().Hub()
			if ok != tt.wantOK {
				t.Fatalf("Hub() ok = %v, want %v", ok, tt.wantOK)
			}
			if s.ID != tt.wantID {
				t.Errorf("Hub() = %q, want %q", s.ID, tt.wantID)
			}
		})
	}
}
