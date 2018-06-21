package dpma

import (
	"time"

	"github.com/scjalliance/astconf"
)

// Network is a DPMA line definition.
type Network struct {
	Name                         astconf.SectionName
	Alias                        string        `astconf:"alias"`
	CIDR                         string        `astconf:"cidr"`
	RegistrationAddress          string        `astconf:"registration_address"`
	RegistrationPort             int           `astconf:"registration_port"`
	Transport                    string        `astconf:"transport"` // "udp", "tcp", "tls"
	AlternateRegistrationAddress string        `astconf:"alternate_registration_address"`
	AlternateRegistrationPort    string        `astconf:"alternate_registration_port"`
	AlternateTransport           string        `astconf:"alternate_transport"`
	FileURLPrefix                string        `astconf:"file_url_prefix"`
	PublicFirmwareURLPrefix      string        `astconf:"public_firmware_url_prefix"`
	NTPServer                    string        `astconf:"ntp_server"`
	SyslogServer                 string        `astconf:"syslog_server"`
	SyslogPort                   int           `astconf:"syslog_port"`
	SyslogLevel                  string        `astconf:"syslog_level"`                // "debug", "error", "warning", "information"
	NetworkVLANDiscoveryMode     string        `astconf:"network_vlan_discovery_mode"` // "NONE", "MANUAL", "LLDP"
	SIPQOS                       int           `astconf:"sip_qos"`
	RTPQOS                       int           `astconf:"rtp_qos"`
	NetworkVLANID                int           `astconf:"network_vlan_id"`
	PCVLANID                     int           `astconf:"pc_vlan_id"`
	PCQOS                        int           `astconf:"pc_qos"`
	SIPDSCP                      int           `astconf:"sip_dscp"`
	RTPDSCP                      int           `astconf:"rtp_dscp"`
	UDPKAInterval                time.Duration `astconf:"udp_ka_interval"` // Seconds
}
