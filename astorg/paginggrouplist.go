package astorg

// PagingGroupList is a slice of paging groups.
type PagingGroupList []PagingGroup

// ByExtension returns a map of paging groups indexed by extension.
func (groups PagingGroupList) ByExtension() map[string]PagingGroup {
	lookup := make(map[string]PagingGroup, len(groups))
	for _, group := range groups {
		if group.Extension == "" {
			continue
		}
		lookup[group.Extension] = group
	}
	return lookup
}
