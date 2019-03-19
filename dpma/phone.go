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

// OverlayPhones returns the overlayed configuration of all the given phones,
// in order of priority from least to greatest.
func OverlayPhones(phones ...Phone) (overlayed Phone) {
	for i := range phones {
		phone := &phones[i]
		astoverlay.StringSlice(&phone.Networks, &overlayed.Networks)
		astoverlay.StringSlice(&phone.Firmware, &overlayed.Firmware)
		astoverlay.String(&phone.MAC, &overlayed.MAC)
		astoverlay.String(&phone.PIN, &overlayed.PIN)
		astoverlay.Int(&phone.GroupPIN, &overlayed.GroupPIN)
		astoverlay.StringSlice(&phone.Lines, &overlayed.Lines)
		astoverlay.String(&phone.ExternalLine, &overlayed.ExternalLine)
		astoverlay.StringSlice(&phone.Applications, &overlayed.Applications)
		astoverlay.String(&phone.ConfigFile, &overlayed.ConfigFile)
		astoverlay.String(&phone.FullName, &overlayed.FullName)
		astoverlay.StringSlice(&phone.Contacts, &overlayed.Contacts)
		astoverlay.String(&phone.ContactsDisplayRules, &overlayed.ContactsDisplayRules)
		astoverlay.String(&phone.BLFContactGroup, &overlayed.BLFContactGroup)
		astoverlay.String(&phone.BLFItems, &overlayed.BLFItems)
		astoverlay.AstSeconds(&phone.BLFPageReturnTimeout, &overlayed.BLFPageReturnTimeout)
		astoverlay.Int(&phone.ContactsMaxSubscriptions, &overlayed.ContactsMaxSubscriptions)
		astoverlay.String(&phone.Timezone, &overlayed.Timezone)
		astoverlay.AstSeconds(&phone.NTPResync, &overlayed.NTPResync)
		astoverlay.Int(&phone.ParkingExtension, &overlayed.ParkingExtension)
		astoverlay.String(&phone.ParkingTransferType, &overlayed.ParkingTransferType)
		astoverlay.AstYesNoNone(&phone.ShowCallParking, &overlayed.ShowCallParking)
		astoverlay.StringSlice(&phone.Ringtones, &overlayed.Ringtones)
		astoverlay.String(&phone.ActiveRingtone, &overlayed.ActiveRingtone)
		astoverlay.AstYesNoNone(&phone.WebUIEnabled, &overlayed.WebUIEnabled)
		astoverlay.AstYesNoNone(&phone.RecordOwnCalls, &overlayed.RecordOwnCalls)
		astoverlay.AstYesNoNone(&phone.CanForwardCalls, &overlayed.CanForwardCalls)
		astoverlay.AstYesNoNone(&phone.ShowCallLog, &overlayed.ShowCallLog)
		astoverlay.AstYesNoNone(&phone.LogOutEnabled, &overlayed.LogOutEnabled)
		astoverlay.StringSlice(&phone.Alerts, &overlayed.Alerts)
		astoverlay.StringSlice(&phone.MulticastPage, &overlayed.MulticastPage)
		astoverlay.AstYesNoNone(&phone.BLFUnusedLineKeys, &overlayed.BLFUnusedLineKeys)
		astoverlay.AstYesNoNone(&phone.SendToVoicemail, &overlayed.SendToVoicemail)
		overlayLogoFileSlice(&phone.LogoFiles, &overlayed.LogoFiles)
		astoverlay.String(&phone.WallpaperFile, &overlayed.WallpaperFile)
		astoverlay.String(&phone.EHS, &overlayed.EHS)
		astoverlay.AstYesNoNone(&phone.LockPreferences, &overlayed.LockPreferences)
		astoverlay.Int(&phone.LoginPassword, &overlayed.LoginPassword)
		astoverlay.String(&phone.AcceptLocalCalls, &overlayed.AcceptLocalCalls)
		astoverlay.AstYesNoNone(&phone.DisplayMCNotification, &overlayed.DisplayMCNotification)
		astoverlay.String(&phone.IdleCompanyText, &overlayed.IdleCompanyText)
		astoverlay.AstYesNoNone(&phone.SmallClock, &overlayed.SmallClock)
		astoverlay.Int(&phone.DefaultFontSize, &overlayed.DefaultFontSize)
		astoverlay.AstInt(&phone.Brightness, &overlayed.Brightness)
		astoverlay.AstInt(&phone.Contrast, &overlayed.Contrast)
		astoverlay.AstYesNoNone(&phone.DimBacklight, &overlayed.DimBacklight)
		astoverlay.AstSeconds(&phone.BacklightTimeout, &overlayed.BacklightTimeout)
		astoverlay.AstInt(&phone.BacklightDimLevel, &overlayed.BacklightDimLevel)
		astoverlay.String(&phone.ActiveLocale, &overlayed.ActiveLocale)
		astoverlay.AstInt(&phone.RingerVolume, &overlayed.RingerVolume)
		astoverlay.AstInt(&phone.SpeakerVolume, &overlayed.SpeakerVolume)
		astoverlay.AstInt(&phone.HandsetVolume, &overlayed.HandsetVolume)
		astoverlay.AstInt(&phone.HeadsetVolume, &overlayed.HeadsetVolume)
		astoverlay.AstYesNoNone(&phone.CallWaitingTone, &overlayed.CallWaitingTone)
		astoverlay.AstInt(&phone.HandsetSidetoneDB, &overlayed.HandsetSidetoneDB)
		astoverlay.AstInt(&phone.HeadsetSidetoneDB, &overlayed.HeadsetSidetoneDB)
		astoverlay.AstYesNoNone(&phone.ResetCallVolume, &overlayed.ResetCallVolume)
		astoverlay.AstYesNoNone(&phone.HeadsetAnswer, &overlayed.HeadsetAnswer)
		astoverlay.AstYesNoNone(&phone.RingHeadsetOnly, &overlayed.RingHeadsetOnly)
		astoverlay.String(&phone.NameFormat, &overlayed.NameFormat)
		astoverlay.String(&phone.LanPortMode, &overlayed.LanPortMode)
		astoverlay.String(&phone.PCPortMode, &overlayed.PCPortMode)
		astoverlay.AstYesNoNone(&phone.EnableCheckSync, &overlayed.EnableCheckSync)
		// TODO: Add 802.1x parameters
		astoverlay.StringSlice(&phone.Codecs, &overlayed.Codecs)
		// TODO: Add OpenVPN parameters
		astoverlay.AstYesNoNone(&phone.TransportTLSAllowed, &overlayed.TransportTLSAllowed)
	}
	return
}
