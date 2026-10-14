package astgen

import (
	"fmt"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dialplan"
)

// PagingGroupExtensions generates a dialplan section that includes all
// paging group extensions in a dataset.
func PagingGroupExtensions(data *astorg.DataSet, context string) dialplan.Section {
	// Build up the set of devices contained in each paging group
	members := pagingMembers(data)

	// Build an extension for each paging group
	section := dialplan.Section{Context: context}

	for _, group := range data.PagingGroups {
		var devices []dialplan.Device
		for _, name := range members[group.Extension] {
			devices = append(devices, dialplan.SIP(name))
		}
		devices = dedupAndSortDevices(devices)

		actions := []dialplan.Action{
			dialplan.Noop(fmt.Sprintf("Paging %s", group.Name)),
			dialplan.Playback("beep"),
		}
		if group.Alert != "" {
			actions = append(actions, dialplan.SIPAddHeader("Alert-Info", group.Alert))
		}
		actions = append(actions, pageCallerIDActions()...)
		if len(devices) > 0 {
			actions = append(actions, dialplan.Page(devices...))
		}
		actions = append(actions, dialplan.Hangup())

		ext := dialplan.Extension{
			Comment: group.Name,
			Number:  group.Extension,
			Actions: actions,
		}

		section.Extensions = append(section.Extensions, ext)
	}

	return section
}
