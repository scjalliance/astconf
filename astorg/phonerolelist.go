package astorg

// PhoneRoleList is a slice of phone roles.
type PhoneRoleList []PhoneRole

// ByUsername returns a map of phone roles indexed by username.
func (roles PhoneRoleList) ByUsername() map[string]PhoneRole {
	lookup := make(map[string]PhoneRole, len(roles))
	for _, role := range roles {
		if role.Username == "" {
			continue
		}
		lookup[role.Username] = role
	}
	return lookup
}

// ByExtension returns a map of phone roles indexed by extension.
func (roles PhoneRoleList) ByExtension() map[string]PhoneRole {
	lookup := make(map[string]PhoneRole, len(roles))
	for _, role := range roles {
		if role.Extension == "" {
			continue
		}
		lookup[role.Extension] = role
	}
	return lookup
}

// Assignments returns a phone assignment map for roles, indexed by
// mac address.
func (roles PhoneRoleList) Assignments() AssignmentMap {
	assignments := make(AssignmentMap, len(roles))
	for _, role := range roles {
		for _, mac := range role.Phones {
			assignments.Add(mac, Assignment{Type: RoleAssigned, Username: role.Username})
		}
	}
	return assignments
}
