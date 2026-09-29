package astconf

import "testing"

// TestUnescapedLen checks the mapping from bytes of escaped output back to
// whole bytes of the original value, including a stop between a backslash
// and the semicolon it escapes.
func TestUnescapedLen(t *testing.T) {
	tests := []struct {
		p       string
		written int
		want    int
	}{
		{p: "abc", written: 0, want: 0},
		{p: "abc", written: 2, want: 2},
		{p: "abc", written: 3, want: 3},
		{p: "a;", written: 1, want: 1},  // "a"
		{p: "a;", written: 2, want: 1},  // "a\" stops inside the escape
		{p: "a;", written: 3, want: 2},  // "a\;"
		{p: ";;b", written: 1, want: 0}, // "\"
		{p: ";;b", written: 3, want: 1}, // "\;\"
		{p: ";;b", written: 4, want: 2}, // "\;\;"
		{p: ";;b", written: 5, want: 3}, // "\;\;b"
	}
	for _, tt := range tests {
		if got := unescapedLen([]byte(tt.p), tt.written); got != tt.want {
			t.Errorf("unescapedLen(%q, %d) = %d, want %d", tt.p, tt.written, got, tt.want)
		}
	}
}
