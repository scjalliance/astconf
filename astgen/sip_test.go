package astgen_test

import (
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/sip"
)

func TestSIP(t *testing.T) {
	base := sip.Entity{
		Type:    sip.Friend,
		Context: "from-internal",
		Host:    "dynamic",
	}
	entities := astgen.SIP(testData(), base, "default")

	var values []interface{}
	for i := range entities {
		values = append(values, &entities[i])
	}
	got := render(t, values...)

	want := "[fred](desk)\n" +
		"type = friend\n" +
		"callerid = \"Fred Flintstone\" <100>\n" +
		"context = from-internal\n" +
		"host = dynamic\n" +
		"mailbox = 100@default\n" +
		"secret = fred-secret\n" +
		"setvar = VMCODE=1234\n" +
		"setvar = MOBILE=3605551234\n" +
		"setvar = USERNAME=fred\n" +
		"setvar = USER_LOCATION=Quarry\n" +
		"setvar = OUTBOUND_CALLERID=3605550100\n" +
		"setvar = AREACODE=360\n" +
		"\n" +
		"[lobby]\n" +
		"type = friend\n" +
		"callerid = \"Lobby\" <200>\n" +
		"context = from-internal\n" +
		"host = dynamic\n" +
		"mailbox = 200@default\n" +
		"secret = lobby-secret\n" +
		"setvar = USER_LOCATION=Quarry\n" +
		"setvar = OUTBOUND_CALLERID=3605550100\n" +
		"setvar = AREACODE=360\n" +
		"\n" +
		"[000fd3000003]\n" +
		"type = friend\n" +
		"callerid = \"QRY-000FD3000003\" <UNAVAILABLE>\n" +
		"context = from-internal\n" +
		"host = dynamic\n" +
		"secret = spare-secret\n" +
		"setvar = USER_LOCATION=Quarry\n" +
		"setvar = OUTBOUND_CALLERID=3605550100\n" +
		"setvar = AREACODE=360\n" +
		"\n" +
		"[fred.soft](soft)\n" +
		"type = friend\n" +
		"callerid = \"Fred Flintstone\" <100>\n" +
		"context = from-internal\n" +
		"host = dynamic\n" +
		"secret = fred-soft-secret\n" +
		"setvar = VMCODE=1234\n" +
		"setvar = MOBILE=3605551234\n" +
		"setvar = USERNAME=fred\n" +
		"setvar = USER_LOCATION=Quarry\n" +
		"setvar = OUTBOUND_CALLERID=3605550100\n" +
		"setvar = AREACODE=360\n" +
		"\n" +
		"[lobby.soft]\n" +
		"type = friend\n" +
		"callerid = \"Lobby\" <200>\n" +
		"context = from-internal\n" +
		"host = dynamic\n" +
		"secret = lobby-soft-secret\n"
	if got != want {
		t.Errorf("SIP() rendered:\n%s\nwant:\n%s", got, want)
	}
}
