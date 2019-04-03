package astgen

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dpma"
)

// Ringtones generates DPMA ringtone entries for a dataset.
func Ringtones(data *astorg.DataSet) []dpma.Ringtone {
	var ringtones []dpma.Ringtone
	for _, ringtone := range data.Ringtones {
		ringtones = append(ringtones, dpma.Ringtone{
			Name:     ringtone.Name,
			Alias:    ringtone.DisplayName,
			Filename: ringtone.Filename,
		})
	}
	return ringtones
}
