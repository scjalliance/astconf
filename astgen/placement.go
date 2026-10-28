package astgen

import (
	"errors"
	"fmt"
	"sort"

	"github.com/scjalliance/astconf/astorg"
)

// Errors returned by Place. Nothing can route without exactly one hub.
var (
	ErrNoHub        = errors.New("astgen: the data set has no hub server")
	ErrMultipleHubs = errors.New("astgen: the data set has more than one hub server")
)

// Problem rules reported by Place.
const (
	RuleLocationNoServer   = "location-no-server"  // A desk phone at a location with no server
	RuleServerMissing      = "server-missing"      // A location names a server that is not in the data set
	RuleExtensionDuplicate = "extension-duplicate" // Two people or roles claim the same extension
	RuleDeviceConflict     = "device-conflict"     // Desk phones that share a line name sit on different servers, or a softphone reuses a device name
	RuleServerNoAddress    = "server-no-address"   // A server has an empty Address, so no peer can reach it
	RuleOwnedNoLocation    = "owned-no-location"   // A desk phone with no location that a person or role owns
)

// Problem describes a record that Place left out of the placement, or a
// server with no address that Place left out of every trunk pair. No
// field ever carries a secret.
type Problem struct {
	Rule   string // One of the Rule constants
	Kind   string // "phone", "softphone", "person", "role" or "server"
	ID     string // The MAC address, the softphone's, person's or role's username, or the server ID
	Detail string
}

// String returns a one-line description of the problem.
func (p Problem) String() string {
	return fmt.Sprintf("%s: %s %s: %s", p.Rule, p.Kind, p.ID, p.Detail)
}

// TrunkPair is an unordered pair of servers joined by a trunk.
type TrunkPair struct {
	ID string // The two server IDs, sorted, joined by "~"
	A  string // The server ID that sorts first
	B  string // The server ID that sorts second
}

// Credential names a secret that the node agent writes into a server's
// local secrets file. astgen never renders the secret itself.
//
// For a trunk, the node agent writes the whole auth section: type=auth,
// auth_type=userpass, a username and the password derived from ID. The
// trunk endpoint uses that one section both as its auth, to check the
// peer's requests, and as its outbound_auth, to answer the peer's
// challenges. So the section on one end of a pair must hold the same
// username and password as the section on the other end, or calls fail
// in both directions. The expected username is the pair ID, which both
// ends already share.
//
// For a desk phone, the node agent writes the auth section of the phone's
// line, with the secret derived from ID, the MAC address. The section is
// the one the line's pjsip endpoint names: "<line>-auth", which for an
// unassigned phone or a spare is "<MAC>-auth". When several phones share
// a line, one credential covers the line, derived from the MAC address of
// the first of them in data set order that Place hosts or keeps as a
// spare. A phone that Place reports as a problem, such as a
// device-conflict, is skipped. A consumer must write the section that
// Section names and never build the name itself.
type Credential struct {
	Kind    string // "trunk" or "phone"
	ID      string // The trunk pair ID, from which both ends derive the same secret, or a desk phone's MAC address
	Section string // The pjsip auth section the secret is written to
}

// Placement records which server hosts each device and which server owns
// each extension. It is computed once per data set by Place and shared by
// every per-server generator.
//
// Device names are pjsip endpoint names: the line name of a desk phone and
// the username of a softphone.
type Placement struct {
	Hub        string                         // ID of the hub server
	Servers    map[string]astorg.Server       // Every server, keyed by ID
	Phones     map[string]string              // Desk phone MAC address to host server ID
	Hosts      map[string]string              // Device name to host server ID
	Homes      map[string]string              // Extension to home server ID
	Owners     map[string]string              // Extension to the username of the person or role that holds it
	Devices    map[string][]string            // Extension to the names of the devices it rings, sorted
	Remote     map[string]map[string][]string // Extension to host server ID to the devices it rings there, sorted, for every host other than its home
	TrunkPairs []TrunkPair                    // Every unordered pair of servers that both have an address, sorted by ID
	Spares     []string                       // MAC addresses of the spare desk phones, sorted; every server renders each one
	LineMACs   map[string]string              // Line name of each hosted desk phone and spare to the MAC address its credential derives from: the first phone on the line, in data set order, that Place hosts or keeps as a spare
}

