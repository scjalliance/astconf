// Package astval provides a set of value types used to serialize
// configuration data for asterisk phone servers.
//
// Many of the types in this package have an "unspecified" zero value
// that will not be serialized. This assists with creation of config
// structures that preserve default values for settings that have not
// been explicitly set.
package astval
