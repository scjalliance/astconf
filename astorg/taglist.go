package astorg

// TagList is a slice of tags.
type TagList []Tag

// ByName returns a map of tags indexed by name.
func (tags TagList) ByName() map[string]Tag {
	lookup := make(map[string]Tag, len(tags))
	for _, tag := range tags {
		if tag.Name == "" {
			continue
		}
		lookup[tag.Name] = tag
	}
	return lookup
}
