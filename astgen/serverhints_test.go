package astgen_test

import (
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
)

// TestServerExtensionsHintOffice checks that an office has a hint for an
// extension homed on another server, naming that extension's devices there,
// so a busy-lamp key for a colleague elsewhere lights from published state.
func TestServerExtensionsHintOffice(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerExtensions(data, p, "qry", "from-internal")
	got := render(t, &section)
	for _, want := range []string{
		"exten => 101,hint,PJSIP/barney\n",     // Homed on den
		"exten => 102,hint,PJSIP/wilma.soft\n", // Homed on the hub
		"exten => 100,hint,PJSIP/fred&PJSIP/fred.soft\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("ServerExtensions(qry) rendered:\n%s\nwant it to contain:\n%s", got, want)
		}
	}
}

// TestServerExtensionsHintHub checks that the hub has hints for extensions
// homed on the offices.
func TestServerExtensionsHintHub(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerExtensions(data, p, "hub", "from-internal")
	got := render(t, &section)
	for _, want := range []string{
		"exten => 100,hint,PJSIP/fred&PJSIP/fred.soft\n",
		"exten => 101,hint,PJSIP/barney\n",
		"exten => 200,hint,PJSIP/lobby\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("ServerExtensions(hub) rendered:\n%s\nwant it to contain:\n%s", got, want)
		}
	}
}

// TestServerExtensionsNoHintWithoutDevices checks that an extension with no
// device still renders but has no hint, on its home and elsewhere.
func TestServerExtensionsNoHintWithoutDevices(t *testing.T) {
	data := multiServerData()
	data.People[2].Softphones = nil // Wilma has no device left
	p, _ := place(t, data)
	for _, serverID := range p.ServerIDs() {
		section := astgen.ServerExtensions(data, p, serverID, "from-internal")
		got := render(t, &section)
		if !strings.Contains(got, "exten => 102,1,") {
			t.Errorf("ServerExtensions(%s) has no extension 102:\n%s", serverID, got)
		}
		if strings.Contains(got, "102,hint") {
			t.Errorf("ServerExtensions(%s) has a hint for 102, which has no device:\n%s", serverID, got)
		}
	}
}

// TestServerTrunkExtensionsNoHints checks that hints stay in the phones'
// context and never appear in from-trunk.
func TestServerTrunkExtensionsNoHints(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	for _, serverID := range p.ServerIDs() {
		section := astgen.ServerTrunkExtensions(data, p, serverID, "default")
		if got := render(t, &section); strings.Contains(got, ",hint,") {
			t.Errorf("ServerTrunkExtensions(%s) has a hint:\n%s", serverID, got)
		}
	}
}
