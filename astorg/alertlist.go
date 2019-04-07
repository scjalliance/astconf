package astorg

// AlertList is a slice of alerts.
type AlertList []Alert

// ByName returns a map of alerts indexed by name.
func (alerts AlertList) ByName() map[string]Alert {
	lookup := make(map[string]Alert, len(alerts))
	for _, alert := range alerts {
		if alert.Name == "" {
			continue
		}
		lookup[alert.Name] = alert
	}
	return lookup
}
