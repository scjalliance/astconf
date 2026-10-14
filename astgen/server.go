package astgen

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/dpma"
	"github.com/scjalliance/astconf/pjsip"
	"github.com/scjalliance/astconf/voicemail"
)

// hostedSet returns the set of device names that serverID renders: the
// devices it hosts and, on every server in the placement, every spare.
func hostedSet(p Placement, serverID string) map[string]bool {
	set := make(map[string]bool)
	for name, host := range p.Hosts {
		if host == serverID {
			set[name] = true
		}
	}
	if _, ok := p.Servers[serverID]; ok {
		for name := range p.spareSet() {
			set[name] = true
		}
	}
	return set
}

// ServerEndpoints generates PJSIP endpoint entries for the devices that
// serverID hosts and for every spare. Each entry is the one Endpoints
// generates for the same device, in the same order, except that the VMCODE
// variable is removed.
// Per-server output is committed to a config repo that every server pulls,
// so it never holds a voicemail PIN.
func ServerEndpoints(data *astorg.DataSet, p Placement, serverID string, base pjsip.Endpoint, context string) []pjsip.Endpoint {
	hosted := hostedSet(p, serverID)
	var out []pjsip.Endpoint
	for _, endpoint := range Endpoints(data, base, context) {
		if !hosted[endpoint.Name] {
			continue
		}
		var vars []astval.Var
		for _, v := range endpoint.Variables {
			if v.Name != "VMCODE" {
				vars = append(vars, v)
			}
		}
		endpoint.Variables = vars
		out = append(out, endpoint)
	}
	return out
}

// ServerAORs generates PJSIP aor entries for the devices that serverID
// hosts and for every spare. Each entry is the one AORs generates for the
// same device, in the same order.
func ServerAORs(data *astorg.DataSet, p Placement, serverID string, base pjsip.AOR) []pjsip.AOR {
	hosted := make(map[string]bool)
	for name := range hostedSet(p, serverID) {
		hosted[aorName(name)] = true
	}
	var out []pjsip.AOR
	for _, aor := range AORs(data, base) {
		if hosted[aor.Name] {
			out = append(out, aor)
		}
	}
	return out
}

// ServerLines generates DPMA line entries for the desk phones that serverID
// hosts and for every spare. Each entry is the one Lines generates for the
// same line, in the same order.
func ServerLines(data *astorg.DataSet, p Placement, serverID string) []dpma.Line {
	hosted := hostedSet(p, serverID)
	var out []dpma.Line
	for _, line := range Lines(data) {
		if hosted[line.Name] {
			out = append(out, line)
		}
	}
	return out
}

// ServerPhones generates DPMA phone entries for the desk phones that
// serverID hosts and for every spare. Each entry is the one Phones
// generates for the same phone, in the same order.
func ServerPhones(data *astorg.DataSet, p Placement, serverID string, base dpma.Phone, contactsURL string) []dpma.Phone {
	var spares map[string]bool
	if _, ok := p.Servers[serverID]; ok {
		spares = p.spareSet()
	}
	var out []dpma.Phone
	for _, phone := range Phones(data, base, contactsURL) {
		if host, ok := p.Phones[phone.MAC]; (ok && host == serverID) || spares[phone.MAC] {
			out = append(out, phone)
		}
	}
	return out
}

// LockedMailboxPassword is the placeholder password every mailbox from
// ServerMailboxes carries. With passwordlocation=spooldir, Asterisk falls
// back to the voicemail.conf password when a mailbox's secret.conf is
// missing. An empty password would leave that mailbox open, so the fallback
// is a value no one can enter: "x" is not a DTMF digit, so it cannot be typed
// on a phone keypad. It is fixed and public, not a secret. A present
// secret.conf overrides it.
const LockedMailboxPassword = "x"

// ServerMailboxes generates the voicemail section for serverID. Voicemail
// lives on the hub alone, so every other server gets an empty section with
// the same context.
//
// The section never holds a real PIN: every mailbox carries
// LockedMailboxPassword. The hub reads PINs from the voicemail spool
// directory (voicemail.conf sets passwordlocation=spooldir in the
// hand-written common config) and its node agent writes each mailbox's
// secret.conf.
func ServerMailboxes(data *astorg.DataSet, p Placement, serverID string, context string) voicemail.Section {
	if serverID != p.Hub {
		return voicemail.Section{Context: context}
	}
	section := Mailboxes(data, context)
	for i := range section.Mailboxes {
		section.Mailboxes[i].Password = LockedMailboxPassword
	}
	return section
}
