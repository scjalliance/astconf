package astgen_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/dpma"
	"github.com/scjalliance/astconf/pjsip"
)

// spareMAC is the spare desk phone in multiServerData: no location and no
// owner.
const spareMAC = "000fd3000009"

// TestPlaceSpareClassification checks how Place sorts desk phones that
// have no server: a phone with no location and no owner is a spare, a
// phone at a location with no server is location-no-server, and a phone
// with no location that a person or role owns is owned-no-location.
func TestPlaceSpareClassification(t *testing.T) {
	data := multiServerData()
	data.Phones = append(data.Phones,
		astorg.Phone{MAC: "000fd3000008", Secret: "second-spare-secret"},     // A second spare, which sorts first
		astorg.Phone{MAC: "000fd300000a", Location: "Remote", Secret: "rmt"}, // At a location with no server
		astorg.Phone{MAC: "000fd300000b", Secret: "wilma-desk-secret"},       // No location, owned by a person
		astorg.Phone{MAC: "000fd300000c", Secret: "yard-desk-secret"},        // No location, owned by a role
	)
	data.People[2].Phones = []string{"000fd300000b"}     // wilma
	data.PhoneRoles[1].Phones = []string{"000fd300000c"} // yard

	p, problems := place(t, data)

	if want := []string{"000fd3000008", spareMAC}; !reflect.DeepEqual(p.Spares, want) {
		t.Errorf("Spares = %v, want %v", p.Spares, want)
	}
	want := []astgen.Problem{
		{Rule: astgen.RuleLocationNoServer, Kind: "phone", ID: "000fd300000a", Detail: `location "Remote" has no server`},
		{Rule: astgen.RuleOwnedNoLocation, Kind: "phone", ID: "000fd300000b", Detail: `the phone has no location but is assigned to "wilma"`},
		{Rule: astgen.RuleOwnedNoLocation, Kind: "phone", ID: "000fd300000c", Detail: `the phone has no location but is assigned to "yard"`},
	}
	if !reflect.DeepEqual(problems, want) {
		t.Errorf("Place() problems =\n%v\nwant:\n%v", problems, want)
	}
	for _, mac := range []string{"000fd3000008", spareMAC, "000fd300000a", "000fd300000b", "000fd300000c"} {
		if host, ok := p.Phones[mac]; ok {
			t.Errorf("Phones[%s] = %q, want no host", mac, host)
		}
		if host, ok := p.Hosts[mac]; ok {
			t.Errorf("Hosts[%s] = %q, want no host", mac, host)
		}
	}
	// The owners' extensions still ring only their hosted devices.
	if got := p.Devices["102"]; !reflect.DeepEqual(got, []string{"wilma.soft"}) {
		t.Errorf("Devices[102] = %v, want [wilma.soft]", got)
	}
}

// TestPlaceSpareConflicts checks that a softphone or a located desk phone
// whose device name is a spare's line is reported and left out, so the
// device never renders twice on a server.
func TestPlaceSpareConflicts(t *testing.T) {
	t.Run("softphone", func(t *testing.T) {
		data := multiServerData()
		data.Softphones = append(data.Softphones, astorg.Softphone{Username: spareMAC, Secret: "clash-secret"})
		p, problems := place(t, data)
		want := []astgen.Problem{{
			Rule:   astgen.RuleDeviceConflict,
			Kind:   "softphone",
			ID:     spareMAC,
			Detail: `device "000fd3000009" is already the line of a spare desk phone on every server`,
		}}
		if !reflect.DeepEqual(problems, want) {
			t.Errorf("Place() problems =\n%v\nwant:\n%v", problems, want)
		}
		if _, ok := p.Hosts[spareMAC]; ok {
			t.Error("the conflicting softphone is hosted")
		}
	})
	t.Run("located phone after the spare", func(t *testing.T) {
		data := multiServerData()
		data.Phones = append(data.Phones, astorg.Phone{MAC: spareMAC, Location: "Quarry", Secret: "clash-secret"})
		p, problems := place(t, data)
		want := []astgen.Problem{{
			Rule:   astgen.RuleDeviceConflict,
			Kind:   "phone",
			ID:     spareMAC,
			Detail: `line "000fd3000009" is already the line of a spare desk phone on every server`,
		}}
		if !reflect.DeepEqual(problems, want) {
			t.Errorf("Place() problems =\n%v\nwant:\n%v", problems, want)
		}
		if _, ok := p.Hosts[spareMAC]; ok {
			t.Error("the conflicting phone is hosted")
		}
	})
	t.Run("located phone before the spare", func(t *testing.T) {
		data := multiServerData()
		data.Phones = append([]astorg.Phone{{MAC: spareMAC, Location: "Quarry", Secret: "clash-secret"}}, data.Phones...)
		p, problems := place(t, data)
		want := []astgen.Problem{{
			Rule:   astgen.RuleDeviceConflict,
			Kind:   "phone",
			ID:     spareMAC,
			Detail: `line "000fd3000009" is already hosted on server "qry", so it cannot be a spare`,
		}}
		if !reflect.DeepEqual(problems, want) {
			t.Errorf("Place() problems =\n%v\nwant:\n%v", problems, want)
		}
		if len(p.Spares) != 0 {
			t.Errorf("Spares = %v, want none", p.Spares)
		}
	})
}

