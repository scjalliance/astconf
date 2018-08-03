package dpma

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
)

// Network is a DPMA network definition.
type Network struct {
	Name                         astconf.SectionName
	Alias                        string         `astconf:"alias"`
	CIDR                         string         `astconf:"cidr"`
	RegistrationAddress          string         `astconf:"registration_address"`
	RegistrationPort             int            `astconf:"registration_port,omitempty"`
	Transport                    string         `astconf:"transport"` // "udp", "tcp", "tls"
	AlternateRegistrationAddress string         `astconf:"alternate_registration_address"`
	AlternateRegistrationPort    int            `astconf:"alternate_registration_port"`
	AlternateTransport           string         `astconf:"alternate_transport"`
	FileURLPrefix                string         `astconf:"file_url_prefix"`
	PublicFirmwareURLPrefix      string         `astconf:"public_firmware_url_prefix"`
	NTPServer                    string         `astconf:"ntp_server"`
	SyslogServer                 string         `astconf:"syslog_server"`
	SyslogPort                   int            `astconf:"syslog_port,omitempty"`
	SyslogLevel                  string         `astconf:"syslog_level"`                // "debug", "error", "warning", "information"
	NetworkVLANDiscoveryMode     string         `astconf:"network_vlan_discovery_mode"` // "NONE", "MANUAL", "LLDP"
	SIPQOS                       astval.Int     `astconf:"sip_qos"`
	RTPQOS                       astval.Int     `astconf:"rtp_qos"`
	NetworkVLANID                astval.Int     `astconf:"network_vlan_id"`
	PCVLANID                     astval.Int     `astconf:"pc_vlan_id"`
	PCQOS                        astval.Int     `astconf:"pc_qos"`
	SIPDSCP                      astval.Int     `astconf:"sip_dscp"`
	RTPDSCP                      astval.Int     `astconf:"rtp_dscp"`
	UDPKAInterval                astval.Seconds `astconf:"udp_ka_interval"` // Seconds
}

// MergeNetworks returns the merged configuration of all the given networks,
// in order of priority from least to greatest.
func MergeNetworks(networks ...*Network) (merged Network) {
	for _, network := range networks {
		mergeSectionName(&network.Name, &merged.Name)
		mergeString(&network.Alias, &merged.Alias)
		mergeString(&network.CIDR, &merged.CIDR)
		mergeString(&network.RegistrationAddress, &merged.RegistrationAddress)
		mergeInt(&network.RegistrationPort, &merged.RegistrationPort)
		mergeString(&network.Transport, &merged.Transport)
		mergeString(&network.AlternateRegistrationAddress, &merged.AlternateRegistrationAddress)
		mergeInt(&network.AlternateRegistrationPort, &merged.AlternateRegistrationPort)
		mergeString(&network.AlternateTransport, &merged.AlternateTransport)
		mergeString(&network.FileURLPrefix, &merged.FileURLPrefix)
		mergeString(&network.PublicFirmwareURLPrefix, &merged.PublicFirmwareURLPrefix)
		mergeString(&network.NTPServer, &merged.NTPServer)
		mergeString(&network.SyslogServer, &merged.SyslogServer)
		mergeInt(&network.SyslogPort, &merged.SyslogPort)
		mergeString(&network.SyslogLevel, &merged.SyslogLevel)
		mergeString(&network.NetworkVLANDiscoveryMode, &merged.NetworkVLANDiscoveryMode)
		mergeAstInt(&network.SIPQOS, &merged.SIPQOS)
		mergeAstInt(&network.RTPQOS, &merged.RTPQOS)
		mergeAstInt(&network.NetworkVLANID, &merged.NetworkVLANID)
		mergeAstInt(&network.PCVLANID, &merged.PCVLANID)
		mergeAstInt(&network.PCQOS, &merged.PCQOS)
		mergeAstInt(&network.SIPDSCP, &merged.SIPDSCP)
		mergeAstInt(&network.RTPDSCP, &merged.RTPDSCP)
		mergeAstSeconds(&network.UDPKAInterval, &merged.UDPKAInterval)
	}
	return
}
