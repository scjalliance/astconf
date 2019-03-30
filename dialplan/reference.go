package dialplan

// Reference is an interface that can be evaluated by the dialplan
// as a reference.
type Reference interface {
	Ref() string
}
