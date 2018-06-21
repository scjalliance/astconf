package dpma

// Config holds general DPMA configuration settings.
type Config struct {
	// Authentication
	ServerUUID   string `astconf:"server_uuid"`
	GlobalPIN    int    `astconf:"globalpin"`
	UserlistAuth string `astconf:"userlist_auth"` // "disabled", "globalpin"
	ConfigAuth   string `astconf:"config_auth"`   // "mac", "pin", "globalpin", "disabled"

	// mDNS Discovery
	MDNSAddr                string `astconf:"mdns_address"`
	MDNSPort                string `astconf:"mdns_port"`
	MDNSTransport           string `astconf:"mdns_transport"`
	ServiceName             string `astconf:"service_name"`
	ServiceDiscoveryEnabled bool   `astconf:"service_discovery_enabled"`

	// Files
	FileDirectory string `astconf:"file_directory"`

	// Other
	PJSIPMessageContext       string `astconf:"pjsip_message_context"`
	TransactionInitialThreads string `astconf:"transaction_initial_threads"`
	TransactionMaxThreads     string `astconf:"transaction_max_threads"`
}

// SectionName returns the asterisk configuration section name.
func (c *Config) SectionName() string {
	return "general"
}
