package dialplan

import "strings"

// Hint is the set of devices whose combined state an extension's hint
// reports. A phone that watches the extension, such as a busy-lamp key,
// lights from that state.
type Hint []Device

// String returns the hint's devices joined by &.
func (h Hint) String() string {
	parts := make([]string, len(h))
	for i, device := range h {
		parts[i] = device.String()
	}
	return strings.Join(parts, "&")
}

// Validate returns an error if any device would change the meaning of the
// hint.
func (h Hint) Validate() error {
	for _, device := range h {
		if err := device.Validate(); err != nil {
			return err
		}
	}
	return nil
}
