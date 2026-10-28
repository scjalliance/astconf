package astgen

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/dpma"
	"github.com/scjalliance/astconf/pjsip"
	"github.com/scjalliance/astconf/voicemail"
)

// hostedSet returns the set of device names that serverID renders: the
// devices it hosts and, on a server in SpareServers, every spare.
func hostedSet(p Placement, serverID string) map[string]bool {
	set := make(map[string]bool)
	for name, host := range p.Hosts {
		if host == serverID {
			set[name] = true
		}
	}
	for name := range p.sparesOn(serverID) {
		set[name] = true
	}
	return set
}

// placedData returns a shallow copy of data that holds only the devices
// Place accepted, for the device generators to render from. Each desk
// phone is kept only for the first record of its MAC address that Place
// hosted or kept as a spare, and each softphone only for the first record
// of its username that Place hosted on the hub. Each person's and role's
// softphone list is filtered the same way, so the owner of a rejected
// softphone never merges its caller ID or variables into another device
// of the same name, such as a spare. data itself is never changed.
func placedData(data *astorg.DataSet, p Placement) *astorg.DataSet {
	lookup := data.Lookup()
	spares := make(map[string]bool, len(p.Spares))
	for _, mac := range p.Spares {
		spares[mac] = true
	}

	out := *data

	out.Phones = nil
	keptPhone := make(map[string]bool)
	for _, phone := range data.Phones {
		if keptPhone[phone.MAC] {
			continue
		}
		var accepted bool
		if phone.Location == "" {
			accepted = spares[phone.MAC]
		} else if host, ok := p.Phones[phone.MAC]; ok {
			id, _, _ := lookup.LocationServer(phone.Location)
			accepted = id == host
		}
		if accepted {
			keptPhone[phone.MAC] = true
			out.Phones = append(out.Phones, phone)
		}
	}

	// A softphone is accepted when it is hosted on the hub and its username
	// is not a desk phone line or a spare's line.
	acceptedSoftphone := func(username string) bool {
		if _, line := p.LineMACs[username]; line {
			return false
		}
		host, ok := p.Hosts[username]
		return ok && host == p.Hub
	}
	out.Softphones = nil
	keptSoftphone := make(map[string]bool)
	for _, softphone := range data.Softphones {
		if keptSoftphone[softphone.Username] || !acceptedSoftphone(softphone.Username) {
			continue
		}
		keptSoftphone[softphone.Username] = true
		out.Softphones = append(out.Softphones, softphone)
	}

	filter := func(names []string) []string {
		var kept []string
		for _, name := range names {
			if keptSoftphone[name] {
				kept = append(kept, name)
			}
		}
		return kept
	}
	out.People = make(astorg.PersonList, len(data.People))
	for i, person := range data.People {
		person.Softphones = filter(person.Softphones)
		out.People[i] = person
	}
	out.PhoneRoles = make(astorg.PhoneRoleList, len(data.PhoneRoles))
	for i, role := range data.PhoneRoles {
		role.Softphones = filter(role.Softphones)
		out.PhoneRoles[i] = role
	}
	return &out
}

// ServerEndpoints generates PJSIP endpoint entries for the devices that
// serverID hosts and, on a server in SpareServers, for every spare. Each
// entry is the one Endpoints generates for the same device, in the same
// order, except that the VMCODE variable is removed.  Per-server output
// is committed to a config repo that every server pulls, so it never
// holds a voicemail PIN.
func ServerEndpoints(data *astorg.DataSet, p Placement, serverID string, base pjsip.Endpoint, context string) []pjsip.Endpoint {
	hosted := hostedSet(p, serverID)
	var out []pjsip.Endpoint
	for _, endpoint := range Endpoints(placedData(data, p), base, context) {
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
// hosts and, on a server in SpareServers, for every spare. Each entry is
// the one AORs generates for the same device, in the same order.
func ServerAORs(data *astorg.DataSet, p Placement, serverID string, base pjsip.AOR) []pjsip.AOR {
	hosted := make(map[string]bool)
	for name := range hostedSet(p, serverID) {
		hosted[aorName(name)] = true
	}
	var out []pjsip.AOR
	for _, aor := range AORs(placedData(data, p), base) {
		if hosted[aor.Name] {
			out = append(out, aor)
		}
	}
	return out
}

// ServerLines generates DPMA line entries for the desk phones that serverID
// hosts and, on a server in SpareServers, for every spare. Each entry is
// the one Lines generates for the same line, in the same order.
func ServerLines(data *astorg.DataSet, p Placement, serverID string) []dpma.Line {
	hosted := hostedSet(p, serverID)
	var out []dpma.Line
	for _, line := range Lines(placedData(data, p)) {
		if hosted[line.Name] {
			out = append(out, line)
		}
	}
	return out
}

// ServerPhones generates DPMA phone entries for the desk phones that
// serverID hosts and, on a server in SpareServers, for every spare. Each
// entry is the one Phones generates for the same phone, in the same order.
func ServerPhones(data *astorg.DataSet, p Placement, serverID string, base dpma.Phone, contactsURL string) []dpma.Phone {
	spares := p.sparesOn(serverID)
	var out []dpma.Phone
	for _, phone := range Phones(placedData(data, p), base, contactsURL) {
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
