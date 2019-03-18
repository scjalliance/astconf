package dpma

import (
	"strings"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astoverlay"
	"github.com/scjalliance/astconf/astval"
)

// https://wiki.asterisk.org/wiki/display/DIGIUM/DPMA+Configuration#DPMAConfiguration-PhoneConfigurationOptions

// Phone is a DPMA phone definition.
type Phone struct {
	Networks                 []string         `astconf:"network"`
	Firmware                 []string         `astconf:"firmware"`
	MAC                      string           `astconf:"mac,omitempty"`
	PIN                      string           `astconf:"pin,omitempty"`
	GroupPIN                 int              `astconf:"group_pin,omitempty"`
	Lines                    []string         `astconf:"line"` // First entry is primary line
	ExternalLine             string           `astconf:"external_line,omitempty"`
	Applications             []string         `astconf:"application"`
	ConfigFile               string           `astconf:"config_file,omitempty"`
	FullName                 string           `astconf:"full_name,omitempty"`
	Contacts                 []string         `astconf:"contact"`
	ContactsDisplayRules     string           `astconf:"contacts_display_rules,omitempty"`
	BLFContactGroup          string           `astconf:"blf_contact_group,omitempty"`
	BLFItems                 string           `astconf:"blf_items,omitempty"`
	BLFPageReturnTimeout     astval.Seconds   `astconf:"blf_page_return_timeout,omitempty"`
	ContactsMaxSubscriptions int              `astconf:"contacts_max_subscriptions,omitempty"`
	Timezone                 string           `astconf:"timezone,omitempty"`
	NTPResync                astval.Seconds   `astconf:"ntp_resync"`
	ParkingExtension         int              `astconf:"parking_exten,omitempty"`
	ParkingTransferType      string           `astconf:"parking_transfer_type,omitempty"` // "blind", "attended"
	ShowCallParking          astval.YesNoNone `astconf:"show_call_parking"`
	Ringtones                []string         `astconf:"ringtone"`
	ActiveRingtone           string           `astconf:"active_ringtone,omitempty"`
	WebUIEnabled             astval.YesNoNone `astconf:"web_ui_enabled"`
	RecordOwnCalls           astval.YesNoNone `astconf:"record_own_calls"`
	CanForwardCalls          astval.YesNoNone `astconf:"can_forward_calls"`
	ShowCallLog              astval.YesNoNone `astconf:"show_call_log"`
	LogOutEnabled            astval.YesNoNone `astconf:"logout_enabled"`
	Alerts                   []string         `astconf:"alert"`
	MulticastPage            []string         `astconf:"multicastpage"`
	BLFUnusedLineKeys        astval.YesNoNone `astconf:"blf_unused_linekeys"`
	SendToVoicemail          astval.YesNoNone `astconf:"send_to_vm"`
	LogoFiles                []LogoFile
	WallpaperFile            string           `astconf:"wallpaper_file,omitempty"`
	EHS                      string           `astconf:"ehs,omitempty"` // "auto", "plantronics", "jabra_iq"
	LockPreferences          astval.YesNoNone `astconf:"lock_preferences"`
	LoginPassword            int              `astconf:"login_password,omitempty"`
	AcceptLocalCalls         string           `astconf:"accept_local_calls,omitempty"` // "any", "host"
	DisplayMCNotification    astval.YesNoNone `astconf:"display_mc_notification"`
	IdleCompanyText          string           `astconf:"idle_company_text,omitempty"`
	SmallClock               astval.YesNoNone `astconf:"small_clock"`
	DefaultFontSize          int              `astconf:"default_fontsize,omitempty"`
	Brightness               astval.Int       `astconf:"brightness"`
	Contrast                 astval.Int       `astconf:"contrast"`
	DimBacklight             astval.YesNoNone `astconf:"dim_backlight"`
	BacklightTimeout         astval.Seconds   `astconf:"backlight_timeout"`
	BacklightDimLevel        astval.Int       `astconf:"backlight_dim_level"`
	ActiveLocale             string           `astconf:"active_locale,omitempty"`
	RingerVolume             astval.Int       `astconf:"ringer_volume"`
	SpeakerVolume            astval.Int       `astconf:"speaker_volume"`
	HandsetVolume            astval.Int       `astconf:"handset_volume"`
	HeadsetVolume            astval.Int       `astconf:"headset_volume"`
	CallWaitingTone          astval.YesNoNone `astconf:"call_waiting_tone"`
	HandsetSidetoneDB        astval.Int       `astconf:"handset_sidetone_db"`
	HeadsetSidetoneDB        astval.Int       `astconf:"headset_sidetone_db"`
	ResetCallVolume          astval.YesNoNone `astconf:"reset_call_volume"`
	HeadsetAnswer            astval.YesNoNone `astconf:"headset_answer"`
	RingHeadsetOnly          astval.YesNoNone `astconf:"ring_headset_only"`
	NameFormat               string           `astconf:"name_format,omitempty"`   // "first_last", "last_first"
	LanPortMode              string           `astconf:"lan_port_mode,omitempty"` // "auto", "10hd", "10fd", "100hd", "100fd", "1000fd"
	PCPortMode               string           `astconf:"pc_port_mode,omitempty"`  // "auto", "10hd", "10fd", "100hd", "100fd", "1000fd"
	EnableCheckSync          astval.YesNoNone `astconf:"enable_check_sync"`

	// TODO: Add 802.1x parameters

	Codecs []string `astconf:"codecs"`

	// TODO: Add OpenVPN parameters

	TransportTLSAllowed astval.YesNoNone `astconf:"transport_tls_allowed"`
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

// MergePhones returns the merged configuration of all the given phones,
// in order of priority from least to greatest.
func MergePhones(phones ...Phone) (merged Phone) {
	for i := range phones {
		phone := &phones[i]
		astoverlay.StringSlice(&phone.Networks, &merged.Networks)
		astoverlay.StringSlice(&phone.Firmware, &merged.Firmware)
		astoverlay.String(&phone.MAC, &merged.MAC)
		astoverlay.String(&phone.PIN, &merged.PIN)
		astoverlay.Int(&phone.GroupPIN, &merged.GroupPIN)
		astoverlay.StringSlice(&phone.Lines, &merged.Lines)
		astoverlay.String(&phone.ExternalLine, &merged.ExternalLine)
		astoverlay.StringSlice(&phone.Applications, &merged.Applications)
		astoverlay.String(&phone.ConfigFile, &merged.ConfigFile)
		astoverlay.String(&phone.FullName, &merged.FullName)
		astoverlay.StringSlice(&phone.Contacts, &merged.Contacts)
		astoverlay.String(&phone.ContactsDisplayRules, &merged.ContactsDisplayRules)
		astoverlay.String(&phone.BLFContactGroup, &merged.BLFContactGroup)
		astoverlay.String(&phone.BLFItems, &merged.BLFItems)
		astoverlay.AstSeconds(&phone.BLFPageReturnTimeout, &merged.BLFPageReturnTimeout)
		astoverlay.Int(&phone.ContactsMaxSubscriptions, &merged.ContactsMaxSubscriptions)
		astoverlay.String(&phone.Timezone, &merged.Timezone)
		astoverlay.AstSeconds(&phone.NTPResync, &merged.NTPResync)
		astoverlay.Int(&phone.ParkingExtension, &merged.ParkingExtension)
		astoverlay.String(&phone.ParkingTransferType, &merged.ParkingTransferType)
		astoverlay.AstYesNoNone(&phone.ShowCallParking, &merged.ShowCallParking)
		astoverlay.StringSlice(&phone.Ringtones, &merged.Ringtones)
		astoverlay.String(&phone.ActiveRingtone, &merged.ActiveRingtone)
		astoverlay.AstYesNoNone(&phone.WebUIEnabled, &merged.WebUIEnabled)
		astoverlay.AstYesNoNone(&phone.RecordOwnCalls, &merged.RecordOwnCalls)
		astoverlay.AstYesNoNone(&phone.CanForwardCalls, &merged.CanForwardCalls)
		astoverlay.AstYesNoNone(&phone.ShowCallLog, &merged.ShowCallLog)
		astoverlay.AstYesNoNone(&phone.LogOutEnabled, &merged.LogOutEnabled)
		astoverlay.StringSlice(&phone.Alerts, &merged.Alerts)
		astoverlay.StringSlice(&phone.MulticastPage, &merged.MulticastPage)
		astoverlay.AstYesNoNone(&phone.BLFUnusedLineKeys, &merged.BLFUnusedLineKeys)
		astoverlay.AstYesNoNone(&phone.SendToVoicemail, &merged.SendToVoicemail)
		mergeLogoFileSlice(&phone.LogoFiles, &merged.LogoFiles)
		astoverlay.String(&phone.WallpaperFile, &merged.WallpaperFile)
		astoverlay.String(&phone.EHS, &merged.EHS)
		astoverlay.AstYesNoNone(&phone.LockPreferences, &merged.LockPreferences)
		astoverlay.Int(&phone.LoginPassword, &merged.LoginPassword)
		astoverlay.String(&phone.AcceptLocalCalls, &merged.AcceptLocalCalls)
		astoverlay.AstYesNoNone(&phone.DisplayMCNotification, &merged.DisplayMCNotification)
		astoverlay.String(&phone.IdleCompanyText, &merged.IdleCompanyText)
		astoverlay.AstYesNoNone(&phone.SmallClock, &merged.SmallClock)
		astoverlay.Int(&phone.DefaultFontSize, &merged.DefaultFontSize)
		astoverlay.AstInt(&phone.Brightness, &merged.Brightness)
		astoverlay.AstInt(&phone.Contrast, &merged.Contrast)
		astoverlay.AstYesNoNone(&phone.DimBacklight, &merged.DimBacklight)
		astoverlay.AstSeconds(&phone.BacklightTimeout, &merged.BacklightTimeout)
		astoverlay.AstInt(&phone.BacklightDimLevel, &merged.BacklightDimLevel)
		astoverlay.String(&phone.ActiveLocale, &merged.ActiveLocale)
		astoverlay.AstInt(&phone.RingerVolume, &merged.RingerVolume)
		astoverlay.AstInt(&phone.SpeakerVolume, &merged.SpeakerVolume)
		astoverlay.AstInt(&phone.HandsetVolume, &merged.HandsetVolume)
		astoverlay.AstInt(&phone.HeadsetVolume, &merged.HeadsetVolume)
		astoverlay.AstYesNoNone(&phone.CallWaitingTone, &merged.CallWaitingTone)
		astoverlay.AstInt(&phone.HandsetSidetoneDB, &merged.HandsetSidetoneDB)
		astoverlay.AstInt(&phone.HeadsetSidetoneDB, &merged.HeadsetSidetoneDB)
		astoverlay.AstYesNoNone(&phone.ResetCallVolume, &merged.ResetCallVolume)
		astoverlay.AstYesNoNone(&phone.HeadsetAnswer, &merged.HeadsetAnswer)
		astoverlay.AstYesNoNone(&phone.RingHeadsetOnly, &merged.RingHeadsetOnly)
		astoverlay.String(&phone.NameFormat, &merged.NameFormat)
		astoverlay.String(&phone.LanPortMode, &merged.LanPortMode)
		astoverlay.String(&phone.PCPortMode, &merged.PCPortMode)
		astoverlay.AstYesNoNone(&phone.EnableCheckSync, &merged.EnableCheckSync)
		// TODO: Add 802.1x parameters
		astoverlay.StringSlice(&phone.Codecs, &merged.Codecs)
		// TODO: Add OpenVPN parameters
		astoverlay.AstYesNoNone(&phone.TransportTLSAllowed, &merged.TransportTLSAllowed)
	}
	return
}
