package dpma

import (
	"strings"

	"github.com/scjalliance/astconf"
)

// https://wiki.asterisk.org/wiki/display/DIGIUM/DPMA+Configuration#DPMAConfiguration-PhoneConfigurationOptions

// Phone is a DPMA phone definition.
type Phone struct {
	Networks                 []string `astconf:"network"`
	Firmware                 []string `astconf:"firmware"`
	MAC                      string   `astconf:"mac,omitempty"`
	PIN                      int      `astconf:"pin,omitempty"`
	GroupPIN                 int      `astconf:"group_pin,omitempty"`
	Lines                    []string `astconf:"line"` // First entry is primary line
	ExternalLine             string   `astconf:"external_line,omitempty"`
	Applications             []string `astconf:"application"`
	ConfigFile               string   `astconf:"config_file,omitempty"`
	FullName                 string   `astconf:"full_name,omitempty"`
	Contacts                 []string `astconf:"contact"`
	ContactsDisplayRules     string   `astconf:"contacts_display_rules,omitempty"`
	BLFContactGroup          string   `astconf:"blf_contact_group,omitempty"`
	BLFItems                 string   `astconf:"blf_items,omitempty"`
	BLFPageReturnTimeout     int      `astconf:"blf_page_return_timeout,omitempty"`
	ContactsMaxSubscriptions int      `astconf:"contacts_max_subscriptions,omitempty"`
	Timezone                 string   `astconf:"timezone,omitempty"`
	NTPResync                int      `astconf:"ntp_resync,omitempty"`
	ParkingExtension         int      `astconf:"parking_exten,omitempty"`
	ParkingTransferType      string   `astconf:"parking_transfer_type,omitempty"` // "blind", "attended"
	ShowCallParking          bool     `astconf:"show_call_parking,omitempty"`
	Ringtones                []string `astconf:"ringtone,omitempty"`
	ActiveRingtone           string   `astconf:"active_ringtone,omitempty"`
	WebUIEnabled             bool     `astconf:"web_ui_enabled,omitempty"`
	RecordOwnCalls           bool     `astconf:"record_own_calls,omitempty"`
	CanForwardCalls          bool     `astconf:"can_forward_calls,omitempty"`
	ShowCallLog              bool     `astconf:"show_call_log,omitempty"`
	LogOutEnabled            bool     `astconf:"logout_enabled,omitempty"`
	Alerts                   []string `astconf:"alert"`
	MulticastPage            []string `astconf:"multicastpage"`
	BLFUnusedLineKeys        bool     `astconf:"blf_unused_linekeys,omitempty"`
	SendToVoicemail          bool     `astconf:"send_to_vm,omitempty"`
	LogoFiles                []LogoFile
	WallpaperFile            string `astconf:"wallpaper_file,omitempty"`
	EHS                      string `astconf:"ehs,omitempty"` // "auto", "plantronics", "jabra_iq"
	LockPreferences          bool   `astconf:"lock_preferences,omitempty"`
	LoginPassword            int    `astconf:"login_password,omitempty"`
	AcceptLocalCalls         string `astconf:"accept_local_calls,omitempty"` // "any", "host"
	DisplayMCNotification    bool   `astconf:"display_mc_notification,omitempty"`
	IdleCompanyText          string `astconf:"idle_company_text,omitempty"`
	SmallClock               bool   `astconf:"small_clock,omitempty"`
	DefaultFontSize          int    `astconf:"default_fontsize,omitempty"`
	Brightness               int    `astconf:"brightness,omitempty"`
	Contrast                 int    `astconf:"contrast,omitempty"`
	DimBacklight             bool   `astconf:"dim_backlight,omitempty"`
	BacklightTimeout         int    `astconf:"backlight_timeout,omitempty"`
	BacklightDimLevel        int    `astconf:"backlight_dim_level,omitempty"`
	ActiveLocale             string `astconf:"active_locale,omitempty"`
	RingerVolume             int    `astconf:"ringer_volume,omitempty"`
	SpeakerVolume            int    `astconf:"speaker_volume,omitempty"`
	HandsetVolume            int    `astconf:"handset_volume,omitempty"`
	HeadsetVolume            int    `astconf:"headset_volume,omitempty"`
	CallWaitingTone          bool   `astconf:"call_waiting_tone,omitempty"`
	HandsetSidetoneDB        int    `astconf:"handset_sidetone_db,omitempty"`
	HeadsetSidetoneDB        int    `astconf:"headset_sidetone_db,omitempty"`
	ResetCallVolume          bool   `astconf:"reset_call_volume,omitempty"`
	HeadsetAnswer            bool   `astconf:"headset_answer,omitempty"`
	RingHeadsetOnly          bool   `astconf:"ring_headset_only,omitempty"`
	NameFormat               string `astconf:"name_format,omitempty"`   // "first_last", "last_first"
	LanPortMode              string `astconf:"lan_port_mode,omitempty"` // "auto", "10hd", "10fd", "100hd", "100fd", "1000fd"
	PCPortMode               string `astconf:"pc_port_mode,omitempty"`  // "auto", "10hd", "10fd", "100hd", "100fd", "1000fd"
	EnableCheckSync          bool   `astconf:"enable_check_sync,omitempty"`

	// TODO: Add 802.1x parameters

	Codecs []string

	// TODO: Add OpenVPN parameters

	TransportTLSAllowed bool `astconf:"transport_tls_allowed"`
}

// SectionName returns the name of the ringtone section.
func (p *Phone) SectionName() string {
	// FIXME: Generate an alternate name if the phone doesn't have a MAC
	return strings.ToLower(strings.Replace(string(p.MAC), ":", "", -1))
}

// MarshalAsteriskPreamble marshals the type.
func (p *Phone) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "phone")
}
