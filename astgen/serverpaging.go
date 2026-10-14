package astgen

import (
	"fmt"
	"sort"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dialplan"
)

// ServerPagingExtensions generates the paging group extensions for
// serverID. A page reaches the members that serverID hosts directly, and
// reaches the members on each other server through one relay leg per
// server.
//
// Each relay leg is a local channel into PageRelayContext rather than the
// trunk itself. Page's "s" option skips a recipient that is not idle, and a
// trunk endpoint reads as in use while it carries any call, so a direct
// trunk leg could be skipped. A local channel is idle unless the same relay
// is already running, so "s" judges only the local member devices.
func ServerPagingExtensions(data *astorg.DataSet, p Placement, serverID string, context string) dialplan.Section {
	section := dialplan.Section{Context: context}
	members := pagingMembers(data)
	for _, group := range data.PagingGroups {
		local, peers := splitMembers(p, serverID, members[group.Extension])
		devices := local
		for _, peer := range peers {
			devices = append(devices, dialplan.Device{
				Technology: "Local",
				Resource:   pageRelayName(group.Extension, peer) + "@" + PageRelayContext,
			})
		}
		section.Extensions = append(section.Extensions, dialplan.Extension{
			Comment: group.Name,
			Number:  group.Extension,
			Actions: pageActions(group, devices),
		})
	}
	return section
}

// ServerPageRelayExtensions generates PageRelayContext for serverID: one
// extension for each paging group and each other server that hosts members
// of it, which carries the page over the trunk to that server's from-trunk
// page- entry.
func ServerPageRelayExtensions(data *astorg.DataSet, p Placement, serverID string) dialplan.Section {
	section := dialplan.Section{Context: PageRelayContext}
	members := pagingMembers(data)
	for _, group := range data.PagingGroups {
		_, peers := splitMembers(p, serverID, members[group.Extension])
		for _, peer := range peers {
			section.Extensions = append(section.Extensions, dialplan.Extension{
				Comment: fmt.Sprintf("Relay %s to %s", group.Name, peer),
				Number:  pageRelayName(group.Extension, peer),
				Actions: []dialplan.Action{
					dialplan.Dial(dialplan.PJSIP("page-"+group.Extension+"@"+trunkName(peer)), 0),
				},
			})
		}
	}
	return section
}

// PageAlertSection generates the subroutine in PageAlertContext that Page
// runs on each callee channel before dialing it. It adds the paging
// group's Alert-Info header, passed as the first argument, so the phone
// answers the page on its own. A relay leg is a local channel with no SIP
// headers, so the header is added only on a PJSIP channel; the receiving
// server adds it for its own members.
func PageAlertSection() dialplan.Section {
	return dialplan.Section{
		Context: PageAlertContext,
		Extensions: []dialplan.Extension{
			{
				Comment: "Add a paging group's Alert-Info header to a paged phone",
				Number:  "s",
				Actions: []dialplan.Action{
					dialplan.ExecIf(
						dialplan.Equal(dialplan.Channel("channeltype"), dialplan.String("PJSIP")),
						dialplan.Set(dialplan.PJSIPHeader("add", "Alert-Info"), dialplan.Var("ARG1")),
					),
					dialplan.Return(),
				},
			},
		},
	}
}

// pageRelayName returns the name of the relay extension that carries a
// page for the group with the given extension to peer.
func pageRelayName(group, peer string) string {
	return "page-relay-" + group + "-" + peer
}

// splitMembers divides the members of a paging group into the devices that
// serverID hosts, sorted, and the IDs of the other servers that host at
// least one member, sorted. A member that no server hosts is dropped.
func splitMembers(p Placement, serverID string, names []string) (local []dialplan.Device, peers []string) {
	seenPeer := make(map[string]bool)
	for _, name := range sortedUnique(names) {
		host, ok := p.Hosts[name]
		switch {
		case !ok:
			continue
		case host == serverID:
			local = append(local, dialplan.PJSIP(name))
		case !seenPeer[host]:
			seenPeer[host] = true
			peers = append(peers, host)
		}
	}
	sort.Strings(peers)
	return local, peers
}

// pageActions returns the actions of a paging group extension that pages
// the given devices. They follow PagingGroupExtensions, except that the
// Alert-Info header is added on each callee channel through
// PageAlertContext, because SIPAddHeader has no effect on a PJSIP channel.
func pageActions(group astorg.PagingGroup, devices []dialplan.Device) []dialplan.Action {
	actions := []dialplan.Action{
		dialplan.Noop(fmt.Sprintf("Paging %s", group.Name)),
		dialplan.Playback("beep"),
	}
	actions = append(actions, pageCallerIDActions()...)
	if len(devices) > 0 {
		page := dialplan.Page(devices...)
		if group.Alert != "" {
			page.Predial = &dialplan.PredialTarget{
				Context:   PageAlertContext,
				Extension: "s",
				Priority:  1,
				Args:      []string{group.Alert},
			}
		}
		actions = append(actions, page)
	}
	actions = append(actions, dialplan.Hangup())
	return actions
}
