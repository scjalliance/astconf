package astgen

// TrunkContext is the dialplan context that calls arriving over a trunk
// from another server enter.
const TrunkContext = "from-trunk"

// TrunkTransport is the name of the pjsip transport that every trunk uses.
const TrunkTransport = "transport-tcp"

// trunkPort is the SIP port that every server listens on for trunks.
const trunkPort = "5060"

// trunkName returns the name of the endpoint, aor, identify and
// inbound-publication sections for the trunk to peer.
func trunkName(peer string) string {
	return "trunk-" + peer
}

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

// HubFallthroughContext is the dialplan context that sends every unknown
// extension on an office server to the hub. ServerHubFallthrough renders
// it. The phones' context must include it last, after paging and anything
// else it includes, because Asterisk matches a context's own patterns
// before any included context and _X. would shadow them.
const HubFallthroughContext = "astgen-hub-fallthrough"

// publishName returns the name of the outbound-publish section that sends
// device state to peer.
func publishName(peer string) string {
	return "publish-" + peer
}

// publishMWIName returns the name of the outbound-publish section that sends
// mailbox state to peer.
func publishMWIName(peer string) string {
	return "publish-mwi-" + peer
}
