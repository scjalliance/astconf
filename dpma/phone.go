package dpma

// Phone is a DPMA phone definition.
type Phone struct {
	Network                  string
	Firmware                 string
	MAC                      string
	PIN                      int
	GroupPIN                 int
	Line                     string
	ExternalLine             string
	Application              string
	ConfigFile               string
	FullName                 string
	Contact                  string
	ContactsDisplayRules     string
	BLFContactGroup          string
	BLFItems                 string
	BLFPageReturnTimeout     string
	ContactsMaxSubscriptions int
	Timezone                 string
	NTPResync                int
	ParkingExtension         string
	ParkingTransferType      string // "blind", "attended"
	ShowCallParking          bool
	Ringtones                []string
	WebUIEnabled             bool
	RecordOwnCalls           bool
	CanForwardCalls          bool
	ShowCallLog              bool
	Alerts                   []string
	MulticastPage            []string
	BLFUnusedLineKeys        bool
	SendToVoicemail          bool
	D40LogoFile              string
	D45LogoFile              string
	D50LogoFile              string
	D60LogoFile              string
	D62LogoFile              string
	D65LogoFile              string
	D70LogoFile              string
	D80LogoFile              string
	WallpaperFile            string
	EHS                      string // "auto", "plantronics", "jabra_iq"
	LockPreferences          bool
	LoginPassword            int
	AcceptLocalCalls         string // "any", "host"
	DisplayMCNotification    string
	IdleCompanyText          string
	SmallClock               bool
	DefaultFontSize          int
	Brightness               int
	Contrast                 int
	DimBacklight             bool
	BacklightTimeout         int
	BacklightDimLevel        int
	ActiveLocale             string
	RingerVolume             int
	SpeakerVolume            int
	HandsetVolume            int
	HeadsetVolume            int
	CallWaitingTone          bool
	HandsetSidetoneDB        int
	HeadsetSidetoneDB        int
	ResetCallVolume          bool
	HeadsetAnswer            bool
	RingHeadsetOnly          bool
	NameFormat               string // "first_last", "last_first"
	LanPortMode              string // "auto", "10hd", "10fd", "100hd", "100fd", "1000fd"
	PCPortMode               string // "auto", "10hd", "10fd", "100hd", "100fd", "1000fd"
	EnableCheckSync          bool

	// TODO: Add 802.1x parameters

	Codecs []string

	// TODO: Add OpenVPN parameters

	TransportTLSAllowed bool
}
