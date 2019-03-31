package dpma

import (
	"strings"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astmerge"
	"github.com/scjalliance/astconf/astoverlay"
	"github.com/scjalliance/astconf/astval"
)

// https://wiki.asterisk.org/wiki/display/DIGIUM/DPMA+Configuration#DPMAConfiguration-PhoneConfigurationOptions

// Phone is a DPMA phone definition.
type Phone struct {
	Username                 string           `astconf:"-"`
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
	UseLocalStorage          astval.YesNoNone `astconf:"use_local_storage"`
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

// SectionName returns the name of the phone section.
func (p *Phone) SectionName() string {
	if p.Username != "" {
		return p.Username
	}
	return strings.ToLower(strings.Replace(p.MAC, ":", "", -1))
}

// MarshalAsteriskPreamble marshals the type.
func (p *Phone) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "phone")
}

// OverlayPhones returns the overlayed configuration of all the given phones,
// in order of priority from least to greatest.
func OverlayPhones(phones ...Phone) (overlayed Phone) {
	for i := range phones {
		overlayPhoneScalars(&phones[i], &overlayed)
		overlayPhoneVectors(&phones[i], &overlayed)
	}
	return
}

// MergePhones returns the merged configuration of all the given phones,
// in order of priority from least to greatest.
func MergePhones(phones ...Phone) (merged Phone) {
	for i := range phones {
		overlayPhoneScalars(&phones[i], &merged)
		mergePhoneVectors(&phones[i], &merged)
	}
	return
}

// overlayPhoneScalars overlays all scalar values in from with values from to.
func overlayPhoneScalars(from, to *Phone) {
	astoverlay.String(&from.Username, &to.Username)
	astoverlay.String(&from.MAC, &to.MAC)
	astoverlay.String(&from.PIN, &to.PIN)
	astoverlay.Int(&from.GroupPIN, &to.GroupPIN)
	astoverlay.String(&from.ExternalLine, &to.ExternalLine)
	astoverlay.String(&from.ConfigFile, &to.ConfigFile)
	astoverlay.String(&from.FullName, &to.FullName)
	astoverlay.String(&from.ContactsDisplayRules, &to.ContactsDisplayRules)
	astoverlay.String(&from.BLFContactGroup, &to.BLFContactGroup)
	astoverlay.String(&from.BLFItems, &to.BLFItems)
	astoverlay.AstSeconds(&from.BLFPageReturnTimeout, &to.BLFPageReturnTimeout)
	astoverlay.Int(&from.ContactsMaxSubscriptions, &to.ContactsMaxSubscriptions)
	astoverlay.String(&from.Timezone, &to.Timezone)
	astoverlay.AstSeconds(&from.NTPResync, &to.NTPResync)
	astoverlay.Int(&from.ParkingExtension, &to.ParkingExtension)
	astoverlay.String(&from.ParkingTransferType, &to.ParkingTransferType)
	astoverlay.AstYesNoNone(&from.ShowCallParking, &to.ShowCallParking)
	astoverlay.String(&from.ActiveRingtone, &to.ActiveRingtone)
	astoverlay.AstYesNoNone(&from.WebUIEnabled, &to.WebUIEnabled)
	astoverlay.AstYesNoNone(&from.RecordOwnCalls, &to.RecordOwnCalls)
	astoverlay.AstYesNoNone(&from.CanForwardCalls, &to.CanForwardCalls)
	astoverlay.AstYesNoNone(&from.ShowCallLog, &to.ShowCallLog)
	astoverlay.AstYesNoNone(&from.LogOutEnabled, &to.LogOutEnabled)
	astoverlay.AstYesNoNone(&from.BLFUnusedLineKeys, &to.BLFUnusedLineKeys)
	astoverlay.AstYesNoNone(&from.SendToVoicemail, &to.SendToVoicemail)
	astoverlay.AstYesNoNone(&from.UseLocalStorage, &to.UseLocalStorage)
	astoverlay.String(&from.WallpaperFile, &to.WallpaperFile)
	astoverlay.String(&from.EHS, &to.EHS)
	astoverlay.AstYesNoNone(&from.LockPreferences, &to.LockPreferences)
	astoverlay.Int(&from.LoginPassword, &to.LoginPassword)
	astoverlay.String(&from.AcceptLocalCalls, &to.AcceptLocalCalls)
	astoverlay.AstYesNoNone(&from.DisplayMCNotification, &to.DisplayMCNotification)
	astoverlay.String(&from.IdleCompanyText, &to.IdleCompanyText)
	astoverlay.AstYesNoNone(&from.SmallClock, &to.SmallClock)
	astoverlay.Int(&from.DefaultFontSize, &to.DefaultFontSize)
	astoverlay.AstInt(&from.Brightness, &to.Brightness)
	astoverlay.AstInt(&from.Contrast, &to.Contrast)
	astoverlay.AstYesNoNone(&from.DimBacklight, &to.DimBacklight)
	astoverlay.AstSeconds(&from.BacklightTimeout, &to.BacklightTimeout)
	astoverlay.AstInt(&from.BacklightDimLevel, &to.BacklightDimLevel)
	astoverlay.String(&from.ActiveLocale, &to.ActiveLocale)
	astoverlay.AstInt(&from.RingerVolume, &to.RingerVolume)
	astoverlay.AstInt(&from.SpeakerVolume, &to.SpeakerVolume)
	astoverlay.AstInt(&from.HandsetVolume, &to.HandsetVolume)
	astoverlay.AstInt(&from.HeadsetVolume, &to.HeadsetVolume)
	astoverlay.AstYesNoNone(&from.CallWaitingTone, &to.CallWaitingTone)
	astoverlay.AstInt(&from.HandsetSidetoneDB, &to.HandsetSidetoneDB)
	astoverlay.AstInt(&from.HeadsetSidetoneDB, &to.HeadsetSidetoneDB)
	astoverlay.AstYesNoNone(&from.ResetCallVolume, &to.ResetCallVolume)
	astoverlay.AstYesNoNone(&from.HeadsetAnswer, &to.HeadsetAnswer)
	astoverlay.AstYesNoNone(&from.RingHeadsetOnly, &to.RingHeadsetOnly)
	astoverlay.String(&from.NameFormat, &to.NameFormat)
	astoverlay.String(&from.LanPortMode, &to.LanPortMode)
	astoverlay.String(&from.PCPortMode, &to.PCPortMode)
	astoverlay.AstYesNoNone(&from.EnableCheckSync, &to.EnableCheckSync)
	// TODO: Add 802.1x parameters
	// TODO: Add OpenVPN parameters
	astoverlay.AstYesNoNone(&from.TransportTLSAllowed, &to.TransportTLSAllowed)
}

