package digium

// Model specifies a model of digium phone.
type Model struct {
	Name      string
	MainPanel Panel
	SidePanel Panel
}

// Panel returns the panel with the given name. If the model doesn't
// have such a panel an empty panel will be returned.
func (m Model) Panel(name string) Panel {
	switch name {
	case "main":
		return m.MainPanel
	case "side":
		return m.SidePanel
	default:
		return Panel{}
	}
}
