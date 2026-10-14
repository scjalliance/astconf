package astgen_test

import (
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
)

// TestServerPagingAllMembersRemote checks a paging group with no members on
// the hub. A page from the hub reaches each office through one relay leg,
// and each leg is a local channel, not the trunk.
func TestServerPagingAllMembersRemote(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerPagingExtensions(data, p, "hub", "paging")
	got := render(t, &section)
	want := "[paging]\n" +
		"\n" +
		"; All Call\n" +
		"exten => 500,1,Noop(Paging All Call)\n" +
		"same => n,Playback(beep)\n" +
		"same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})\n" +
		"same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})\n" +
		"same => n,Page(Local/page-relay-500-den@astgen-page-relay&Local/page-relay-500-qry@astgen-page-relay,isb(astgen-page-alert^s^1(intercom)))\n" +
		"same => n,Hangup()\n"
	if got != want {
		t.Errorf("ServerPagingExtensions(hub) rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestServerPagingOffice checks that an office pages its local members
// directly and relays to the one other office that hosts members.
func TestServerPagingOffice(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerPagingExtensions(data, p, "den", "paging")
	got := render(t, &section)
	want := "same => n,Page(PJSIP/barney&PJSIP/lobby&Local/page-relay-500-qry@astgen-page-relay,isb(astgen-page-alert^s^1(intercom)))\n"
	if !strings.Contains(got, want) {
		t.Errorf("ServerPagingExtensions(den) rendered:\n%s\nwant it to contain:\n%s", got, want)
	}
	if strings.Contains(got, "SIPAddHeader") {
		t.Errorf("ServerPagingExtensions(den) uses SIPAddHeader, which does nothing on PJSIP:\n%s", got)
	}
}

// TestServerPagingNoAlert checks that a paging group with no alert pages
// without the predial subroutine.
func TestServerPagingNoAlert(t *testing.T) {
	data := multiServerData()
	data.PagingGroups[0].Alert = ""
	p, _ := place(t, data)
	section := astgen.ServerPagingExtensions(data, p, "den", "paging")
	want := "same => n,Page(PJSIP/barney&PJSIP/lobby&Local/page-relay-500-qry@astgen-page-relay,is)\n"
	if got := render(t, &section); !strings.Contains(got, want) {
		t.Errorf("ServerPagingExtensions(den) rendered:\n%s\nwant it to contain:\n%s", got, want)
	}
}

// TestServerPageRelayExtensions checks the relay context. Each relay
// extension dials one server's from-trunk page- entry over the trunk; a
// server with no relays renders an empty context.
func TestServerPageRelayExtensions(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)

	hub := astgen.ServerPageRelayExtensions(data, p, "hub")
	got := render(t, &hub)
	want := "[astgen-page-relay]\n" +
		"\n" +
		"; Relay All Call to den\n" +
		"exten => page-relay-500-den,1,Dial(PJSIP/page-500@trunk-den)\n" +
		"\n" +
		"; Relay All Call to qry\n" +
		"exten => page-relay-500-qry,1,Dial(PJSIP/page-500@trunk-qry)\n"
	if got != want {
		t.Errorf("ServerPageRelayExtensions(hub) rendered:\n%s\nwant:\n%s", got, want)
	}

	data.Locations[1].PagingGroups = nil // Den's phones leave the group
	p, _ = place(t, data)
	qry := astgen.ServerPageRelayExtensions(data, p, "qry")
	if got, want := render(t, &qry), "[astgen-page-relay]\n"; got != want {
		t.Errorf("ServerPageRelayExtensions(qry) with no remote members rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestPageAlertSection checks the predial subroutine that adds Alert-Info
// on PJSIP callee channels only.
func TestPageAlertSection(t *testing.T) {
	section := astgen.PageAlertSection()
	got := render(t, &section)
	want := "[astgen-page-alert]\n" +
		"\n" +
		"; Add a paging group's Alert-Info header to a paged phone\n" +
		"exten => s,1,ExecIf($[\"${CHANNEL(channeltype)}\"=\"PJSIP\"]?Set(PJSIP_HEADER(add,Alert-Info)=${ARG1}))\n" +
		"same => n,Return()\n"
	if got != want {
		t.Errorf("PageAlertSection() rendered:\n%s\nwant:\n%s", got, want)
	}
}
