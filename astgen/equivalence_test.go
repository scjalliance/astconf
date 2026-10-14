package astgen_test

import (
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/dialplan"
	"github.com/scjalliance/astconf/dpma"
	"github.com/scjalliance/astconf/pjsip"
)

// This file proves that with one server, which is also the hub, the
// per-server generators render what the whole-dataset generators render,
// apart from the deliberate differences below. Any other difference is a
// regression. The allowed differences are exactly these:
//
//   - Endpoints: per-server endpoints never carry the VMCODE set_var,
//     because per-server output goes to a config repo every server pulls
//     and must hold no voicemail PIN (ruling P1). VMCODE is removed from the
//     legacy side before comparing.
//   - Mailboxes: per-server hub mailboxes carry the locked sentinel password
//     astgen.LockedMailboxPassword, rendered as -x, instead of a PIN (ruling
//     P1). The legacy password is replaced with the sentinel before comparing.
//   - Paging: devices are named PJSIP/ and not SIP/, and the Alert-Info
//     header moves from SIPAddHeader, which does nothing on PJSIP, to Page's
//     predial subroutine. Both changes are applied to the legacy output.
//   - Phone extensions: only the numbers, comments and the set of devices
//     are compared. The ring logic differs on purpose: no call waiting, each
//     device drops out of the ring alone when busy, and Gosub replaces Macro.
//   - Phone extensions also carry busy-lamp hint lines (ruling P2). Hints
//     are ignored when comparing.

// singleServerData returns testData with one server, qry, that is also the
// hub, a paging group at Quarry, and a spare phone with no location and no
// owner. This is the current deployment: every device and extension is
// local, and the spare renders as the legacy unassigned phone does.
func singleServerData() *astorg.DataSet {
	data := testData()
	data.Servers = astorg.ServerList{{ID: "qry", Name: "Quarry", Hub: true, Address: "qry-phone"}}
	data.Locations[0].Server = "qry"
	data.Locations[0].PagingGroups = []string{"500"}
	data.PagingGroups = astorg.PagingGroupList{{Name: "All Call", Extension: "500", Alert: "intercom"}}
	data.Phones = append(data.Phones, astorg.Phone{MAC: "000fd3000009", Secret: "roaming-spare-secret", Templates: []string{"desk"}})
	return data
}

// withoutVMCODE returns a copy of the endpoints with the VMCODE variable
// removed. It fails the test when no endpoint had one, so the comparison
// cannot pass for want of anything to remove.
func withoutVMCODE(t *testing.T, endpoints []pjsip.Endpoint) []pjsip.Endpoint {
	t.Helper()
	removed := 0
	out := make([]pjsip.Endpoint, len(endpoints))
	for i, endpoint := range endpoints {
		var vars []astval.Var
		for _, v := range endpoint.Variables {
			if v.Name == "VMCODE" {
				removed++
				continue
			}
			vars = append(vars, v)
		}
		endpoint.Variables = vars
		out[i] = endpoint
	}
	if removed == 0 {
		t.Fatal("Endpoints() has no VMCODE variable, so the comparison proves nothing")
	}
	return out
}

