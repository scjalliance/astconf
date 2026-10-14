package astgen

// TrunkContext is the dialplan context that calls arriving over a trunk
// from another server enter.
const TrunkContext = "from-trunk"

// TrunkTransport is the name of the pjsip transport that every trunk uses.
const TrunkTransport = "transport-tcp"

// trunkAuthName returns the name of the auth section for the trunk to peer.
func trunkAuthName(peer string) string {
	return "auth-trunk-" + peer
}

// PageAlertContext is the dialplan context of the subroutine that adds a
// paging group's Alert-Info header to each paged PJSIP device.
// PageAlertSection renders it.
const PageAlertContext = "astgen-page-alert"

// PageRelayContext is the dialplan context that holds one relay extension
// for each paging group and each other server with members of it.
// ServerPageRelayExtensions renders it.
const PageRelayContext = "astgen-page-relay"
