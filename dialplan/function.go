package dialplan

// Function is an interface that can be evaluated by the dialplan
// as a function.
type Function interface {
	Func() FuncCall
}
