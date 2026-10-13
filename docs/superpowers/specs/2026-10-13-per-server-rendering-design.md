# Per-server rendering in astgen

Status: design, 2026-10-13. Not yet reviewed by the astconf maintainer.

## Why

SCJ's phone system is moving from one Asterisk server to one server per office
plus a hub in GCP, all on one tailnet. `astorg` gained a `Server` model in PR #9
(`Server{ID, Name, Hub, Address}`, `Location.Server`, `Lookup.LocationServer`,
`Lookup.Hub`). Nothing in `astgen` reads it yet: every generator renders the
whole dataset as if every device were local.

This design adds what a multi-server render needs. A rendered config must be
able to say "extension 1234 lives on the Quarry server, and every other server
reaches it over the tailnet trunk", and it must keep the behavior people have
today: one extension rings all of a person's devices, busy lamps light, and
message-waiting lamps light.

The consumer is `scj-phone-gen`, which renders one tree per server
(`nodes/<server-id>/`) and commits it to the config repo every server pulls.

## Decisions

Made 2026-10-13 with Emmaly:

1. **Forked ringing.** An extension rings every device its owner has, on
   whatever server each device registers to. The first device to answer wins
   and the others stop. A busy device drops out of that one ring; busy never
   synchronizes across devices or servers. "No answer" means every device timed
   out or was busy.
2. **Cross-server paging relays.** A page reaches every member of the group,
   on every server.
3. **Trunks are SIP over TCP on the tailnet**, matched by source address and
   authenticated with a per-pair secret derived from the master key.
4. **Busy and message-waiting state is shared now**, by Asterisk's own
   `res_pjsip_publish_asterisk` over the same trunks. No extra component.
5. **API shape: a placement step plus per-server generators.** The existing
   whole-dataset generators stay as they are, so the current single server is
   unaffected until cutover.

Taken from the phone system's decision log (`_phone/docs/decisions.md`):

- Every softphone registers to the hub. A desk phone registers to its office's
  server.
- Offices call each other directly over the tailnet, never through the hub.
- The hub alone carries Twilio, inbound DIDs, outbound PSTN and voicemail.
- A server's ID is permanent and is the node's folder name.

## Placement

```go
// Place decides which server hosts each device and which server owns each
// extension. It is computed once per dataset and shared by every per-server
// generator.
func Place(data *astorg.DataSet) (Placement, []Problem, error)
```

`Placement` is plain data keyed by server ID, so `scj-phone-gen` can also use
it to lay out each node's tree and secret manifest:

- **Device host.** A desk phone is hosted by its location's server
  (`Lookup.LocationServer`). Every softphone is hosted by the hub.
- **Extension home.** A person's or role's extension lives on its location's
  server, or on the hub when that location has no server.
- **Remote devices.** For each extension, the devices hosted on a server other
  than its home, grouped by host.
- **Trunk pairs.** Every unordered pair of servers. The pair id is the two
  server IDs, sorted, joined by `~` (`den~lcy`). This is the `id` of the
  `trunk` credential in the secret manifest, so both ends derive the same
  secret; the sort is the normalization rule the secret-derivation contract
  requires to be fixed.

`Place` returns an error only when the dataset has no hub or more than one,
since nothing can route without exactly one. Everything else is a `Problem`
and the device or extension is left out:

| Problem | Example |
|---|---|
| `location-no-server` | A desk phone at an office that has no server |
| `server-missing` | A location names a server that is not in the dataset |
| `extension-duplicate` | Two people or roles claim the same extension |

`Problem` carries a rule name, the record kind and id, and a detail string,
with no secret in any field. `scj-phone-gen` decides whether problems block a
commit.

With a single server that is also the hub, every device and extension is local
and the per-server output matches the existing generators (see Testing).

## Trunks

For each peer, a server renders:

- `endpoint` `trunk-<peer>`: `transport=transport-tcp`, `context=from-trunk`,
  `aors=trunk-<peer>`, `auth` and `outbound_auth` = `auth-trunk-<peer>`,
  `direct_media=no`.
- `aor` `trunk-<peer>`: `contact=sip:<peer address>:5060;transport=tcp`, where
  the address is `Server.Address` (the tailnet name, `<id>-phone`).
- `identify` `trunk-<peer>`: `endpoint=trunk-<peer>`, `match=<peer address>`.
- One shared `transport` `transport-tcp`, `protocol=tcp`, bound to an address
  the caller passes in (the server's tailnet IP, which `Server.Address` does
  not carry since it is a hostname). The address is a render input, not a
  dataset field.

**Secrets.** astgen never renders a password. Each trunk endpoint names its
auth section, and `Placement` lists the credential as `{kind: "trunk", id:
<pair id>, section: "auth-trunk-<peer>"}`. The node agent writes that section
into its local secrets file, as it does for phone auth (`_phone` contracts,
section 4).

**State publication.** For each peer:

- `outbound-publish` `publish-<peer>`: `server_uri=sip:<peer
  address>:5060;transport=tcp`, `event=asterisk-devicestate`,
  `outbound_auth=auth-trunk-<peer>`. The hub adds a second one,
  `publish-mwi-<peer>`, with `event=asterisk-mwi`.
- `inbound-publication` `trunk-<peer>`: `endpoint=trunk-<peer>`,
  `event_asterisk-devicestate=<peer>`, and on office servers
  `event_asterisk-mwi=<peer>` from the hub.
