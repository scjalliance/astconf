package astgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dialplan"
)

// dialVar is the channel variable that a home extension builds its dial
// string in.
const dialVar = dialplan.Var("ASTGEN_DIAL")

// extensionOwner is a person or role that holds an extension.
type extensionOwner struct {
	Extension string
	Username  string
	Comment   string
	Mailbox   string // Empty when the owner has no mailbox
	Role      bool
}

// extensionOwners returns every person and then every role that holds an
// extension in p, in data set order.
func extensionOwners(data *astorg.DataSet, p Placement) []extensionOwner {
	var owners []extensionOwner
	for _, person := range data.People {
		if p.Owners[person.Extension] != person.Username || person.Username == "" {
			continue
		}
		owner := extensionOwner{
			Extension: person.Extension,
			Username:  person.Username,
			Comment:   person.FullName,
		}
		if person.VoicemailCode != "" {
			owner.Mailbox = person.Extension
		}
		owners = append(owners, owner)
	}
	for _, role := range data.PhoneRoles {
		if p.Owners[role.Extension] != role.Username || role.Username == "" {
			continue
		}
		owners = append(owners, extensionOwner{
			Extension: role.Extension,
			Username:  role.Username,
			Comment:   role.DisplayName,
			Mailbox:   role.MailboxNumber,
			Role:      true,
		})
	}
	return owners
}

// remoteDevice returns the dial device that rings a device hosted on
// another server through the trunk to that server.
func remoteDevice(device, host string) string {
	return "PJSIP/dev-" + device + "@" + trunkName(host)
}

// homeExtension builds the extension that rings every device of an owner
// whose extension lives on this server.
//
// Each local device joins the ring only when it is idle. Each remote device
// always joins, and its own server answers busy for it when it is in use.
//
// The pre-answer and no-answer subroutines are hand-written on each server
// and called with Gosub, because Asterisk 21 removed Macro. Each reads its
// arguments as ${ARG1} and ${ARG2} and must end in Return.
func homeExtension(p Placement, owner extensionOwner) dialplan.Extension {
	home := p.Homes[owner.Extension]

	var remote strings.Builder
	hosts := make([]string, 0, len(p.Remote[owner.Extension]))
	for host := range p.Remote[owner.Extension] {
		hosts = append(hosts, host)
	}
	sort.Strings(hosts)
	for _, host := range hosts {
		for _, device := range p.Remote[owner.Extension][host] {
			remote.WriteString("&" + remoteDevice(device, host))
		}
	}

	actions := []dialplan.Action{
		dialplan.Noop(fmt.Sprintf("Call %s", owner.Username)),
		dialplan.Subroutine("pre-answer", "s", 1, dialplan.Var("EXTEN")),
		dialplan.Set(dialVar, dialplan.String(remote.String())),
	}
	for _, device := range p.Devices[owner.Extension] {
		if p.Hosts[device] != home {
			continue
		}
		local := dialplan.PJSIP(device)
		actions = append(actions, dialplan.ExecIf(
			dialplan.Equal(dialplan.DeviceState(local), dialplan.String("NOT_INUSE")),
			dialplan.Set(dialVar, dialplan.Cat(dialVar, dialplan.String("&"+local.String()))),
		))
	}
	actions = append(actions, dialplan.ExecIf(
		dialplan.GreaterThan(dialplan.Len(dialVar), dialplan.Int(0)),
		dialplan.DialExpr(dialplan.Substring(dialVar, 1), 20),
	))

	switch {
	case owner.Mailbox != "":
		actions = append(actions, dialplan.Subroutine("no-answer", "s", 1, dialplan.Var("EXTEN"), dialplan.String(voicemailTarget(p, home, owner.Mailbox))))
	case owner.Role:
		actions = append(actions, dialplan.Congestion())
	default:
		actions = append(actions, dialplan.Subroutine("no-answer", "s", 1, dialplan.Var("EXTEN")))
	}

	return dialplan.Extension{
		Comment: owner.Comment,
		Number:  owner.Extension,
		Actions: actions,
	}
}

