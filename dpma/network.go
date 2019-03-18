package dpma

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astoverlay"
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
		astoverlay.SectionName(&network.Name, &merged.Name)
		astoverlay.String(&network.Alias, &merged.Alias)
		astoverlay.String(&network.CIDR, &merged.CIDR)
		astoverlay.String(&network.RegistrationAddress, &merged.RegistrationAddress)
		astoverlay.Int(&network.RegistrationPort, &merged.RegistrationPort)
		astoverlay.String(&network.Transport, &merged.Transport)
		astoverlay.String(&network.AlternateRegistrationAddress, &merged.AlternateRegistrationAddress)
		astoverlay.Int(&network.AlternateRegistrationPort, &merged.AlternateRegistrationPort)
		astoverlay.String(&network.AlternateTransport, &merged.AlternateTransport)
		astoverlay.String(&network.FileURLPrefix, &merged.FileURLPrefix)
		astoverlay.String(&network.PublicFirmwareURLPrefix, &merged.PublicFirmwareURLPrefix)
		astoverlay.String(&network.NTPServer, &merged.NTPServer)
		astoverlay.String(&network.SyslogServer, &merged.SyslogServer)
		astoverlay.Int(&network.SyslogPort, &merged.SyslogPort)
		astoverlay.String(&network.SyslogLevel, &merged.SyslogLevel)
		astoverlay.String(&network.NetworkVLANDiscoveryMode, &merged.NetworkVLANDiscoveryMode)
		astoverlay.AstInt(&network.SIPQOS, &merged.SIPQOS)
		astoverlay.AstInt(&network.RTPQOS, &merged.RTPQOS)
		astoverlay.AstInt(&network.NetworkVLANID, &merged.NetworkVLANID)
		astoverlay.AstInt(&network.PCVLANID, &merged.PCVLANID)
		astoverlay.AstInt(&network.PCQOS, &merged.PCQOS)
		astoverlay.AstInt(&network.SIPDSCP, &merged.SIPDSCP)
		astoverlay.AstInt(&network.RTPDSCP, &merged.RTPDSCP)
		astoverlay.AstSeconds(&network.UDPKAInterval, &merged.UDPKAInterval)
	}
	return
}
