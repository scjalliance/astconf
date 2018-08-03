package dpmastatus

import "github.com/scjalliance/astconf/astval"

// Config holds configuration for the DPMA status application.
type Config struct {
	Status    string       `astconf:"status"` // "available", "dnd", "away", "xa", "chat", "unavailable"
	SubStatus string       `astconf:"substatus,omitempty"`
	Send486   astval.YesNo `astconf:"send486,omitempty"`
}