// voicemailTarget returns the dial target that the no-answer subroutine on
// home sends a caller to for mailbox. Voicemail lives on the hub, so the hub
// reaches its own from-trunk entry through a local channel and every other
// server reaches it over the trunk.
func voicemailTarget(p Placement, home, mailbox string) string {
	if home == p.Hub {
		return "Local/vm-" + mailbox + "@" + TrunkContext
	}
	return "PJSIP/vm-" + mailbox + "@" + trunkName(p.Hub)
}

// ServerExtensions generates the phones' dialplan section for serverID.
//
// It holds a home extension for each extension that lives on serverID and
// a route over the trunk for each extension that lives elsewhere, in data
// set order. An office server also sends every unknown extension to the
// hub. On the hub an unknown extension is invalid.
func ServerExtensions(data *astorg.DataSet, p Placement, serverID string, context string) dialplan.Section {
	section := dialplan.Section{Context: context}

	for _, owner := range extensionOwners(data, p) {
		home := p.Homes[owner.Extension]
		if home == serverID {
			section.Extensions = append(section.Extensions, homeExtension(p, owner))
			continue
		}
		section.Extensions = append(section.Extensions, dialplan.Extension{
			Comment: owner.Comment,
			Number:  owner.Extension,
			Actions: []dialplan.Action{
				dialplan.Dial(dialplan.PJSIP(owner.Extension+"@"+trunkName(home)), 0),
			},
		})
	}

	if serverID != p.Hub {
		section.Extensions = append(section.Extensions, dialplan.Extension{
			Comment: "Unknown extensions go to the hub",
			Number:  "_X.",
			Actions: []dialplan.Action{
				dialplan.DialExpr(dialplan.Cat(dialplan.String("PJSIP/"), dialplan.Var("EXTEN"), dialplan.String("@"+trunkName(p.Hub))), 0),
			},
		})
	}

	return section
}

// ServerTrunkExtensions generates the from-trunk dialplan section for
// serverID, which every call from another server enters.
//
// It holds the home extensions of serverID, a dev- entry for each device
// serverID hosts, a vm- entry for each mailbox on the hub, and a page- entry
// for each paging group with members on serverID. It never holds a route to
// another server, so a call cannot bounce between servers.
//
// vmContext is the voicemail context that the hub's mailboxes live in.
func ServerTrunkExtensions(data *astorg.DataSet, p Placement, serverID string, vmContext string) dialplan.Section {
	section := dialplan.Section{Context: TrunkContext}

	// Step 1: Home extensions
	for _, owner := range extensionOwners(data, p) {
		if p.Homes[owner.Extension] == serverID {
			section.Extensions = append(section.Extensions, homeExtension(p, owner))
		}
	}

	// Step 2: One leg of another server's ring for each hosted device
	for _, name := range p.HostedDevices(serverID) {
		device := dialplan.PJSIP(name)
		section.Extensions = append(section.Extensions, dialplan.Extension{
			Comment: fmt.Sprintf("Ring %s for another server", name),
			Number:  "dev-" + name,
			Actions: []dialplan.Action{
				dialplan.ExecIf(dialplan.NotEqual(dialplan.DeviceState(device), dialplan.String("NOT_INUSE")), dialplan.Busy()),
				dialplan.Dial(device, 0),
				dialplan.Hangup(),
			},
		})
	}

	// Step 3: Voicemail, on the hub alone
	if serverID == p.Hub {
		seen := make(map[string]bool)
		for _, box := range Mailboxes(data, vmContext).Mailboxes {
			if seen[box.Extension] {
				continue
			}
			seen[box.Extension] = true
			section.Extensions = append(section.Extensions, dialplan.Extension{
				Comment: fmt.Sprintf("Voicemail for %s", box.Name),
				Number:  "vm-" + box.Extension,
				Actions: []dialplan.Action{
					dialplan.VoiceMail(box.Extension, vmContext, "u"),
					dialplan.Hangup(),
				},
			})
		}
	}

	// Step 4: Pages relayed from another server, to local members only
	members := pagingMembers(data)
	for _, group := range data.PagingGroups {
		local, _ := splitMembers(p, serverID, members[group.Extension])
		if len(local) == 0 {
			continue
		}
		section.Extensions = append(section.Extensions, dialplan.Extension{
			Comment: group.Name,
			Number:  "page-" + group.Extension,
			Actions: pageActions(group, local),
		})
	}

	return section
}