// TestSpareRenderedOnEveryServer checks that every server renders the
// spare's DPMA phone and line and its pjsip endpoint and aor, each exactly
// as the whole-dataset generators render the unassigned phone.
func TestSpareRenderedOnEveryServer(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)

	endpointBase := pjsip.Endpoint{Context: "from-internal"}
	aorBase := pjsip.AOR{MaxContacts: astval.NewInt(5)}
	const contactsURL = "https://contacts.example.com/"

	pick := func(t *testing.T, what string, rendered []string, names []string, name string) string {
		t.Helper()
		var found []string
		for i, n := range names {
			if n == name {
				found = append(found, rendered[i])
			}
		}
		if len(found) != 1 {
			t.Fatalf("%s has %d entries named %s, want 1", what, len(found), name)
		}
		return found[0]
	}

	var legacy struct{ phone, line, endpoint, aor string }
	{
		var r, n []string
		for _, ph := range astgen.Phones(data, dpma.Phone{}, contactsURL) {
			r, n = append(r, render(t, &ph)), append(n, ph.MAC)
		}
		legacy.phone = pick(t, "Phones()", r, n, spareMAC)
		r, n = nil, nil
		for _, l := range astgen.Lines(data) {
			r, n = append(r, render(t, &l)), append(n, l.Name)
		}
		legacy.line = pick(t, "Lines()", r, n, spareMAC)
		r, n = nil, nil
		for _, e := range astgen.Endpoints(data, endpointBase, "default") {
			r, n = append(r, render(t, &e)), append(n, e.Name)
		}
		legacy.endpoint = pick(t, "Endpoints()", r, n, spareMAC)
		r, n = nil, nil
		for _, a := range astgen.AORs(data, aorBase) {
			r, n = append(r, render(t, &a)), append(n, a.Name)
		}
		legacy.aor = pick(t, "AORs()", r, n, spareMAC+"-aor")
	}
	if !strings.Contains(legacy.phone, "line = "+spareMAC+"\n") {
		t.Fatalf("the legacy phone does not name its line by its MAC:\n%s", legacy.phone)
	}

	for _, serverID := range p.ServerIDs() {
		t.Run(serverID, func(t *testing.T) {
			var r, n []string
			for _, ph := range astgen.ServerPhones(data, p, serverID, dpma.Phone{}, contactsURL) {
				r, n = append(r, render(t, &ph)), append(n, ph.MAC)
			}
			if got := pick(t, "ServerPhones()", r, n, spareMAC); got != legacy.phone {
				t.Errorf("ServerPhones() spare =\n%s\nwant:\n%s", got, legacy.phone)
			}
			r, n = nil, nil
			for _, l := range astgen.ServerLines(data, p, serverID) {
				r, n = append(r, render(t, &l)), append(n, l.Name)
			}
			if got := pick(t, "ServerLines()", r, n, spareMAC); got != legacy.line {
				t.Errorf("ServerLines() spare =\n%s\nwant:\n%s", got, legacy.line)
			}
			r, n = nil, nil
			for _, e := range astgen.ServerEndpoints(data, p, serverID, endpointBase, "default") {
				r, n = append(r, render(t, &e)), append(n, e.Name)
			}
			if got := pick(t, "ServerEndpoints()", r, n, spareMAC); got != legacy.endpoint {
				t.Errorf("ServerEndpoints() spare =\n%s\nwant:\n%s", got, legacy.endpoint)
			}
			r, n = nil, nil
			for _, a := range astgen.ServerAORs(data, p, serverID, aorBase) {
				r, n = append(r, render(t, &a)), append(n, a.Name)
			}
			if got := pick(t, "ServerAORs()", r, n, spareMAC+"-aor"); got != legacy.aor {
				t.Errorf("ServerAORs() spare =\n%s\nwant:\n%s", got, legacy.aor)
			}

			wantCred := astgen.Credential{Kind: "phone", ID: spareMAC, Section: spareMAC + "-auth"}
			found := false
			for _, cred := range p.Credentials(serverID) {
				if cred == wantCred {
					found = true
				}
			}
			if !found {
				t.Errorf("Credentials(%s) = %v, want it to include %v", serverID, p.Credentials(serverID), wantCred)
			}
		})
	}
}

