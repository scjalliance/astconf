package astgen_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
)

// place runs Place on data and fails the test on an error.
func place(t *testing.T, data *astorg.DataSet) (astgen.Placement, []astgen.Problem) {
	t.Helper()
	p, problems, err := astgen.Place(data)
	if err != nil {
		t.Fatalf("Place() error: %v", err)
	}
	return p, problems
}

func TestPlaceHubCount(t *testing.T) {
	t.Run("no hub", func(t *testing.T) {
		data := multiServerData()
		data.Servers[2].Hub = false
		if _, _, err := astgen.Place(data); !errors.Is(err, astgen.ErrNoHub) {
			t.Errorf("Place() error = %v, want ErrNoHub", err)
		}
	})
	t.Run("two hubs", func(t *testing.T) {
		data := multiServerData()
		data.Servers[0].Hub = true
		if _, _, err := astgen.Place(data); !errors.Is(err, astgen.ErrMultipleHubs) {
			t.Errorf("Place() error = %v, want ErrMultipleHubs", err)
		}
	})
	t.Run("no servers", func(t *testing.T) {
		if _, _, err := astgen.Place(testData()); !errors.Is(err, astgen.ErrNoHub) {
			t.Errorf("Place() error = %v, want ErrNoHub", err)
		}
	})
}

func TestPlaceHosts(t *testing.T) {
	p, problems := place(t, multiServerData())
	if len(problems) != 0 {
		t.Fatalf("Place() problems: %v", problems)
	}
	if p.Hub != "hub" {
		t.Errorf("Hub = %q, want hub", p.Hub)
	}
	wantHosts := map[string]string{
		"fred":         "qry",
		"lobby":        "den",
		"000fd3000003": "qry",
		"barney":       "den",
		"fred.soft":    "hub", // A softphone is hosted by the hub even when its location has a server
		"wilma.soft":   "hub",
		"yard.soft":    "hub",
	}
	if !reflect.DeepEqual(p.Hosts, wantHosts) {
		t.Errorf("Hosts = %v, want %v", p.Hosts, wantHosts)
	}
	wantPhones := map[string]string{
		"000fd3000001": "qry",
		"000fd3000002": "den",
		"000fd3000003": "qry",
		"000fd3000004": "den",
	}
	if !reflect.DeepEqual(p.Phones, wantPhones) {
		t.Errorf("Phones = %v, want %v", p.Phones, wantPhones)
	}
	if got, want := p.HostedDevices("qry"), []string{"000fd3000003", "fred"}; !reflect.DeepEqual(got, want) {
		t.Errorf("HostedDevices(qry) = %v, want %v", got, want)
	}
}

func TestPlaceHomes(t *testing.T) {
	p, _ := place(t, multiServerData())
	wantHomes := map[string]string{
		"100": "qry",
		"101": "den",
		"102": "hub", // Remote has no server
		"200": "den",
		"201": "hub",
	}
	if !reflect.DeepEqual(p.Homes, wantHomes) {
		t.Errorf("Homes = %v, want %v", p.Homes, wantHomes)
	}
	wantOwners := map[string]string{
		"100": "fred",
		"101": "barney",
		"102": "wilma",
		"200": "lobby",
		"201": "yard",
	}
	if !reflect.DeepEqual(p.Owners, wantOwners) {
		t.Errorf("Owners = %v, want %v", p.Owners, wantOwners)
	}
	if got, want := p.HomedExtensions("hub"), []string{"102", "201"}; !reflect.DeepEqual(got, want) {
		t.Errorf("HomedExtensions(hub) = %v, want %v", got, want)
	}
}

func TestPlaceRemoteDevices(t *testing.T) {
	p, _ := place(t, multiServerData())
	wantDevices := map[string][]string{
		"100": {"fred", "fred.soft"},
		"101": {"barney"},
		"102": {"wilma.soft"},
		"200": {"lobby"},
		"201": {"yard.soft"},
	}
	if !reflect.DeepEqual(p.Devices, wantDevices) {
		t.Errorf("Devices = %v, want %v", p.Devices, wantDevices)
	}
	wantRemote := map[string]map[string][]string{
		"100": {"hub": {"fred.soft"}},
	}
	if !reflect.DeepEqual(p.Remote, wantRemote) {
		t.Errorf("Remote = %v, want %v", p.Remote, wantRemote)
	}
}