- `asterisk-publication` `<peer>`: `devicestate_publish=publish-<peer>`,
  `device_state=yes`, `device_state_filter=<regex of this server's own
  devices>`. On the hub it also sets `mailboxstate_publish=publish-mwi-<peer>`
  (the `asterisk-mwi` outbound-publish) and `mailbox_state=yes`.

Each server publishes only the devices it hosts (the filter), so an update
never echoes back. A busy-lamp hint for a colleague on another server names
the colleague's device, and lights from the published state.

**astconf additions.** The `pjsip` package gains typed sections for
`identify`, `outbound-publish`, `inbound-publication` and
`asterisk-publication`, each with encode tests in the existing style.

## Dialplan

Every per-server generator names devices `PJSIP/...`. The existing generators
keep `SIP/...` for the current server.

**Home extension.** For each extension homed on this server:

1. `Macro(pre-answer)`, unchanged and hand-written.
2. Build one dial string: each local device is added only when
   `DEVICE_STATE(PJSIP/<device>)` is `NOT_INUSE`, and each remote device is
   added as `PJSIP/dev-<device>@trunk-<host>`.
3. `Dial(<string>,20)` when the string is not empty.
4. Voicemail (below) when nobody answered or the string was empty.

**Device on this server for another home.** `from-trunk` has `dev-<device>`
for each device this server hosts: ring it when idle, otherwise `Busy()`. The
caller's ring loses that one leg and keeps the rest.

**Voicemail.** On the hub, `vm-<mailbox>` in `from-trunk` runs `VoiceMail`.
On other servers, the no-answer step dials `vm-<mailbox>@trunk-<hub>`. The
`no-answer` macro stays hand-written; astgen supplies the target.

**Routing table.** Every server renders, in the phones' context, each
extension homed elsewhere as `Dial(PJSIP/<ext>@trunk-<home>)`: offices call
each other directly. An office server falls through unknown extensions to the
hub (`_X.` to `trunk-<hub>`). On the hub an unknown extension is invalid.

**Paging.** The server where a page starts pages its local members and adds
one leg `PJSIP/page-<group>@trunk-<peer>` for each other server that hosts
members. That server's `from-trunk` has `page-<group>`, which pages its local
members only and never relays.

**No loops.** `from-trunk` holds only the extensions homed on this server and
the `dev-`, `vm-` and `page-` entries. It never includes the routing table or
the fallthrough, so a call cannot bounce between servers.

## DPMA

Phone, line and busy-lamp sections render only on the phone's host server. A
busy-lamp key for someone on another server names the same `PJSIP/<device>`
hint as a local one and lights from the published state.

## API

```go
func Place(data *astorg.DataSet) (Placement, []Problem, error)

// Per-server generators. Each renders only what serverID needs.
func ServerEndpoints(data *astorg.DataSet, p Placement, serverID string, base pjsip.Endpoint, context string) []pjsip.Endpoint
func ServerAORs(data *astorg.DataSet, p Placement, serverID string, base pjsip.AOR) []pjsip.AOR
func ServerTrunks(p Placement, serverID, bindAddr string) TrunkConfig // endpoints, aors, identifies, transport, publications
func ServerExtensions(data *astorg.DataSet, p Placement, serverID string, context string) dialplan.Section
func ServerTrunkExtensions(data *astorg.DataSet, p Placement, serverID string) dialplan.Section // from-trunk
func ServerPagingExtensions(data *astorg.DataSet, p Placement, serverID string, context string) dialplan.Section
func ServerPhones(data *astorg.DataSet, p Placement, serverID string, base dpma.Phone, contactsURL string) []dpma.Phone
func ServerLines(data *astorg.DataSet, p Placement, serverID string) []dpma.Line
func ServerMailboxes(data *astorg.DataSet, p Placement, serverID string, context string) voicemail.Section // hub only
```

Names are indicative; the plan settles them. Every function is deterministic:
the same input renders the same bytes, with maps iterated in sorted order.

## Testing

- **Unit tests** in the existing style (render to a string, compare with an
  inline `want`) for each new `pjsip` section type and for `Place`: home and
  host rules, the softphone-on-hub rule, sorted pair ids, each problem, and the
  hub-count error.
- **Golden tree.** A multi-server fixture (Quarry, a second office, and the
  hub, with a person who has a desk phone and a softphone, a role, a
  cross-office paging group, and a mailbox) renders every server's output to
  `astgen/testdata/<server>/`. A `-update` flag rewrites the files.
- **Determinism.** Render the fixture twice and compare byte for byte.
- **No regression for the current server.** A single-server fixture renders
  through the per-server path, apart from the deliberate `SIP/` to `PJSIP/`
  device rename, to the same extensions, endpoints, lines and mailboxes as the
  existing generators.

## Out of scope

- Twilio, inbound DIDs, outbound PSTN and E911 caller ID: hand-written on the
  hub.
- Secrets, the node agent and reload policy: `scj-phone`.
- The tree layout, manifest and commit: `scj-phone-gen`.
- The refactor `astgen/doc.go` asks for (generators that do not take whole
  datasets). This design adds new functions beside the old ones and does not
  block that refactor.

## Delivery

Two pull requests to `scjalliance/astconf`, both for the maintainer's review:

1. The four new `pjsip` section types.
2. `Place` and the per-server generators, with the golden fixture.
