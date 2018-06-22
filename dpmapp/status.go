package dpmaapp

// Status is a DPMA status application definition.
type Status struct {
	Status    string `astconf:"status"` // "available", "dnd", "away", "xa", "chat", "unavailable"
	SubStatus string `astconf:"substatus,omitempty"`
	Send486   bool   `astconf:"send486,omitempty"`
}
