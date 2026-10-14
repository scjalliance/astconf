package astgen_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/dpma"
	"github.com/scjalliance/astconf/pjsip"
)

// update rewrites the golden files under testdata from the current output.
// Run "go test ./astgen -run TestGolden -update", then review the diff
// before committing it.
var update = flag.Bool("update", false, "rewrite the golden files in testdata")

// bindAddrs holds the tailnet IP each fixture server binds its trunk
// transport to.
var bindAddrs = map[string]string{
	"den": "100.64.0.2",
	"hub": "100.64.0.1",
	"qry": "100.64.0.3",
}

// renderServer renders every per-server generator for serverID and returns
// the output keyed by file name.
func renderServer(t *testing.T, data *astorg.DataSet, p astgen.Placement, serverID string) map[string]string {
	t.Helper()

	endpoints := astgen.ServerEndpoints(data, p, serverID, pjsip.Endpoint{Context: "from-internal"}, "default")
	aors := astgen.ServerAORs(data, p, serverID, pjsip.AOR{MaxContacts: astval.NewInt(5), RemoveExisting: astval.Yes})
	trunks := astgen.ServerTrunks(p, serverID, bindAddrs[serverID])
	var pjsipValues []interface{}
	pjsipValues = append(pjsipValues, &trunks.Transport)
	for i := range endpoints {
		pjsipValues = append(pjsipValues, &endpoints[i])
	}
	for i := range aors {
		pjsipValues = append(pjsipValues, &aors[i])
	}
	for i := range trunks.Endpoints {
		pjsipValues = append(pjsipValues, &trunks.Endpoints[i])
	}
	for i := range trunks.AORs {
		pjsipValues = append(pjsipValues, &trunks.AORs[i])
	}
	for i := range trunks.Identifies {
		pjsipValues = append(pjsipValues, &trunks.Identifies[i])
	}
	for i := range trunks.OutboundPublishes {
		pjsipValues = append(pjsipValues, &trunks.OutboundPublishes[i])
	}
	for i := range trunks.InboundPublications {
		pjsipValues = append(pjsipValues, &trunks.InboundPublications[i])
	}
	for i := range trunks.AsteriskPublications {
		pjsipValues = append(pjsipValues, &trunks.AsteriskPublications[i])
	}

	extensions := astgen.ServerExtensions(data, p, serverID, "from-internal")
	paging := astgen.ServerPagingExtensions(data, p, serverID, "paging")
	relay := astgen.ServerPageRelayExtensions(data, p, serverID)
	alert := astgen.PageAlertSection()
	fromTrunk := astgen.ServerTrunkExtensions(data, p, serverID, "default")

	phones := astgen.ServerPhones(data, p, serverID, dpma.Phone{}, "https://contacts.example.com/")
	lines := astgen.ServerLines(data, p, serverID)
	var dpmaValues []interface{}
	for i := range phones {
		dpmaValues = append(dpmaValues, &phones[i])
	}
	for i := range lines {
		dpmaValues = append(dpmaValues, &lines[i])
	}

	mailboxes := astgen.ServerMailboxes(data, p, serverID, "default")

	return map[string]string{
		"pjsip.conf":            render(t, pjsipValues...),
		"extensions.conf":       render(t, &extensions, &paging, &relay, &alert, &fromTrunk),
		"res_digium_phone.conf": render(t, dpmaValues...),
		"voicemail.conf":        render(t, &mailboxes),
	}
}

// TestGolden renders every server in the multi-server fixture and compares
// each file with its golden copy in testdata/<server-id>/.
func TestGolden(t *testing.T) {
	data := multiServerData()
	p, problems, err := astgen.Place(data)
	if err != nil {
		t.Fatalf("Place() error: %v", err)
	}
	if len(problems) != 0 {
		t.Fatalf("Place() problems: %v", problems)
	}

	for _, serverID := range p.ServerIDs() {
		for name, got := range renderServer(t, data, p, serverID) {
			path := filepath.Join("testdata", serverID, name)
			if *update {
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatalf("MkdirAll() error: %v", err)
				}
				if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
					t.Fatalf("WriteFile() error: %v", err)
				}
				continue
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile() error: %v (run with -update to create it)", err)
			}
			if got != string(want) {
				t.Errorf("%s does not match its golden file.\ngot:\n%s\nwant:\n%s", path, got, want)
			}
		}
	}
}

