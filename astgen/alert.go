package astgen

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dpma"
)

// Alerts generates DPMA alert entries for a dataset.
func Alerts(data *astorg.DataSet) []dpma.Alert {
	var alerts []dpma.Alert
	for _, alert := range data.Alerts {
		alerts = append(alerts, dpma.Alert{
			Name:       alert.Name,
			InfoHeader: alert.Name,
			RingType:   alert.Type,
			Ringtone:   alert.Ringtone,
		})
	}
	return alerts
}
