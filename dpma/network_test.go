package dpma_test

import (
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/dpma"
)

func TestNetworkMarshal(t *testing.T) {
	network := dpma.Network{
		Name:                "quarry",
		Alias:               "Quarry",
		CIDR:                "10.0.0.0/24",
		RegistrationAddress: "10.0.0.1",
		RegistrationPort:    5060,
		Transport:           "tcp",
		SIPQOS:              astval.NewInt(26),
	}
	got, err := astconf.Marshal(&network)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	// Fields without omitempty are written even when empty.
	want := "[quarry]\n" +
		"alias = Quarry\n" +
		"cidr = 10.0.0.0/24\n" +
		"registration_address = 10.0.0.1\n" +
		"registration_port = 5060\n" +
		"transport = tcp\n" +
		"alternate_registration_address = \n" +
		"alternate_registration_port = 0\n" +
		"alternate_transport = \n" +
		"file_url_prefix = \n" +
		"public_firmware_url_prefix = \n" +
		"ntp_server = \n" +
		"syslog_server = \n" +
		"syslog_level = \n" +
		"network_vlan_discovery_mode = \n" +
		"sip_qos = 26\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestOverlayNetworks(t *testing.T) {
	base := dpma.Network{Name: "base", Transport: "udp", NTPServer: "pool.ntp.org"}
	override := dpma.Network{Name: "quarry", Transport: "tcp", SIPQOS: astval.NewInt(26)}
	got := dpma.OverlayNetworks(&base, &override)
	if got.Name != "quarry" || got.Transport != "tcp" || got.NTPServer != "pool.ntp.org" || got.SIPQOS.Value() != 26 {
		t.Errorf("OverlayNetworks() = %+v", got)
	}
}
