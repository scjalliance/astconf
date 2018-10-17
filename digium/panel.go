package digium

// Panel describes the properties of a display panel on a model
// of Digium phone.
type Panel struct {
	Width  int  // Display width in pixels
	Height int  // Display height in pixels
	Color  bool // Does it have a color display?
	Fields int  // The number of busy lamp fields present on the display
}

// Digium display panel names used for Smart BLF instructions.
const (
	MainPanel = "main"
	SidePanel = "side"
)
