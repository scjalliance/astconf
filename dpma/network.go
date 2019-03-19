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

// OverlayNetworks returns the overlayed configuration of all the given
// networks, in order of priority from least to greatest.
func OverlayNetworks(networks ...*Network) (overlayed Network) {
	for _, network := range networks {
		astoverlay.SectionName(&network.Name, &overlayed.Name)
		astoverlay.String(&network.Alias, &overlayed.Alias)
		astoverlay.String(&network.CIDR, &overlayed.CIDR)
		astoverlay.String(&network.RegistrationAddress, &overlayed.RegistrationAddress)
		astoverlay.Int(&network.RegistrationPort, &overlayed.RegistrationPort)
		astoverlay.String(&network.Transport, &overlayed.Transport)
		astoverlay.String(&network.AlternateRegistrationAddress, &overlayed.AlternateRegistrationAddress)
		astoverlay.Int(&network.AlternateRegistrationPort, &overlayed.AlternateRegistrationPort)
		astoverlay.String(&network.AlternateTransport, &overlayed.AlternateTransport)
		astoverlay.String(&network.FileURLPrefix, &overlayed.FileURLPrefix)
		astoverlay.String(&network.PublicFirmwareURLPrefix, &overlayed.PublicFirmwareURLPrefix)
		astoverlay.String(&network.NTPServer, &overlayed.NTPServer)
		astoverlay.String(&network.SyslogServer, &overlayed.SyslogServer)
		astoverlay.Int(&network.SyslogPort, &overlayed.SyslogPort)
		astoverlay.String(&network.SyslogLevel, &overlayed.SyslogLevel)
		astoverlay.String(&network.NetworkVLANDiscoveryMode, &overlayed.NetworkVLANDiscoveryMode)
		astoverlay.AstInt(&network.SIPQOS, &overlayed.SIPQOS)
		astoverlay.AstInt(&network.RTPQOS, &overlayed.RTPQOS)
		astoverlay.AstInt(&network.NetworkVLANID, &overlayed.NetworkVLANID)
		astoverlay.AstInt(&network.PCVLANID, &overlayed.PCVLANID)
		astoverlay.AstInt(&network.PCQOS, &overlayed.PCQOS)
		astoverlay.AstInt(&network.SIPDSCP, &overlayed.SIPDSCP)
		astoverlay.AstInt(&network.RTPDSCP, &overlayed.RTPDSCP)
		astoverlay.AstSeconds(&network.UDPKAInterval, &overlayed.UDPKAInterval)
	}
	return
}