// TestSingleServerEquivalence renders the single-server fixture through the
// per-server generators and checks the result against the existing
// generators, allowing only the differences listed at the top of this file.
func TestSingleServerEquivalence(t *testing.T) {
	data := singleServerData()
	p, problems, err := astgen.Place(data)
	if err != nil {
		t.Fatalf("Place() error: %v", err)
	}
	if len(problems) != 0 {
		t.Fatalf("Place() problems: %v", problems)
	}
	const id = "qry"

	t.Run("endpoints", func(t *testing.T) {
		base := pjsip.Endpoint{Context: "from-internal"}
		got := astgen.ServerEndpoints(data, p, id, base, "default")
		want := withoutVMCODE(t, astgen.Endpoints(data, base, "default"))
		if len(got) == 0 {
			t.Fatal("ServerEndpoints() returned nothing")
		}
		if a, b := renderEach(t, got), renderEach(t, want); a != b {
			t.Errorf("ServerEndpoints() rendered:\n%s\nEndpoints() rendered, without VMCODE:\n%s", a, b)
		}
	})

	t.Run("aors", func(t *testing.T) {
		base := pjsip.AOR{MaxContacts: astval.NewInt(5)}
		got := astgen.ServerAORs(data, p, id, base)
		want := astgen.AORs(data, base)
		if len(got) == 0 {
			t.Fatal("ServerAORs() returned nothing")
		}
		if a, b := renderEach(t, got), renderEach(t, want); a != b {
			t.Errorf("ServerAORs() rendered:\n%s\nAORs() rendered:\n%s", a, b)
		}
	})

	t.Run("lines", func(t *testing.T) {
		got := astgen.ServerLines(data, p, id)
		want := astgen.Lines(data)
		if len(got) == 0 {
			t.Fatal("ServerLines() returned nothing")
		}
		if a, b := renderEach(t, got), renderEach(t, want); a != b {
			t.Errorf("ServerLines() rendered:\n%s\nLines() rendered:\n%s", a, b)
		}
	})

	t.Run("phones", func(t *testing.T) {
		got := astgen.ServerPhones(data, p, id, dpma.Phone{}, "https://contacts.example.com/")
		want := astgen.Phones(data, dpma.Phone{}, "https://contacts.example.com/")
		if len(got) == 0 {
			t.Fatal("ServerPhones() returned nothing")
		}
		if a, b := renderEach(t, got), renderEach(t, want); a != b {
			t.Errorf("ServerPhones() rendered:\n%s\nPhones() rendered:\n%s", a, b)
		}
	})

	t.Run("mailboxes", func(t *testing.T) {
		got := astgen.ServerMailboxes(data, p, id, "default")
		want := astgen.Mailboxes(data, "default")
		if len(want.Mailboxes) == 0 {
			t.Fatal("Mailboxes() returned nothing")
		}
		for i := range want.Mailboxes {
			if want.Mailboxes[i].Password == astgen.LockedMailboxPassword {
				t.Fatalf("Mailboxes() box %s already has the sentinel password", want.Mailboxes[i].Extension)
			}
			want.Mailboxes[i].Password = astgen.LockedMailboxPassword
		}
		if a, b := render(t, &got), render(t, &want); a != b {
			t.Errorf("ServerMailboxes() rendered:\n%s\nMailboxes() rendered, with the sentinel password:\n%s", a, b)
		}
	})

	t.Run("paging", func(t *testing.T) {
		got := astgen.ServerPagingExtensions(data, p, id, "paging")
		want := astgen.PagingGroupExtensions(data, "paging")
		a := render(t, &got)
		// Apply the two deliberate changes to the old output: devices are
		// PJSIP/, and the Alert-Info header moves from SIPAddHeader, which
		// does nothing on PJSIP, to Page's predial subroutine.
		b := render(t, &want)
		b = strings.ReplaceAll(b, "SIP/", "PJSIP/")
		b = strings.ReplaceAll(b, "same => n,SIPAddHeader(Alert-Info:intercom)\n", "")
		b = strings.ReplaceAll(b, ",is)\n", ",isb(astgen-page-alert^s^1(intercom)))\n")
		if a != b {
			t.Errorf("ServerPagingExtensions() rendered:\n%s\nPagingGroupExtensions() changed:\n%s", a, b)
		}
		if !strings.Contains(b, "Page(PJSIP/") {
			t.Errorf("the old paging output has no Page line to compare:\n%s", b)
		}
	})

	t.Run("extensions", func(t *testing.T) {
		got := astgen.ServerExtensions(data, p, id, "from-internal")
		want := astgen.PhoneExtensions(data, "from-internal")
		if len(want.Extensions) == 0 {
			t.Fatal("PhoneExtensions() returned nothing")
		}
		if len(got.Extensions) != len(want.Extensions) {
			t.Fatalf("ServerExtensions() returned %d extensions, PhoneExtensions() returned %d", len(got.Extensions), len(want.Extensions))
		}
		for i := range want.Extensions {
			g, w := got.Extensions[i], want.Extensions[i]
			if g.Number != w.Number || g.Comment != w.Comment {
				t.Errorf("extension %d = %s %q, want %s %q", i, g.Number, g.Comment, w.Number, w.Comment)
			}
			// Hint lines are allowed to differ, so drop them.
			g.Hint = nil
			gotDevices := devicesIn(t, g, "PJSIP/")
			wantDevices := devicesIn(t, w, "SIP/")
			if len(wantDevices) == 0 {
				t.Errorf("extension %s has no devices in the old output", w.Number)
			}
			if strings.Join(gotDevices, " ") != strings.Join(wantDevices, " ") {
				t.Errorf("extension %s rings %v, want %v", g.Number, gotDevices, wantDevices)
			}
		}
		if text := render(t, &got); strings.Contains(text, "@trunk-") {
			t.Errorf("single-server extensions reach a trunk:\n%s", text)
		}
	})
}

// renderEach renders every element of a slice of section values in order.
func renderEach[T any](t *testing.T, items []T) string {
	t.Helper()
	values := make([]interface{}, 0, len(items))
	for i := range items {
		values = append(values, &items[i])
	}
	return render(t, values...)
}

// devicePattern matches a device in a dial string, such as SIP/fred.
var devicePattern = regexp.MustCompile(`(?:PJ)?SIP/[^,&)\s]+`)

// devicesIn returns the sorted, distinct resource names of the devices with
// the given technology prefix that an extension refers to.
func devicesIn(t *testing.T, ext dialplan.Extension, prefix string) []string {
	t.Helper()
	section := dialplan.Section{Context: "x", Extensions: []dialplan.Extension{ext}}
	seen := make(map[string]bool)
	for _, match := range devicePattern.FindAllString(render(t, &section), -1) {
		if strings.HasPrefix(match, prefix) {
			seen[strings.TrimPrefix(match, prefix)] = true
		}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
