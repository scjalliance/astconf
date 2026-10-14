package astgen_test

import (
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
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
	got := renderTrunks(t, astgen.ServerTrunks(p, "qry", "100.64.0.3"))
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
	cfg := astgen.ServerTrunks(p, "hub", "100.64.0.1")
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

	cfg := astgen.ServerTrunks(p, "lcy", "100.64.0.4")
	if got := len(cfg.Endpoints); got != 3 {
		t.Errorf("lcy has %d trunk endpoints, want 3", got)
	}
	for _, publish := range cfg.OutboundPublishes {
		if publish.DeviceStateFilter != "^$" {
			t.Errorf("%s @device_state_filter = %q, want ^$", publish.Name, publish.DeviceStateFilter)
		}
	}
	found := false
	for _, publication := range astgen.ServerTrunks(p, "qry", "100.64.0.3").AsteriskPublications {
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