// Place decides which server hosts each device and which server owns each
// extension.
//
// A desk phone is hosted by its location's server. Every softphone is
// hosted by the hub. A person's or role's extension lives on its location's
// server, or on the hub when that location has no server.
//
// A spare is a desk phone with an empty Location that no person or role
// owns. It may be plugged in at any office, so it is hosted by no server
// and every server renders it, as the whole-dataset generators render an
// unassigned phone: a DPMA phone with no extension and no user whose line,
// pjsip endpoint and aor are named by its MAC address. A spare rings for
// no extension, has no hint and no route, joins no paging group, and is
// left out of every device state filter, so no server publishes state for
// a device that may register to any of them. A desk phone with an empty
// Location that a person or role owns is not a spare: it is reported as
// RuleOwnedNoLocation and left out, since its extension's home cannot
// tell which server should ring it.
//
// Place returns an error only when the data set does not have exactly one
// hub. Any other record that cannot be placed is reported as a Problem and
// left out of the placement.
//
// A server with no address is reported as a Problem and joins no trunk
// pair, since no peer could reach it. It still hosts its devices and homes
// its extensions, so its own output renders, but no server renders a trunk
// to it and it renders none, and calls that other servers route to it fail
// until it has an address. A consumer should treat that problem as one that
// blocks a commit.
func Place(data *astorg.DataSet) (Placement, []Problem, error) {
	lookup := data.Lookup()

	var hubs []string
	for id, server := range lookup.ServerByID {
		if server.Hub {
			hubs = append(hubs, id)
		}
	}
	switch {
	case len(hubs) == 0:
		return Placement{}, nil, ErrNoHub
	case len(hubs) > 1:
		return Placement{}, nil, ErrMultipleHubs
	}

	p := Placement{
		Hub:     hubs[0],
		Servers: make(map[string]astorg.Server, len(lookup.ServerByID)),
		Phones:  make(map[string]string),
		Hosts:   make(map[string]string),
		Homes:   make(map[string]string),
		Owners:  make(map[string]string),
		Devices: make(map[string][]string),
		Remote:  make(map[string]map[string][]string),

		LineMACs: make(map[string]string),
	}
	for id, server := range lookup.ServerByID {
		p.Servers[id] = server
	}

	var problems []Problem

	// Step 1: Host each desk phone on its location's server, and set aside
	// each spare
	spares := make(map[string]bool)
	for _, phone := range data.Phones {
		name := lineUsername(phone.MAC, lookup)
		if name == "" {
			continue
		}
		if phone.Location == "" {
			if phoneOwned(phone.MAC, lookup) {
				problems = append(problems, Problem{
					Rule:   RuleOwnedNoLocation,
					Kind:   "phone",
					ID:     phone.MAC,
					Detail: fmt.Sprintf("the phone has no location but is assigned to %q", name),
				})
				continue
			}
			if host, found := p.Hosts[name]; found {
				problems = append(problems, Problem{
					Rule:   RuleDeviceConflict,
					Kind:   "phone",
					ID:     phone.MAC,
					Detail: fmt.Sprintf("line %q is already hosted on server %q, so it cannot be a spare", name, host),
				})
				continue
			}
			// An unassigned phone's line is named by its MAC address.
			if !spares[name] {
				spares[name] = true
				p.Spares = append(p.Spares, phone.MAC)
				p.LineMACs[name] = phone.MAC
			}
			continue
		}
		id, _, ok := lookup.LocationServer(phone.Location)
		switch {
		case id == "":
			problems = append(problems, Problem{
				Rule:   RuleLocationNoServer,
				Kind:   "phone",
				ID:     phone.MAC,
				Detail: fmt.Sprintf("location %q has no server", phone.Location),
			})
			continue
		case !ok:
			problems = append(problems, Problem{
				Rule:   RuleServerMissing,
				Kind:   "phone",
				ID:     phone.MAC,
				Detail: fmt.Sprintf("location %q names server %q, which is not in the data set", phone.Location, id),
			})
			continue
		}
		if host, found := p.Hosts[name]; found && host != id {
			problems = append(problems, Problem{
				Rule:   RuleDeviceConflict,
				Kind:   "phone",
				ID:     phone.MAC,
				Detail: fmt.Sprintf("line %q is already hosted on server %q, not %q", name, host, id),
			})
			continue
		}
		if spares[name] {
			problems = append(problems, Problem{
				Rule:   RuleDeviceConflict,
				Kind:   "phone",
				ID:     phone.MAC,
				Detail: fmt.Sprintf("line %q is already the line of a spare desk phone on every server", name),
			})
			continue
		}
		p.Phones[phone.MAC] = id
		p.Hosts[name] = id
		if _, found := p.LineMACs[name]; !found {
			p.LineMACs[name] = phone.MAC
		}
	}
	sort.Strings(p.Spares)

	// Step 2: Host every softphone on the hub
	for _, softphone := range data.Softphones {
		if softphone.Username == "" {
			continue
		}
		if host, found := p.Hosts[softphone.Username]; found {
			problems = append(problems, Problem{
				Rule:   RuleDeviceConflict,
				Kind:   "softphone",
				ID:     softphone.Username,
				Detail: fmt.Sprintf("device %q is already hosted on server %q", softphone.Username, host),
			})
			continue
		}
		if spares[softphone.Username] {
			problems = append(problems, Problem{
				Rule:   RuleDeviceConflict,
				Kind:   "softphone",
				ID:     softphone.Username,
				Detail: fmt.Sprintf("device %q is already the line of a spare desk phone on every server", softphone.Username),
			})
			continue
		}
		p.Hosts[softphone.Username] = p.Hub
	}

	// Step 3: Home each extension, people first and then roles
	type owner struct {
		kind       string
		username   string
		extension  string
		location   string
		softphones []string
	}
	var owners []owner
	for _, person := range data.People {
		owners = append(owners, owner{"person", person.Username, person.Extension, person.Location, person.Softphones})
	}
	for _, role := range data.PhoneRoles {
		owners = append(owners, owner{"role", role.Username, role.Extension, role.Location, role.Softphones})
	}

	for _, o := range owners {
		if o.extension == "" || o.username == "" {
			continue
		}
		if holder, found := p.Owners[o.extension]; found {
			problems = append(problems, Problem{
				Rule:   RuleExtensionDuplicate,
				Kind:   o.kind,
				ID:     o.username,
				Detail: fmt.Sprintf("extension %s is already held by %q", o.extension, holder),
			})
			continue
		}
		home := p.Hub
		id, _, ok := lookup.LocationServer(o.location)
		switch {
		case id != "" && !ok:
			problems = append(problems, Problem{
				Rule:   RuleServerMissing,
				Kind:   o.kind,
				ID:     o.username,
				Detail: fmt.Sprintf("location %q names server %q, which is not in the data set", o.location, id),
			})
			continue
		case id != "":
			home = id
		}
		p.Homes[o.extension] = home
		p.Owners[o.extension] = o.username

		// The devices an extension rings are its desk phone line, which is
		// named after the username, and its softphones.
		var devices []string
		for _, name := range append([]string{o.username}, o.softphones...) {
			if _, hosted := p.Hosts[name]; hosted {
				devices = append(devices, name)
			}
		}
		devices = sortedUnique(devices)
		p.Devices[o.extension] = devices

		for _, name := range devices {
			host := p.Hosts[name]
			if host == home {
				continue
			}
			if p.Remote[o.extension] == nil {
				p.Remote[o.extension] = make(map[string][]string)
			}
			p.Remote[o.extension][host] = append(p.Remote[o.extension][host], name)
		}
	}

	// Step 4: Pair every two servers that have an address
	var ids []string
	for _, id := range p.ServerIDs() {
		if p.Servers[id].Address == "" {
			problems = append(problems, Problem{
				Rule:   RuleServerNoAddress,
				Kind:   "server",
				ID:     id,
				Detail: "the server has no address, so it gets no trunk to any other server",
			})
			continue
		}
		ids = append(ids, id)
	}
	for i := range ids {
		for j := i + 1; j < len(ids); j++ {
			p.TrunkPairs = append(p.TrunkPairs, TrunkPair{
				ID: PairID(ids[i], ids[j]),
				A:  ids[i],
				B:  ids[j],
			})
		}
	}

	return p, problems, nil
}