// TestGoldenDeterministic renders the fixture ten times, with a fresh
// placement each time, and checks that every file is the same byte for byte.
func TestGoldenDeterministic(t *testing.T) {
	first := renderAll(t)
	// Go randomizes map iteration order, so repeated renders catch a map
	// that is ranged over without sorting.
	for run := 1; run < 10; run++ {
		again := renderAll(t)
		if len(again) != len(first) {
			t.Fatalf("render %d produced %d files, want %d", run, len(again), len(first))
		}
		for path, want := range first {
			if got := again[path]; got != want {
				t.Fatalf("render %d of %s differs:\n%s\nfirst render:\n%s", run, path, got, want)
			}
		}
	}
}

// renderAll places the fixture and renders every server, keyed by
// "<server-id>/<file>".
func renderAll(t *testing.T) map[string]string {
	t.Helper()
	data := multiServerData()
	p, _, err := astgen.Place(data)
	if err != nil {
		t.Fatalf("Place() error: %v", err)
	}
	out := make(map[string]string)
	for _, serverID := range p.ServerIDs() {
		for name, text := range renderServer(t, data, p, serverID) {
			out[serverID+"/"+name] = text
		}
	}
	return out
}

// containsToken reports whether text holds code as a whole token, not as a
// piece of a longer run of letters or digits. Short numeric codes would
// otherwise match inside MAC addresses.
func containsToken(text, code string) bool {
	isWord := func(b byte) bool {
		return b >= '0' && b <= '9' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
	}
	for from := 0; ; {
		i := strings.Index(text[from:], code)
		if i < 0 {
			return false
		}
		start, end := from+i, from+i+len(code)
		if (start == 0 || !isWord(text[start-1])) && (end == len(text) || !isWord(text[end])) {
			return true
		}
		from = start + 1
	}
}

// TestGoldenNoPasswords checks that no per-server output carries a secret
// from the fixture: no phone or softphone secret, no VMCODE variable, no
// voicemail code or mailbox access code, and no password setting other than
// the locked sentinel in a mailbox line.
func TestGoldenNoPasswords(t *testing.T) {
	data := multiServerData()

	// Secrets long and distinctive enough to match as plain substrings.
	secrets := []string{"VMCODE"}
	for _, phone := range data.Phones {
		secrets = append(secrets, phone.Secret)
	}
	for _, softphone := range data.Softphones {
		secrets = append(secrets, softphone.Secret)
	}
	// Short numeric codes match only as whole tokens.
	var codes []string
	for _, person := range data.People {
		codes = append(codes, person.VoicemailCode)
	}
	for _, mailbox := range data.Mailboxes {
		codes = append(codes, mailbox.AccessCode)
	}
	if len(secrets) < 2 || len(codes) < 2 {
		t.Fatal("the fixture lists no secrets or codes to check for")
	}

	for path, text := range renderAll(t) {
		for _, secret := range secrets {
			if secret != "" && strings.Contains(text, secret) {
				t.Errorf("%s contains the secret %q", path, secret)
			}
		}
		for _, code := range codes {
			if code != "" && containsToken(text, code) {
				t.Errorf("%s contains the code %q", path, code)
			}
		}
		if strings.Contains(strings.ToLower(text), "password") {
			t.Errorf("%s contains a password setting", path)
		}
		if !strings.HasSuffix(path, "/voicemail.conf") {
			continue
		}
		// A mailbox line is "<number> => <password>,<name>,...".
		for _, line := range strings.Split(text, "\n") {
			_, rest, ok := strings.Cut(line, " => ")
			if !ok {
				continue
			}
			password, _, _ := strings.Cut(rest, ",")
			if password != "-"+astgen.LockedMailboxPassword {
				t.Errorf("%s has a mailbox password other than the sentinel: %q", path, line)
			}
		}
	}
}
