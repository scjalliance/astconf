package dpma

import (
	"strings"

	"github.com/scjalliance/astconf"
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
		mergeStringSlice(&phone.Networks, &merged.Networks)
		mergeStringSlice(&phone.Firmware, &merged.Firmware)
		mergeString(&phone.MAC, &merged.MAC)
		mergeString(&phone.PIN, &merged.PIN)
		mergeInt(&phone.GroupPIN, &merged.GroupPIN)
		mergeStringSlice(&phone.Lines, &merged.Lines)
		mergeString(&phone.ExternalLine, &merged.ExternalLine)
		mergeStringSlice(&phone.Applications, &merged.Applications)
		mergeString(&phone.ConfigFile, &merged.ConfigFile)
		mergeString(&phone.FullName, &merged.FullName)
		mergeStringSlice(&phone.Contacts, &merged.Contacts)
		mergeString(&phone.ContactsDisplayRules, &merged.ContactsDisplayRules)
		mergeString(&phone.BLFContactGroup, &merged.BLFContactGroup)
		mergeString(&phone.BLFItems, &merged.BLFItems)
		mergeAstSeconds(&phone.BLFPageReturnTimeout, &merged.BLFPageReturnTimeout)
		mergeInt(&phone.ContactsMaxSubscriptions, &merged.ContactsMaxSubscriptions)
		mergeString(&phone.Timezone, &merged.Timezone)
		mergeAstSeconds(&phone.NTPResync, &merged.NTPResync)
		mergeInt(&phone.ParkingExtension, &merged.ParkingExtension)
		mergeString(&phone.ParkingTransferType, &merged.ParkingTransferType)
		mergeAstYesNoNone(&phone.ShowCallParking, &merged.ShowCallParking)
		mergeStringSlice(&phone.Ringtones, &merged.Ringtones)
		mergeString(&phone.ActiveRingtone, &merged.ActiveRingtone)
		mergeAstYesNoNone(&phone.WebUIEnabled, &merged.WebUIEnabled)
		mergeAstYesNoNone(&phone.RecordOwnCalls, &merged.RecordOwnCalls)
		mergeAstYesNoNone(&phone.CanForwardCalls, &merged.CanForwardCalls)
		mergeAstYesNoNone(&phone.ShowCallLog, &merged.ShowCallLog)
		mergeAstYesNoNone(&phone.LogOutEnabled, &merged.LogOutEnabled)
		mergeStringSlice(&phone.Alerts, &merged.Alerts)
		mergeStringSlice(&phone.MulticastPage, &merged.MulticastPage)
		mergeAstYesNoNone(&phone.BLFUnusedLineKeys, &merged.BLFUnusedLineKeys)
		mergeAstYesNoNone(&phone.SendToVoicemail, &merged.SendToVoicemail)
		mergeLogoFileSlice(&phone.LogoFiles, &merged.LogoFiles)
		mergeString(&phone.WallpaperFile, &merged.WallpaperFile)
		mergeString(&phone.EHS, &merged.EHS)
		mergeAstYesNoNone(&phone.LockPreferences, &merged.LockPreferences)
		mergeInt(&phone.LoginPassword, &merged.LoginPassword)
		mergeString(&phone.AcceptLocalCalls, &merged.AcceptLocalCalls)
		mergeAstYesNoNone(&phone.DisplayMCNotification, &merged.DisplayMCNotification)
		mergeString(&phone.IdleCompanyText, &merged.IdleCompanyText)
		mergeAstYesNoNone(&phone.SmallClock, &merged.SmallClock)
		mergeInt(&phone.DefaultFontSize, &merged.DefaultFontSize)
		mergeAstInt(&phone.Brightness, &merged.Brightness)
		mergeAstInt(&phone.Contrast, &merged.Contrast)
		mergeAstYesNoNone(&phone.DimBacklight, &merged.DimBacklight)
		mergeAstSeconds(&phone.BacklightTimeout, &merged.BacklightTimeout)
		mergeAstInt(&phone.BacklightDimLevel, &merged.BacklightDimLevel)
		mergeString(&phone.ActiveLocale, &merged.ActiveLocale)
		mergeAstInt(&phone.RingerVolume, &merged.RingerVolume)
		mergeAstInt(&phone.SpeakerVolume, &merged.SpeakerVolume)
		mergeAstInt(&phone.HandsetVolume, &merged.HandsetVolume)
		mergeAstInt(&phone.HeadsetVolume, &merged.HeadsetVolume)
		mergeAstYesNoNone(&phone.CallWaitingTone, &merged.CallWaitingTone)
		mergeAstInt(&phone.HandsetSidetoneDB, &merged.HandsetSidetoneDB)
		mergeAstInt(&phone.HeadsetSidetoneDB, &merged.HeadsetSidetoneDB)
		mergeAstYesNoNone(&phone.ResetCallVolume, &merged.ResetCallVolume)
		mergeAstYesNoNone(&phone.HeadsetAnswer, &merged.HeadsetAnswer)
		mergeAstYesNoNone(&phone.RingHeadsetOnly, &merged.RingHeadsetOnly)
		mergeString(&phone.NameFormat, &merged.NameFormat)
		mergeString(&phone.LanPortMode, &merged.LanPortMode)
		mergeString(&phone.PCPortMode, &merged.PCPortMode)
		mergeAstYesNoNone(&phone.EnableCheckSync, &merged.EnableCheckSync)
		// TODO: Add 802.1x parameters
		mergeStringSlice(&phone.Codecs, &merged.Codecs)
		// TODO: Add OpenVPN parameters
		mergeAstYesNoNone(&phone.TransportTLSAllowed, &merged.TransportTLSAllowed)
	}
	return
}
