package astconf

import "testing"

func TestTagParsing(t *testing.T) {
	tag := parseTag("field,world,oyster")
	if tag.name != "field" {
		t.Fatalf("name = %q, want field", tag.name)
	}
	for _, tt := range []struct {
		opt  string
		want bool
	}{
		{"world", true},
		{"oyster", true},
		{"baz", false},
	} {
		if tag.opts.Contains(tt.opt) != tt.want {
			t.Errorf("Contains(%q) = %v", tt.opt, !tt.want)
		}
	}
}
