package digium

// Models is a lookup of Digium phone models by name.
var Models = map[string]Model{
	"D80": D80,
	"D70": D70,
	"D65": D65,
	"D62": D62,
}

// D80 holds the specifications of the D80 model of phone.
var D80 = Model{
	Name:      "D80",
	MainPanel: Panel{Width: 800, Height: 1280, Fields: 20, Color: true},
}

// D70 holds the specifications of the D70 model of phone.
var D70 = Model{
	Name:      "D70",
	MainPanel: Panel{Width: 320, Height: 160, Fields: 6},
	SidePanel: Panel{Width: 113, Height: 357, Fields: 10},
}

// D65 holds the specifications of the D65 model of phone.
var D65 = Model{
	Name:      "D65",
	MainPanel: Panel{Width: 480, Height: 272, Fields: 5, Color: true},
}

// D62 holds the specifications of the D62 model of phone.
var D62 = Model{
	Name:      "D62",
	MainPanel: Panel{Width: 480, Height: 272, Fields: 1, Color: true},
}

// D60 holds the specifications of the D60 model of phone.
var D60 = Model{
	Name:      "D60",
	MainPanel: Panel{Width: 480, Height: 272, Fields: 1, Color: true},
}
