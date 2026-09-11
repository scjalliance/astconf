// Package pjsip provides types for writing chan_pjsip configuration
// sections to pjsip.conf.
//
// Each PJSIP object type (endpoint, auth, aor, transport) is a separate
// Go type. Unlike chan_sip, a single device is described by an endpoint
// section, an auth section and an aor section that reference each other
// by name. The astgen package generates all three for each phone.
//
// The section name, templates and type preamble are provided by pointer
// receiver methods, so pass a pointer to the encoder, as in
// e.Encode(&endpoint). Encoding a value instead of a pointer writes the
// settings without the section header. This matches the sip and dpma
// packages.
//
// Only the settings needed by SCJ configurations are modeled.
package pjsip