// TestSpareUnknownServer checks that a server ID not in the placement
// renders no spare and lists no credential for one.
func TestSpareUnknownServer(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	if got := astgen.ServerPhones(data, p, "gone", dpma.Phone{}, ""); len(got) != 0 {
		t.Errorf("ServerPhones(gone) = %v, want none", got)
	}
	if got := astgen.ServerEndpoints(data, p, "gone", pjsip.Endpoint{}, "default"); len(got) != 0 {
		t.Errorf("ServerEndpoints(gone) = %v, want none", got)
	}
	if got := p.Credentials("gone"); len(got) != 0 {
		t.Errorf("Credentials(gone) = %v, want none", got)
	}
}

// TestSpareNoDialplanOrState checks that no server gives the spare an
// extension, hint, route, dev- entry or paging membership, and that no
// device state filter names it.
func TestSpareNoDialplanOrState(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)

	for _, serverID := range p.ServerIDs() {
		t.Run(serverID, func(t *testing.T) {
			extensions := astgen.ServerExtensions(data, p, serverID, "from-internal")
			fromTrunk := astgen.ServerTrunkExtensions(data, p, serverID, "default")
			paging := astgen.ServerPagingExtensions(data, p, serverID, "paging")
			relay := astgen.ServerPageRelayExtensions(data, p, serverID)
			hubFallthrough := astgen.ServerHubFallthrough(p, serverID)
			text := render(t, &extensions, &fromTrunk, &paging, &relay, &hubFallthrough)
			if strings.Contains(text, spareMAC) {
				t.Errorf("the dialplan of %s names the spare:\n%s", serverID, text)
			}

			cfg := astgen.ServerTrunks(p, serverID, pjsip.Endpoint{}, bindAddrs[serverID])
			for _, publish := range cfg.OutboundPublishes {
				if strings.Contains(publish.DeviceStateFilter, spareMAC) {
					t.Errorf("outbound-publish %s filter names the spare: %s", publish.Name, publish.DeviceStateFilter)
				}
			}
			for _, publication := range cfg.AsteriskPublications {
				if strings.Contains(publication.DeviceStateFilter, spareMAC) {
					t.Errorf("asterisk-publication %s filter names the spare: %s", publication.Name, publication.DeviceStateFilter)
				}
			}
			for _, name := range p.HostedDevices(serverID) {
				if name == spareMAC {
					t.Errorf("HostedDevices(%s) includes the spare", serverID)
				}
			}
		})
	}
}

// TestSpareSingleServer checks that the single-server fixture holds a
// spare and that the hub renders it, so TestSingleServerEquivalence covers
// a spare against the legacy unassigned-phone output.
func TestSpareSingleServer(t *testing.T) {
	data := singleServerData()
	p, problems := place(t, data)
	if len(problems) != 0 {
		t.Fatalf("Place() problems: %v", problems)
	}
	if len(p.Spares) == 0 {
		t.Fatal("the single-server fixture has no spare")
	}
	for _, mac := range p.Spares {
		found := false
		for _, phone := range astgen.ServerPhones(data, p, "qry", dpma.Phone{}, "") {
			if phone.MAC == mac {
				found = true
			}
		}
		if !found {
			t.Errorf("ServerPhones(qry) has no entry for the spare %s", mac)
		}
	}
}