// TestPlaceRoleAtOfficeWithoutServer checks that a role at an office with no
// server lives on the hub with its softphone, so it has no remote devices.
func TestPlaceRoleAtOfficeWithoutServer(t *testing.T) {
	p, _ := place(t, multiServerData())
	if got := p.Homes["201"]; got != "hub" {
		t.Errorf("Homes[201] = %q, want hub", got)
	}
	if got, want := p.Devices["201"], []string{"yard.soft"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Devices[201] = %v, want %v", got, want)
	}
	if got := p.Remote["201"]; got != nil {
		t.Errorf("Remote[201] = %v, want none", got)
	}
}

// TestPlaceSortsByID checks that servers, peers and trunk pairs are ordered
// by server ID. The fixture's names sort Hub, Quarry, Rock Den, which is not
// the ID order den, hub, qry.
func TestPlaceSortsByID(t *testing.T) {
	p, _ := place(t, multiServerData())
	if got, want := p.ServerIDs(), []string{"den", "hub", "qry"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ServerIDs() = %v, want %v", got, want)
	}
	var peers []string
	for _, peer := range p.Peers("qry") {
		peers = append(peers, peer.ID)
	}
	if want := []string{"den", "hub"}; !reflect.DeepEqual(peers, want) {
		t.Errorf("Peers(qry) = %v, want %v", peers, want)
	}
	wantPairs := []astgen.TrunkPair{
		{ID: "den~hub", A: "den", B: "hub"},
		{ID: "den~qry", A: "den", B: "qry"},
		{ID: "hub~qry", A: "hub", B: "qry"},
	}
	if !reflect.DeepEqual(p.TrunkPairs, wantPairs) {
		t.Errorf("TrunkPairs = %v, want %v", p.TrunkPairs, wantPairs)
	}
	if got := astgen.PairID("qry", "den"); got != "den~qry" {
		t.Errorf("PairID(qry, den) = %q, want den~qry", got)
	}
}

func TestPlaceCredentials(t *testing.T) {
	p, _ := place(t, multiServerData())
	want := []astgen.Credential{
		{Kind: "trunk", ID: "den~hub", Section: "auth-trunk-den"},
		{Kind: "trunk", ID: "hub~qry", Section: "auth-trunk-qry"},
	}
	if got := p.Credentials("hub"); !reflect.DeepEqual(got, want) {
		t.Errorf("Credentials(hub) = %v, want %v", got, want)
	}
}

