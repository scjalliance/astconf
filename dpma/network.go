package dpma

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astmerge"
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
		astmerge.SectionName(&network.Name, &merged.Name)
		astmerge.String(&network.Alias, &merged.Alias)
		astmerge.String(&network.CIDR, &merged.CIDR)
		astmerge.String(&network.RegistrationAddress, &merged.RegistrationAddress)
		astmerge.Int(&network.RegistrationPort, &merged.RegistrationPort)
		astmerge.String(&network.Transport, &merged.Transport)
		astmerge.String(&network.AlternateRegistrationAddress, &merged.AlternateRegistrationAddress)
		astmerge.Int(&network.AlternateRegistrationPort, &merged.AlternateRegistrationPort)
		astmerge.String(&network.AlternateTransport, &merged.AlternateTransport)
		astmerge.String(&network.FileURLPrefix, &merged.FileURLPrefix)
		astmerge.String(&network.PublicFirmwareURLPrefix, &merged.PublicFirmwareURLPrefix)
		astmerge.String(&network.NTPServer, &merged.NTPServer)
		astmerge.String(&network.SyslogServer, &merged.SyslogServer)
		astmerge.Int(&network.SyslogPort, &merged.SyslogPort)
		astmerge.String(&network.SyslogLevel, &merged.SyslogLevel)
		astmerge.String(&network.NetworkVLANDiscoveryMode, &merged.NetworkVLANDiscoveryMode)
		astmerge.AstInt(&network.SIPQOS, &merged.SIPQOS)
		astmerge.AstInt(&network.RTPQOS, &merged.RTPQOS)
		astmerge.AstInt(&network.NetworkVLANID, &merged.NetworkVLANID)
		astmerge.AstInt(&network.PCVLANID, &merged.PCVLANID)
		astmerge.AstInt(&network.PCQOS, &merged.PCQOS)
		astmerge.AstInt(&network.SIPDSCP, &merged.SIPDSCP)
		astmerge.AstInt(&network.RTPDSCP, &merged.RTPDSCP)
		astmerge.AstSeconds(&network.UDPKAInterval, &merged.UDPKAInterval)
	}
	return
}
