// Package pjsip provides types for writing chan_pjsip configuration
// sections to pjsip.conf.
//
// Each PJSIP object type (endpoint, auth, aor, transport, identify,
// outbound-publish, inbound-publication, asterisk-publication) is a
// separate Go type. Unlike chan_sip, a single device is described by an
// endpoint section, an auth section and an aor section that reference each
// other by name. The astgen package generates all three for each phone.
//
// A trunk to another Asterisk server adds an identify section, which
// recognizes the server by its address, and the three publication types,
// which share device and mailbox state with it through
// res_pjsip_publish_asterisk.
//
// The section name, templates and type preamble are provided by pointer
// receiver methods, so pass a pointer to the encoder, as in
// e.Encode(&endpoint). Encoding a value instead of a pointer writes the
// settings without the section header. This matches the sip and dpma
// packages.
//
// Only a subset of the PJSIP settings is modeled. Fields are added as
// configurations need them.
package pjsip