// overlayPhoneVectors overlays all vector values in from with values from to.
func overlayPhoneVectors(from, to *Phone) {
	astoverlay.StringSlice(&from.Networks, &to.Networks)
	astoverlay.StringSlice(&from.Firmware, &to.Firmware)
	astoverlay.StringSlice(&from.Lines, &to.Lines)
	astoverlay.StringSlice(&from.Applications, &to.Applications)
	astoverlay.StringSlice(&from.Contacts, &to.Contacts)
	astoverlay.StringSlice(&from.Ringtones, &to.Ringtones)
	astoverlay.StringSlice(&from.Alerts, &to.Alerts)
	astoverlay.StringSlice(&from.MulticastPage, &to.MulticastPage)
	overlayLogoFileSlice(&from.LogoFiles, &to.LogoFiles)
	astoverlay.StringSlice(&from.Codecs, &to.Codecs)
}

// mergePhoneVectors merges all vector values in from with values from to.
func mergePhoneVectors(from, to *Phone) {
	astmerge.StringSlice(&from.Networks, &to.Networks)
	astmerge.StringSlice(&from.Firmware, &to.Firmware)
	astmerge.StringSlice(&from.Lines, &to.Lines)
	astmerge.StringSlice(&from.Applications, &to.Applications)
	astmerge.StringSlice(&from.Contacts, &to.Contacts)
	astmerge.StringSlice(&from.Ringtones, &to.Ringtones)
	astmerge.StringSlice(&from.Alerts, &to.Alerts)
	astmerge.StringSlice(&from.MulticastPage, &to.MulticastPage)
	overlayLogoFileSlice(&from.LogoFiles, &to.LogoFiles)
	astmerge.StringSlice(&from.Codecs, &to.Codecs)
}
