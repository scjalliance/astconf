// Package pjsip provides types for writing chan_pjsip configuration
// sections to pjsip.conf.
//
// Each PJSIP object type (endpoint, auth, aor, transport) is a separate
// Go type. Unlike chan_sip, a single device is described by an endpoint
// section, an auth section and an aor section that reference each other
// by name. The astgen package generates all three for each phone.
//
// Only the settings needed by SCJ configurations are modeled.
package pjsip
