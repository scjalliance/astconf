package astgen_test

import (
	"testing"

	"github.com/scjalliance/astconf/astgen"
)

func TestLines(t *testing.T) {
	lines := astgen.Lines(testData())

	var values []interface{}
	for i := range lines {
		values = append(values, &lines[i])
	}
	got := render(t, values...)

	want := "[fred]\n" +
		"type = line\n" +
		"line_label = 100 Fred Flintstone\n" +
		"\n" +
		"[lobby]\n" +
		"type = line\n" +
		"line_label = 200 Lobby\n" +
		"\n" +
		"[000fd3000003]\n" +
		"type = line\n" +
		"line_label = QRY-000fd3000003\n"
	if got != want {
		t.Errorf("Lines() rendered:\n%s\nwant:\n%s", got, want)
	}
}
