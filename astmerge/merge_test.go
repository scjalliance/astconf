package astmerge_test

import (
	"reflect"
	"testing"

	"github.com/scjalliance/astconf/astmerge"
	"github.com/scjalliance/astconf/astval"
)

func TestStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		from, to []string
		want     []string
	}{
		{name: "both nil", want: nil},
		{name: "nil from", to: []string{"a"}, want: []string{"a"}},
		{name: "nil to", from: []string{"a"}, want: []string{"a"}},
		{name: "from values come first", from: []string{"a"}, to: []string{"b"}, want: []string{"a", "b"}},
		{name: "duplicates removed", from: []string{"a", "b"}, to: []string{"b", "c", "a"}, want: []string{"a", "b", "c"}},
		{name: "duplicates within from removed", from: []string{"a", "a"}, want: []string{"a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astmerge.StringSlice(&tt.from, &tt.to)
			if !reflect.DeepEqual(tt.to, tt.want) {
				t.Errorf("to = %v, want %v", tt.to, tt.want)
			}
		})
	}
}

func TestAstVarSlice(t *testing.T) {
	a := astval.NewVar("A", "1")
	b := astval.NewVar("B", "2")
	a2 := astval.NewVar("A", "2") // Same name, different value, is a distinct entry
	tests := []struct {
		name     string
		from, to []astval.Var
		want     []astval.Var
	}{
		{name: "both nil", want: nil},
		{name: "from values come first", from: []astval.Var{a}, to: []astval.Var{b}, want: []astval.Var{a, b}},
		{name: "duplicates removed", from: []astval.Var{a, b}, to: []astval.Var{b, a}, want: []astval.Var{a, b}},
		{name: "same name different value kept", from: []astval.Var{a}, to: []astval.Var{a2}, want: []astval.Var{a, a2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astmerge.AstVarSlice(&tt.from, &tt.to)
			if !reflect.DeepEqual(tt.to, tt.want) {
				t.Errorf("to = %v, want %v", tt.to, tt.want)
			}
		})
	}
}
