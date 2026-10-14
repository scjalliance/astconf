package astgen_test

import (
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
)

// TestServerExtensionsOffice renders the phones' context of an office. Its
// own extension rings its local desk phone when idle and the owner's
// softphone on the hub, every other extension routes to its home, and
// unknown extensions go to the hub.
func TestServerExtensionsOffice(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerExtensions(data, p, "qry", "from-internal")
	got := render(t, &section)
	want := "[from-internal]\n" +
		"\n" +
		"; Fred Flintstone\n" +
		"exten => 100,1,Noop(Call fred)\n" +
		"same => n,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Set(ASTGEN_DIAL=&PJSIP/dev-fred.soft@trunk-hub)\n" +
		"same => n,ExecIf($[${DEVICE_STATE(PJSIP/fred)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/fred))\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-100@trunk-hub))\n" +
		"\n" +
		"; Barney Rubble\n" +
		"exten => 101,1,Dial(PJSIP/101@trunk-den)\n" +
		"\n" +
		"; Wilma Flintstone\n" +
		"exten => 102,1,Dial(PJSIP/102@trunk-hub)\n" +
		"\n" +
		"; Lobby\n" +
		"exten => 200,1,Dial(PJSIP/200@trunk-den)\n" +
		"\n" +
		"; Yard\n" +
		"exten => 201,1,Dial(PJSIP/201@trunk-hub)\n" +
		"\n" +
		"; Unknown extensions go to the hub\n" +
		"exten => _X.,1,Dial(PJSIP/${EXTEN}@trunk-hub)\n"
	if got != want {
		t.Errorf("ServerExtensions(qry) rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestServerExtensionsHub checks the hub's own extensions. Wilma has no
// mailbox, so the no-answer subroutine gets no target. The yard role has no
// mailbox, so it ends in congestion as a role does today. The hub has no
// fallthrough, so an unknown extension is invalid.
func TestServerExtensionsHub(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerExtensions(data, p, "hub", "from-internal")
	got := render(t, &section)
	want := "[from-internal]\n" +
		"\n" +
		"; Fred Flintstone\n" +
		"exten => 100,1,Dial(PJSIP/100@trunk-qry)\n" +
		"\n" +
		"; Barney Rubble\n" +
		"exten => 101,1,Dial(PJSIP/101@trunk-den)\n" +
		"\n" +
		"; Wilma Flintstone\n" +
		"exten => 102,1,Noop(Call wilma)\n" +
		"same => n,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Set(ASTGEN_DIAL=)\n" +
		"same => n,ExecIf($[${DEVICE_STATE(PJSIP/wilma.soft)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/wilma.soft))\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Gosub(no-answer,s,1(${EXTEN}))\n" +
		"\n" +
		"; Lobby\n" +
		"exten => 200,1,Dial(PJSIP/200@trunk-den)\n" +
		"\n" +
		"; Yard\n" +
		"exten => 201,1,Noop(Call yard)\n" +
		"same => n,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Set(ASTGEN_DIAL=)\n" +
		"same => n,ExecIf($[${DEVICE_STATE(PJSIP/yard.soft)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/yard.soft))\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Congestion()\n"
	if got != want {
		t.Errorf("ServerExtensions(hub) rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestHomeExtensionOnlyRemoteDevices checks an extension that lives on an
// office whose owner has no desk phone there, only a softphone on the hub.
// The ring has one remote leg, no local leg, and still dials.
func TestHomeExtensionOnlyRemoteDevices(t *testing.T) {
	data := multiServerData()
	data.Phones = data.Phones[1:] // Drop fred's desk phone
	data.People[0].Phones = nil
	p, _ := place(t, data)
	section := astgen.ServerExtensions(data, p, "qry", "from-internal")
	got := render(t, &section)
	want := "; Fred Flintstone\n" +
		"exten => 100,1,Noop(Call fred)\n" +
		"same => n,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Set(ASTGEN_DIAL=&PJSIP/dev-fred.soft@trunk-hub)\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-100@trunk-hub))\n" +
		"\n"
	if !strings.Contains(got, want) {
		t.Errorf("ServerExtensions(qry) rendered:\n%s\nwant it to contain:\n%s", got, want)
	}
}

// TestServerTrunkExtensionsOffice renders an office's from-trunk context:
// its home extension, a dev- entry per hosted device, and the page- entry
// for its paging members.
func TestServerTrunkExtensionsOffice(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerTrunkExtensions(data, p, "qry", "default")
	got := render(t, &section)
	want := "[from-trunk]\n" +
		"\n" +
		"; Fred Flintstone\n" +
		"exten => 100,1,Noop(Call fred)\n" +
		"same => n,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Set(ASTGEN_DIAL=&PJSIP/dev-fred.soft@trunk-hub)\n" +
		"same => n,ExecIf($[${DEVICE_STATE(PJSIP/fred)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/fred))\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-100@trunk-hub))\n" +
		"\n" +
		"; Ring 000fd3000003 for another server\n" +
		"exten => dev-000fd3000003,1,ExecIf($[${DEVICE_STATE(PJSIP/000fd3000003)}!=NOT_INUSE]?Busy())\n" +
		"same => n,Dial(PJSIP/000fd3000003)\n" +
		"same => n,Hangup()\n" +
		"\n" +
		"; Ring fred for another server\n" +
		"exten => dev-fred,1,ExecIf($[${DEVICE_STATE(PJSIP/fred)}!=NOT_INUSE]?Busy())\n" +
		"same => n,Dial(PJSIP/fred)\n" +
		"same => n,Hangup()\n" +
		"\n" +
		"; All Call\n" +
		"exten => page-500,1,Noop(Paging All Call)\n" +
		"same => n,Playback(beep)\n" +
		"same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})\n" +
		"same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})\n" +
		"same => n,Page(PJSIP/000fd3000003&PJSIP/fred,isb(astgen-page-alert^s^1(intercom)))\n" +
		"same => n,Hangup()\n"
	if got != want {
		t.Errorf("ServerTrunkExtensions(qry) rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestServerTrunkExtensionsHubVoicemail checks the hub's vm- entries, one
// per mailbox, in the order Mailboxes writes them.
func TestServerTrunkExtensionsHubVoicemail(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerTrunkExtensions(data, p, "hub", "default")
	got := render(t, &section)
	want := "; Voicemail for Sales\n" +
		"exten => vm-300,1,VoiceMail(300@default,u)\n" +
		"same => n,Hangup()\n" +
		"\n" +
		"; Voicemail for fred\n" +
		"exten => vm-100,1,VoiceMail(100@default,u)\n" +
		"same => n,Hangup()\n" +
		"\n" +
		"; Voicemail for barney\n" +
		"exten => vm-101,1,VoiceMail(101@default,u)\n" +
		"same => n,Hangup()\n"
	if !strings.HasSuffix(got, want) {
		t.Errorf("ServerTrunkExtensions(hub) rendered:\n%s\nwant it to end with:\n%s", got, want)
	}
	office := astgen.ServerTrunkExtensions(data, p, "qry", "default")
	if text := render(t, &office); strings.Contains(text, "exten => vm-") {
		t.Errorf("ServerTrunkExtensions(qry) has a vm- entry:\n%s", text)
	}
}

// TestServerTrunkExtensionsNoLoops checks every server's from-trunk context.
// It may hold only extensions homed on that server and dev-, vm- and page-
// entries, and nothing in it may dial another trunk except a home
// extension's remote ring legs and voicemail target.
func TestServerTrunkExtensionsNoLoops(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	for _, serverID := range p.ServerIDs() {
		section := astgen.ServerTrunkExtensions(data, p, serverID, "default")
		for _, ext := range section.Extensions {
			switch {
			case strings.HasPrefix(ext.Number, "dev-"),
				strings.HasPrefix(ext.Number, "vm-"),
				strings.HasPrefix(ext.Number, "page-"):
			case p.Homes[ext.Number] == serverID:
			default:
				t.Errorf("%s from-trunk holds %s, which is not homed there", serverID, ext.Number)
			}
		}
		text := render(t, &section)
		for _, line := range strings.Split(text, "\n") {
			if strings.Contains(line, "Dial(PJSIP/") && strings.Contains(line, "@trunk-") {
				t.Errorf("%s from-trunk routes to another server: %s", serverID, line)
			}
		}
	}
}

// TestServerTrunkExtensionsNoPageWithoutMembers checks that a server with
// no members of a paging group has no page- entry for it, so a relayed page
// never arrives where nobody would hear it.
func TestServerTrunkExtensionsNoPageWithoutMembers(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerTrunkExtensions(data, p, "hub", "default")
	if text := render(t, &section); strings.Contains(text, "page-500,") {
		t.Errorf("hub from-trunk has a page- entry with no local members:\n%s", text)
	}
}
