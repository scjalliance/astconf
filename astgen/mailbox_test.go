package astgen_test

import (
	"testing"

	"github.com/scjalliance/astconf/astgen"
)

func TestMailboxes(t *testing.T) {
	section := astgen.Mailboxes(testData(), "default")
	got := render(t, &section)

	want := "[default]\n" +
		"300 => -3000,Sales,sales@example.com,,attach=yes|saycid=yes|delete=yes\n" +
		"100 => -1234,fred,fred@example.com,,attach=yes|saycid=yes|delete=no\n"
	if got != want {
		t.Errorf("Mailboxes() rendered:\n%s\nwant:\n%s", got, want)
	}
}
