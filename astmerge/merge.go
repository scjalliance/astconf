package astmerge

import "github.com/scjalliance/astconf/astval"

// StringSlice merges two string slices, combining and deduplicating
// their values.
func StringSlice(from, to *[]string) {
	*to = dedupStringSlice(append(*from, (*to)...))
}

// AstVarSlice merges two asterisk variable value slices, combining and
// deduplicating their values.
func AstVarSlice(from, to *[]astval.Var) {
	*to = dedupAstVarSlice(append(*from, (*to)...))
}
