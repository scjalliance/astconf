package astgen

import (
	"net"
	"regexp"
	"strings"

	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

// TrunkConfig holds the pjsip sections one server needs for its trunks to
// every other server. None of them carries a password: each auth section is
// written by the node agent from the credentials Placement lists.
type TrunkConfig struct {
	Transport            pjsip.Transport
	Endpoints            []pjsip.Endpoint
	AORs                 []pjsip.AOR
	Identifies           []pjsip.Identify
	OutboundPublishes    []pjsip.OutboundPublish
	InboundPublications  []pjsip.InboundPublication
	AsteriskPublications []pjsip.AsteriskPublication
}

// ServerTrunks generates the trunk sections for serverID, one set per peer
// in peer ID order.
//
// The base endpoint is overlaid under each trunk endpoint, which is how a
// trunk gets settings such as its codec list (allow and disallow). The
// trunk's own transport, context, aors, auth, outbound_auth and
// direct_media always win over the base.
//
// bindAddr is the address the trunk transport listens on, which is the
// server's tailnet IP. It is a render input rather than a data set field.
func ServerTrunks(p Placement, serverID string, base pjsip.Endpoint, bindAddr string) TrunkConfig {
	cfg := TrunkConfig{
		Transport: pjsip.Transport{
			Name:     TrunkTransport,
			Protocol: "tcp",
			Bind:     net.JoinHostPort(bindAddr, trunkPort),
		},
	}

	isHub := serverID == p.Hub
	ownFilter := deviceStateFilter(p.HostedDevices(serverID))

	for _, peer := range p.Peers(serverID) {
		name := trunkName(peer.ID)
		auth := trunkAuthName(peer.ID)
		peerURI := "sip:" + net.JoinHostPort(peer.Address, trunkPort) + ";transport=tcp"

		// Overlay rather than merge, so the trunk's aors replace any in the
		// base instead of joining them.
		cfg.Endpoints = append(cfg.Endpoints, pjsip.OverlayEndpoints(base, pjsip.Endpoint{
			Name:         name,
			Transport:    TrunkTransport,
			Context:      TrunkContext,
			Auth:         auth,
			OutboundAuth: auth,
			AORs:         []string{name},
			DirectMedia:  astval.No,
		}))
		cfg.AORs = append(cfg.AORs, pjsip.AOR{
			Name:     name,
			Contacts: []string{peerURI},
		})
		cfg.Identifies = append(cfg.Identifies, pjsip.Identify{
			Name:     name,
			Endpoint: name,
			Match:    []string{peer.Address},
		})

		// The user part of the server URI names the inbound-publication
		// on the peer, which is the peer's trunk to this server.
		publishURI := "sip:" + trunkName(serverID) + "@" + net.JoinHostPort(peer.Address, trunkPort) + ";transport=tcp"
		cfg.OutboundPublishes = append(cfg.OutboundPublishes, pjsip.OutboundPublish{
			Name:              publishName(peer.ID),
			ServerURI:         publishURI,
			Event:             "asterisk-devicestate",
			OutboundAuth:      auth,
			DeviceStateFilter: ownFilter,
		})
		if isHub {
			cfg.OutboundPublishes = append(cfg.OutboundPublishes, pjsip.OutboundPublish{
				Name:         publishMWIName(peer.ID),
				ServerURI:    publishURI,
				Event:        "asterisk-mwi",
				OutboundAuth: auth,
			})
		}

		inbound := pjsip.InboundPublication{
			Name:     name,
			Endpoint: name,
			Events:   pjsip.PublicationEvents{{Event: "asterisk-devicestate", Handler: peer.ID}},
		}
		publication := pjsip.AsteriskPublication{
			Name:               peer.ID,
			DeviceStatePublish: publishName(peer.ID),
			DeviceState:        astval.Yes,
			DeviceStateFilter:  deviceStateFilter(p.HostedDevices(peer.ID)),
		}
		switch {
		case isHub:
			publication.MailboxStatePublish = publishMWIName(peer.ID)
		case peer.ID == p.Hub:
			inbound.Events = append(inbound.Events, pjsip.PublicationEvent{Event: "asterisk-mwi", Handler: peer.ID})
			publication.MailboxState = astval.Yes
		}
		cfg.InboundPublications = append(cfg.InboundPublications, inbound)
		cfg.AsteriskPublications = append(cfg.AsteriskPublications, publication)
	}

	return cfg
}

// deviceStateFilter returns a regular expression that matches the device
// state of exactly the named pjsip devices. With no devices it returns a
// pattern that matches nothing, since every device state name has a
// technology prefix.
func deviceStateFilter(devices []string) string {
	if len(devices) == 0 {
		return "^$"
	}
	quoted := make([]string, len(devices))
	for i, device := range devices {
		quoted[i] = regexp.QuoteMeta(device)
	}
	return "^PJSIP/(" + strings.Join(quoted, "|") + ")$"
}
