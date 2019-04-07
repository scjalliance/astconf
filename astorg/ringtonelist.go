package astorg

// RingtoneList is a slice of ringtones.
type RingtoneList []Ringtone

// ByName returns a map of ringtones indexed by name.
func (ringtones RingtoneList) ByName() map[string]Ringtone {
	lookup := make(map[string]Ringtone, len(ringtones))
	for _, ringtone := range ringtones {
		if ringtone.Name == "" {
			continue
		}
		lookup[ringtone.Name] = ringtone
	}
	return lookup
}
