package dialplan

// Application is an interface that can be evaluated by the dialplan
// as an application.
type Application interface {
	App() AppCall
}