// ServerIDs returns the ID of every server, sorted.
func (p Placement) ServerIDs() []string {
	ids := make([]string, 0, len(p.Servers))
	for id := range p.Servers {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Peers returns every server that shares a trunk pair with serverID,
// sorted by ID. That is every other server, unless serverID or the other
// server has no address.
func (p Placement) Peers(serverID string) []astorg.Server {
	var peers []astorg.Server
	for _, pair := range p.TrunkPairs {
		switch serverID {
		case pair.A:
			peers = append(peers, p.Servers[pair.B])
		case pair.B:
			peers = append(peers, p.Servers[pair.A])
		}
	}
	sort.Slice(peers, func(i, j int) bool { return peers[i].ID < peers[j].ID })
	return peers
}

// HostedDevices returns the names of the devices that serverID hosts,
// sorted.
func (p Placement) HostedDevices(serverID string) []string {
	var names []string
	for name, host := range p.Hosts {
		if host == serverID {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// HomedExtensions returns the extensions that live on serverID, sorted.
func (p Placement) HomedExtensions(serverID string) []string {
	var extensions []string
	for extension, home := range p.Homes {
		if home == serverID {
			extensions = append(extensions, extension)
		}
	}
	sort.Strings(extensions)
	return extensions
}

// Credentials returns the credentials that serverID's node agent writes:
// one trunk credential per peer, sorted by pair ID, and then one phone
// credential for each desk phone line serverID hosts and each spare,
// sorted by MAC address. Every office server lists every spare, whether
// or not it has an address; the hub lists none when there are offices. This is the one list of derived secrets the node
// agent needs: it holds every auth section that serverID's rendered
// endpoints and trunks name, apart from softphones, whose secrets are not
// derived.
func (p Placement) Credentials(serverID string) []Credential {
	var creds []Credential
	for _, pair := range p.TrunkPairs {
		var peer string
		switch serverID {
		case pair.A:
			peer = pair.B
		case pair.B:
			peer = pair.A
		default:
			continue
		}
		creds = append(creds, Credential{
			Kind:    "trunk",
			ID:      pair.ID,
			Section: trunkAuthName(peer),
		})
	}
	if _, ok := p.Servers[serverID]; !ok {
		return creds
	}
	spares := p.sparesOn(serverID)
	var phones []Credential
	for line, mac := range p.LineMACs {
		if p.Hosts[line] != serverID && !spares[line] {
			continue
		}
		phones = append(phones, Credential{
			Kind:    "phone",
			ID:      mac,
			Section: authName(line),
		})
	}
	sort.Slice(phones, func(i, j int) bool { return phones[i].ID < phones[j].ID })
	return append(creds, phones...)
}

// sparesOn returns the spares serverID renders: every spare on an office
// server, and none on a server the placement does not know. The hub gets
// none when there are offices: it need not run DPMA, and a spare is
// provisioned where it is plugged in, at an office. In a deployment of one
// server, the hub is the office and renders every spare.
func (p Placement) sparesOn(serverID string) map[string]bool {
	if _, ok := p.Servers[serverID]; !ok || (serverID == p.Hub && len(p.Servers) > 1) {
		return map[string]bool{}
	}
	return p.spareSet()
}

// spareSet returns the set of device names of the spare desk phones. An
// unassigned phone's line is named by its MAC address.
func (p Placement) spareSet() map[string]bool {
	set := make(map[string]bool, len(p.Spares))
	for _, mac := range p.Spares {
		set[mac] = true
	}
	return set
}

// phoneOwned reports whether a person or role owns the phone with the given
// MAC address, by the data set's phone assignments.
func phoneOwned(mac string, lookup astorg.Lookup) bool {
	assignment, ok := lookup.PhoneAssignments[mac]
	return ok && assignment.Type != astorg.Unassigned
}

// PairID returns the trunk pair ID for two servers: their IDs, sorted,
// joined by "~".
func PairID(a, b string) string {
	if b < a {
		a, b = b, a
	}
	return a + "~" + b
}

func sortedUnique(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	sorted := append([]string(nil), values...)
	sort.Strings(sorted)
	out := sorted[:1]
	for _, v := range sorted[1:] {
		if v != out[len(out)-1] {
			out = append(out, v)
		}
	}
	return out
}
