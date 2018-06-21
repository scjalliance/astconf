package dpma

// Alert is a DPMA alert definition.
type Alert struct {
	Info     string `astconf:"alert_info"`
	RingType string `astconf:"ring_type"`
	RingTone string `astconf:"ring_tone"`
}