func TestPlaceProblems(t *testing.T) {
	tests := []struct {
		name   string
		modify func(data *astorg.DataSet)
		want   []astgen.Problem
		check  func(t *testing.T, p astgen.Placement)
	}{
		{
			name: "desk phone at an office with no server",
			modify: func(data *astorg.DataSet) {
				data.Phones[2].Location = "Remote"
			},
			want: []astgen.Problem{
				{Rule: astgen.RuleLocationNoServer, Kind: "phone", ID: "000fd3000003", Detail: `location "Remote" has no server`},
			},
			check: func(t *testing.T, p astgen.Placement) {
				if _, ok := p.Hosts["000fd3000003"]; ok {
					t.Error("the phone is still hosted")
				}
				if _, ok := p.Phones["000fd3000003"]; ok {
					t.Error("the phone is still placed")
				}
			},
		},
		{
			name: "location names a missing server",
			modify: func(data *astorg.DataSet) {
				data.Locations[1].Server = "gone"
			},
			want: []astgen.Problem{
				{Rule: astgen.RuleServerMissing, Kind: "phone", ID: "000fd3000002", Detail: `location "Den" names server "gone", which is not in the data set`},
				{Rule: astgen.RuleServerMissing, Kind: "phone", ID: "000fd3000004", Detail: `location "Den" names server "gone", which is not in the data set`},
				{Rule: astgen.RuleServerMissing, Kind: "person", ID: "barney", Detail: `location "Den" names server "gone", which is not in the data set`},
				{Rule: astgen.RuleServerMissing, Kind: "role", ID: "lobby", Detail: `location "Den" names server "gone", which is not in the data set`},
			},
			check: func(t *testing.T, p astgen.Placement) {
				if _, ok := p.Homes["101"]; ok {
					t.Error("extension 101 is still placed")
				}
			},
		},
		{
			name: "two people claim one extension",
			modify: func(data *astorg.DataSet) {
				data.People = append(data.People, astorg.Person{Username: "dino", FullName: "Dino", Extension: "100", Location: "Den"})
			},
			want: []astgen.Problem{
				{Rule: astgen.RuleExtensionDuplicate, Kind: "person", ID: "dino", Detail: `extension 100 is already held by "fred"`},
			},
			check: func(t *testing.T, p astgen.Placement) {
				if p.Owners["100"] != "fred" || p.Homes["100"] != "qry" {
					t.Errorf("extension 100 = %q on %q, want fred on qry", p.Owners["100"], p.Homes["100"])
				}
			},
		},
		{
			name: "desk phones sharing a line on two servers",
			modify: func(data *astorg.DataSet) {
				data.Phones = append(data.Phones, astorg.Phone{MAC: "000fd3000005", Location: "Den", Secret: "fred-den-secret"})
				data.People[0].Phones = append(data.People[0].Phones, "000fd3000005")
			},
			want: []astgen.Problem{
				{Rule: astgen.RuleDeviceConflict, Kind: "phone", ID: "000fd3000005", Detail: `line "fred" is already hosted on server "qry", not "den"`},
			},
			check: func(t *testing.T, p astgen.Placement) {
				if p.Hosts["fred"] != "qry" {
					t.Errorf("Hosts[fred] = %q, want qry", p.Hosts["fred"])
				}
				if _, ok := p.Phones["000fd3000005"]; ok {
					t.Error("the conflicting phone is still placed")
				}
			},
		},
		{
			name: "softphone named like a desk phone line",
			modify: func(data *astorg.DataSet) {
				data.Softphones = append(data.Softphones, astorg.Softphone{Username: "fred", Location: "Remote", Secret: "fred-clash-secret"})
			},
			want: []astgen.Problem{
				{Rule: astgen.RuleDeviceConflict, Kind: "softphone", ID: "fred", Detail: `device "fred" is already hosted on server "qry"`},
			},
			check: func(t *testing.T, p astgen.Placement) {
				if p.Hosts["fred"] != "qry" {
					t.Errorf("Hosts[fred] = %q, want qry", p.Hosts["fred"])
				}
			},
		},
		{
			name: "two softphones with one username",
			modify: func(data *astorg.DataSet) {
				data.Softphones = append(data.Softphones, astorg.Softphone{Username: "wilma.soft", Location: "Quarry", Secret: "wilma-soft-2-secret"})
			},
			want: []astgen.Problem{
				{Rule: astgen.RuleDeviceConflict, Kind: "softphone", ID: "wilma.soft", Detail: `device "wilma.soft" is already hosted on server "hub"`},
			},
			check: func(t *testing.T, p astgen.Placement) {
				if p.Hosts["wilma.soft"] != "hub" {
					t.Errorf("Hosts[wilma.soft] = %q, want hub", p.Hosts["wilma.soft"])
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := multiServerData()
			tt.modify(data)
			p, problems := place(t, data)
			if !reflect.DeepEqual(problems, tt.want) {
				t.Errorf("Place() problems =\n%v\nwant:\n%v", problems, tt.want)
			}
			tt.check(t, p)
		})
	}
}

// TestPlaceRoleClaimsPersonExtension checks that people come first in data
// set order, so a role that claims an extension a person holds is reported
// and left out.
func TestPlaceRoleClaimsPersonExtension(t *testing.T) {
	data := multiServerData()
	data.PhoneRoles = append(data.PhoneRoles, astorg.PhoneRole{Username: "gate", DisplayName: "Gate", Extension: "100", Location: "Den"})
	p, problems := place(t, data)
	want := []astgen.Problem{
		{Rule: astgen.RuleExtensionDuplicate, Kind: "role", ID: "gate", Detail: `extension 100 is already held by "fred"`},
	}
	if !reflect.DeepEqual(problems, want) {
		t.Errorf("Place() problems =\n%v\nwant:\n%v", problems, want)
	}
	if p.Owners["100"] != "fred" || p.Homes["100"] != "qry" {
		t.Errorf("extension 100 = %q on %q, want fred on qry", p.Owners["100"], p.Homes["100"])
	}
}
