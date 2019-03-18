package dpma

import (
	"strings"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astmerge"
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
		astmerge.StringSlice(&phone.Networks, &merged.Networks)
		astmerge.StringSlice(&phone.Firmware, &merged.Firmware)
		astmerge.String(&phone.MAC, &merged.MAC)
		astmerge.String(&phone.PIN, &merged.PIN)
		astmerge.Int(&phone.GroupPIN, &merged.GroupPIN)
		astmerge.StringSlice(&phone.Lines, &merged.Lines)
		astmerge.String(&phone.ExternalLine, &merged.ExternalLine)
		astmerge.StringSlice(&phone.Applications, &merged.Applications)
		astmerge.String(&phone.ConfigFile, &merged.ConfigFile)
		astmerge.String(&phone.FullName, &merged.FullName)
		astmerge.StringSlice(&phone.Contacts, &merged.Contacts)
		astmerge.String(&phone.ContactsDisplayRules, &merged.ContactsDisplayRules)
		astmerge.String(&phone.BLFContactGroup, &merged.BLFContactGroup)
		astmerge.String(&phone.BLFItems, &merged.BLFItems)
		astmerge.AstSeconds(&phone.BLFPageReturnTimeout, &merged.BLFPageReturnTimeout)
		astmerge.Int(&phone.ContactsMaxSubscriptions, &merged.ContactsMaxSubscriptions)
		astmerge.String(&phone.Timezone, &merged.Timezone)
		astmerge.AstSeconds(&phone.NTPResync, &merged.NTPResync)
		astmerge.Int(&phone.ParkingExtension, &merged.ParkingExtension)
		astmerge.String(&phone.ParkingTransferType, &merged.ParkingTransferType)
		astmerge.AstYesNoNone(&phone.ShowCallParking, &merged.ShowCallParking)
		astmerge.StringSlice(&phone.Ringtones, &merged.Ringtones)
		astmerge.String(&phone.ActiveRingtone, &merged.ActiveRingtone)
		astmerge.AstYesNoNone(&phone.WebUIEnabled, &merged.WebUIEnabled)
		astmerge.AstYesNoNone(&phone.RecordOwnCalls, &merged.RecordOwnCalls)
		astmerge.AstYesNoNone(&phone.CanForwardCalls, &merged.CanForwardCalls)
		astmerge.AstYesNoNone(&phone.ShowCallLog, &merged.ShowCallLog)
		astmerge.AstYesNoNone(&phone.LogOutEnabled, &merged.LogOutEnabled)
		astmerge.StringSlice(&phone.Alerts, &merged.Alerts)
		astmerge.StringSlice(&phone.MulticastPage, &merged.MulticastPage)
		astmerge.AstYesNoNone(&phone.BLFUnusedLineKeys, &merged.BLFUnusedLineKeys)
		astmerge.AstYesNoNone(&phone.SendToVoicemail, &merged.SendToVoicemail)
		mergeLogoFileSlice(&phone.LogoFiles, &merged.LogoFiles)
		astmerge.String(&phone.WallpaperFile, &merged.WallpaperFile)
		astmerge.String(&phone.EHS, &merged.EHS)
		astmerge.AstYesNoNone(&phone.LockPreferences, &merged.LockPreferences)
		astmerge.Int(&phone.LoginPassword, &merged.LoginPassword)
		astmerge.String(&phone.AcceptLocalCalls, &merged.AcceptLocalCalls)
		astmerge.AstYesNoNone(&phone.DisplayMCNotification, &merged.DisplayMCNotification)
		astmerge.String(&phone.IdleCompanyText, &merged.IdleCompanyText)
		astmerge.AstYesNoNone(&phone.SmallClock, &merged.SmallClock)
		astmerge.Int(&phone.DefaultFontSize, &merged.DefaultFontSize)
		astmerge.AstInt(&phone.Brightness, &merged.Brightness)
		astmerge.AstInt(&phone.Contrast, &merged.Contrast)
		astmerge.AstYesNoNone(&phone.DimBacklight, &merged.DimBacklight)
		astmerge.AstSeconds(&phone.BacklightTimeout, &merged.BacklightTimeout)
		astmerge.AstInt(&phone.BacklightDimLevel, &merged.BacklightDimLevel)
		astmerge.String(&phone.ActiveLocale, &merged.ActiveLocale)
		astmerge.AstInt(&phone.RingerVolume, &merged.RingerVolume)
		astmerge.AstInt(&phone.SpeakerVolume, &merged.SpeakerVolume)
		astmerge.AstInt(&phone.HandsetVolume, &merged.HandsetVolume)
		astmerge.AstInt(&phone.HeadsetVolume, &merged.HeadsetVolume)
		astmerge.AstYesNoNone(&phone.CallWaitingTone, &merged.CallWaitingTone)
		astmerge.AstInt(&phone.HandsetSidetoneDB, &merged.HandsetSidetoneDB)
		astmerge.AstInt(&phone.HeadsetSidetoneDB, &merged.HeadsetSidetoneDB)
		astmerge.AstYesNoNone(&phone.ResetCallVolume, &merged.ResetCallVolume)
		astmerge.AstYesNoNone(&phone.HeadsetAnswer, &merged.HeadsetAnswer)
		astmerge.AstYesNoNone(&phone.RingHeadsetOnly, &merged.RingHeadsetOnly)
		astmerge.String(&phone.NameFormat, &merged.NameFormat)
		astmerge.String(&phone.LanPortMode, &merged.LanPortMode)
		astmerge.String(&phone.PCPortMode, &merged.PCPortMode)
		astmerge.AstYesNoNone(&phone.EnableCheckSync, &merged.EnableCheckSync)
		// TODO: Add 802.1x parameters
		astmerge.StringSlice(&phone.Codecs, &merged.Codecs)
		// TODO: Add OpenVPN parameters
		astmerge.AstYesNoNone(&phone.TransportTLSAllowed, &merged.TransportTLSAllowed)
	}
	return
}
