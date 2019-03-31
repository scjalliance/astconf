package astval

import (
	"fmt"
	"strings"
)

// Var is an asterisk variable value.
type Var struct {
	Name  string
	Value string
}

// NewVar returns a variable value with the given name and value.
func NewVar(name, value string) Var {
	return Var{Name: name, Value: value}
}

// String returns the variable as a string.
func (v Var) String() string {
	return v.Name + "=" + v.Value
}

// MarshalText marshals the variable as a string of utf-8 bytes.
func (v Var) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText parses an asterisk variable value from a string of utf-8 bytes.
func (v *Var) UnmarshalText(text []byte) error {
	parts := strings.SplitN(string(text), "=", 2)
	switch len(parts) {
	case 2:
		v.Name = parts[0]
		v.Value = parts[1]
		return nil
	case 1:
		v.Name = parts[0]
		return nil
	default:
		return fmt.Errorf("cannot unmarshal \"%s\" as a variable value", string(text))
	}
}
