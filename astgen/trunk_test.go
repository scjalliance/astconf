package astgen_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

// renderTrunks renders every section of a trunk configuration in the order
// a pjsip.conf would hold them.
func renderTrunks(t *testing.T, cfg astgen.TrunkConfig) string {
	t.Helper()
	values := []interface{}{&cfg.Transport}
	for i := range cfg.Endpoints {
		values = append(values, &cfg.Endpoints[i])
	}
	for i := range cfg.AORs {
		values = append(values, &cfg.AORs[i])
	}
	for i := range cfg.Identifies {
		values = append(values, &cfg.Identifies[i])
	}
	for i := range cfg.OutboundPublishes {
		values = append(values, &cfg.OutboundPublishes[i])
	}
	for i := range cfg.InboundPublications {
		values = append(values, &cfg.InboundPublications[i])
	}
	for i := range cfg.AsteriskPublications {
		values = append(values, &cfg.AsteriskPublications[i])
	}
	return render(t, values...)
}

// TestServerTrunksOffice renders the trunks of an office server. Its peers
// come in ID order, it publishes only its own devices, it accepts only each
// peer's devices, and it accepts mailbox state from the hub alone.
func TestServerTrunksOffice(t *testing.T) {
	p, _ := place(t, multiServerData())
	got := renderTrunks(t, astgen.ServerTrunks(p, "qry", pjsip.Endpoint{}, "100.64.0.3"))
	want := "[transport-tcp]\n" +
		"type = transport\n" +
		"protocol = tcp\n" +
		"bind = 100.64.0.3:5060\n" +
		"\n" +
		"[trunk-den]\n" +
		"type = endpoint\n" +
		"transport = transport-tcp\n" +
		"context = from-trunk\n" +
		"auth = auth-trunk-den\n" +
		"outbound_auth = auth-trunk-den\n" +
		"aors = trunk-den\n" +
		"direct_media = no\n" +
		"\n" +
		"[trunk-hub]\n" +
		"type = endpoint\n" +
		"transport = transport-tcp\n" +
		"context = from-trunk\n" +
		"auth = auth-trunk-hub\n" +
		"outbound_auth = auth-trunk-hub\n" +
		"aors = trunk-hub\n" +
		"direct_media = no\n" +
		"\n" +
		"[trunk-den]\n" +
		"type = aor\n" +
		"contact = sip:den-phone:5060\\;transport=tcp\n" +
		"\n" +
		"[trunk-hub]\n" +
		"type = aor\n" +
		"contact = sip:hub-phone:5060\\;transport=tcp\n" +
		"\n" +
		"[trunk-den]\n" +
		"type = identify\n" +
		"endpoint = trunk-den\n" +
		"match = den-phone\n" +
		"\n" +
		"[trunk-hub]\n" +
		"type = identify\n" +
		"endpoint = trunk-hub\n" +
		"match = hub-phone\n" +
		"\n" +
		"[publish-den]\n" +
		"type = outbound-publish\n" +
		"server_uri = sip:trunk-qry@den-phone:5060\\;transport=tcp\n" +
		"event = asterisk-devicestate\n" +
		"outbound_auth = auth-trunk-den\n" +
		"@device_state_filter = ^PJSIP/(000fd3000003|fred)$\n" +
		"\n" +
		"[publish-hub]\n" +
		"type = outbound-publish\n" +
		"server_uri = sip:trunk-qry@hub-phone:5060\\;transport=tcp\n" +
		"event = asterisk-devicestate\n" +
		"outbound_auth = auth-trunk-hub\n" +
		"@device_state_filter = ^PJSIP/(000fd3000003|fred)$\n" +
		"\n" +
		"[trunk-den]\n" +
		"type = inbound-publication\n" +
		"endpoint = trunk-den\n" +
		"event_asterisk-devicestate = den\n" +
		"\n" +
		"[trunk-hub]\n" +
		"type = inbound-publication\n" +
		"endpoint = trunk-hub\n" +
		"event_asterisk-devicestate = hub\n" +
		"event_asterisk-mwi = hub\n" +
		"\n" +
		"[den]\n" +
		"type = asterisk-publication\n" +
		"devicestate_publish = publish-den\n" +
		"device_state = yes\n" +
		"device_state_filter = ^PJSIP/(barney|lobby)$\n" +
		"\n" +
		"[hub]\n" +
		"type = asterisk-publication\n" +
		"devicestate_publish = publish-hub\n" +
		"device_state = yes\n" +
		"device_state_filter = ^PJSIP/(fred\\.soft|wilma\\.soft|yard\\.soft)$\n" +
		"mailbox_state = yes\n"
	if got != want {
		t.Errorf("ServerTrunks(qry) rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestServerTrunksHubPublishesMailboxes checks that the hub sends mailbox
// state to each peer through a second outbound-publish.
func TestServerTrunksHubPublishesMailboxes(t *testing.T) {
	p, _ := place(t, multiServerData())
	cfg := astgen.ServerTrunks(p, "hub", pjsip.Endpoint{}, "100.64.0.1")
	var names []string
	for _, publish := range cfg.OutboundPublishes {
		names = append(names, publish.Name+" "+publish.Event)
	}
	want := []string{
		"publish-den asterisk-devicestate",
		"publish-mwi-den asterisk-mwi",
		"publish-qry asterisk-devicestate",
		"publish-mwi-qry asterisk-mwi",
	}
	if len(names) != len(want) {
		t.Fatalf("OutboundPublishes = %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("OutboundPublishes[%d] = %q, want %q", i, names[i], want[i])
		}
	}
	for _, publication := range cfg.AsteriskPublications {
		if publication.MailboxStatePublish != "publish-mwi-"+publication.Name {
			t.Errorf("asterisk-publication %s mailboxstate_publish = %q", publication.Name, publication.MailboxStatePublish)
		}
		if publication.MailboxState.Specified() {
			t.Errorf("asterisk-publication %s accepts mailbox state on the hub", publication.Name)
		}
	}
}

// TestServerTrunksServerWithoutDevices checks a server with no locations.
// It still gets a trunk to every peer, and its device state filters match
// nothing, both on what it sends and on what its peers accept from it.
func TestServerTrunksServerWithoutDevices(t *testing.T) {
	data := multiServerData()
	data.Servers = append(data.Servers, astorg.Server{ID: "lcy", Name: "Lacey", Address: "lcy-phone"})
	p, _ := place(t, data)

	cfg := astgen.ServerTrunks(p, "lcy", pjsip.Endpoint{}, "100.64.0.4")
	if got := len(cfg.Endpoints); got != 3 {
		t.Errorf("lcy has %d trunk endpoints, want 3", got)
	}
	for _, publish := range cfg.OutboundPublishes {
		if publish.DeviceStateFilter != "^$" {
			t.Errorf("%s @device_state_filter = %q, want ^$", publish.Name, publish.DeviceStateFilter)
		}
	}
	found := false
	for _, publication := range astgen.ServerTrunks(p, "qry", pjsip.Endpoint{}, "100.64.0.3").AsteriskPublications {
		if publication.Name != "lcy" {
			continue
		}
		found = true
		if publication.DeviceStateFilter != "^$" {
			t.Errorf("qry accepts %q from lcy, want ^$", publication.DeviceStateFilter)
		}
	}
	if !found {
		t.Error("qry has no asterisk-publication for lcy")
	}
}

// TestServerTrunksBase checks that the base endpoint reaches every trunk,
// so trunks carry a codec list, and that the trunk's own settings win over
// any the base sets.
func TestServerTrunksBase(t *testing.T) {
	p, _ := place(t, multiServerData())
	base := pjsip.Endpoint{
		Transport:    "transport-udp",
		Context:      "from-internal",
		Auth:         "auth-base",
		OutboundAuth: "auth-base",
		AORs:         []string{"aor-base"},
		DirectMedia:  astval.Yes,
		Disallow:     []string{"all"},
		Allow:        []string{"ulaw", "alaw", "g722"},
	}
	cfg := astgen.ServerTrunks(p, "qry", base, "100.64.0.3")
	got := render(t, &cfg.Endpoints[0])
	want := "[trunk-den]\n" +
		"type = endpoint\n" +
		"transport = transport-tcp\n" +
		"context = from-trunk\n" +
		"disallow = all\n" +
		"allow = ulaw\n" +
		"allow = alaw\n" +
		"allow = g722\n" +
		"auth = auth-trunk-den\n" +
		"outbound_auth = auth-trunk-den\n" +
		"aors = trunk-den\n" +
		"direct_media = no\n"
	if got != want {
		t.Errorf("trunk endpoint rendered:\n%s\nwant:\n%s", got, want)
	}
	if len(base.AORs) != 1 || base.AORs[0] != "aor-base" {
		t.Errorf("ServerTrunks changed the base aors to %v", base.AORs)
	}
}

// TestServerTrunksNoAddress checks a server with an empty address. Place
// reports it, and it joins no trunk pair: no peer renders a trunk to it,
// it renders none, and no trunk credential names it.
func TestServerTrunksNoAddress(t *testing.T) {
	data := multiServerData()
	data.Servers[1].Address = "" // den
	p, problems := place(t, data)

	want := []astgen.Problem{{
		Rule:   astgen.RuleServerNoAddress,
		Kind:   "server",
		ID:     "den",
		Detail: "the server has no address, so it gets no trunk to any other server",
	}}
	if !reflect.DeepEqual(problems, want) {
		t.Errorf("Place() problems =\n%v\nwant:\n%v", problems, want)
	}
	if want := []astgen.TrunkPair{{ID: "hub~qry", A: "hub", B: "qry"}}; !reflect.DeepEqual(p.TrunkPairs, want) {
		t.Errorf("TrunkPairs = %v, want %v", p.TrunkPairs, want)
	}
	if p.Homes["101"] != "den" || p.Hosts["barney"] != "den" {
		t.Errorf("den lost its own extension or device: Homes[101] = %q, Hosts[barney] = %q", p.Homes["101"], p.Hosts["barney"])
	}

	for _, serverID := range []string{"qry", "hub"} {
		cfg := astgen.ServerTrunks(p, serverID, pjsip.Endpoint{}, "100.64.0.9")
		if text := renderTrunks(t, cfg); strings.Contains(text, "trunk-den") || strings.Contains(text, "sip::") {
			t.Errorf("%s renders a trunk to den:\n%s", serverID, text)
		}
		if len(cfg.Endpoints) != 1 {
			t.Errorf("%s has %d trunk endpoints, want 1", serverID, len(cfg.Endpoints))
		}
		for _, cred := range p.Credentials(serverID) {
			if strings.Contains(cred.ID, "den") {
				t.Errorf("%s has a credential for den: %v", serverID, cred)
			}
		}
	}
	cfg := astgen.ServerTrunks(p, "den", pjsip.Endpoint{}, "100.64.0.2")
	if len(cfg.Endpoints) != 0 || len(cfg.OutboundPublishes) != 0 || len(cfg.AsteriskPublications) != 0 {
		t.Errorf("den renders trunks with no address: %+v", cfg)
	}
	// den still writes the phone credentials of its own desk phones and of
	// the spare, which every server a location names lists, but no trunk
	// credential.
	wantCreds := []astgen.Credential{
		{Kind: "phone", ID: "000fd3000002", Section: "lobby-auth"},
		{Kind: "phone", ID: "000fd3000004", Section: "barney-auth"},
		{Kind: "phone", ID: "000fd3000009", Section: "000fd3000009-auth"},
	}
	if creds := p.Credentials("den"); !reflect.DeepEqual(creds, wantCreds) {
		t.Errorf("Credentials(den) = %v, want %v", creds, wantCreds)
	}
}

// TestServerTrunksIPv6Address checks that a peer address that is an IPv6
// literal is bracketed in the aor contact and the publish server_uri.
func TestServerTrunksIPv6Address(t *testing.T) {
	data := multiServerData()
	data.Servers[1].Address = "fd7a:115c:a1e0::2" // den
	p, _ := place(t, data)
	cfg := astgen.ServerTrunks(p, "qry", pjsip.Endpoint{}, "fd7a:115c:a1e0::3")

	if got, want := cfg.Transport.Bind, "[fd7a:115c:a1e0::3]:5060"; got != want {
		t.Errorf("transport bind = %q, want %q", got, want)
	}
	if got, want := cfg.AORs[0].Contacts, []string{"sip:[fd7a:115c:a1e0::2]:5060;transport=tcp"}; !reflect.DeepEqual(got, want) {
		t.Errorf("aor %s contact = %v, want %v", cfg.AORs[0].Name, got, want)
	}
	if got, want := cfg.OutboundPublishes[0].ServerURI, "sip:trunk-qry@[fd7a:115c:a1e0::2]:5060;transport=tcp"; got != want {
		t.Errorf("outbound-publish %s server_uri = %q, want %q", cfg.OutboundPublishes[0].Name, got, want)
	}
	if got, want := cfg.Identifies[0].Match, []string{"fd7a:115c:a1e0::2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("identify %s match = %v, want %v", cfg.Identifies[0].Name, got, want)
	}
}
