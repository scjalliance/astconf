# Per-server rendering in astgen

Status: design, 2026-10-13. Not yet reviewed by the astconf maintainer.

## Why

The phone system this serves is moving from one Asterisk server to one server per office
plus a hub in GCP, all on one tailnet. `astorg` gained a `Server` model in PR #9
(`Server{ID, Name, Hub, Address}`, `Location.Server`, `Lookup.LocationServer`,
`Lookup.Hub`). Nothing in `astgen` reads it yet: every generator renders the
whole dataset as if every device were local.

This design adds what a multi-server render needs. A rendered config must be
able to say "extension 1234 lives on the Quarry server, and every other server
reaches it over the tailnet trunk", and it must keep the behavior people have
today: one extension rings all of a person's devices, busy lamps light, and
message-waiting lamps light.

The consumer is the consuming generator, which renders one tree per server
(`nodes/<server-id>/`) and commits it to the config repo every server pulls.

## Decisions

Made 2026-10-13 by the project owner:

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
6. **No call waiting in the per-server path** (added 2026-10-13). A busy
   device never rings. The `CustomPresence` "with Call Waiting Enabled"
   override that `PhoneExtensions` renders is not rendered by the per-server
   generators. The existing generators keep theirs.

Taken from the system's decision log:

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

`Placement` is plain data keyed by device and extension, read through per-server
methods (`HostedDevices`, `HomedExtensions`, `Peers`, `Credentials`), so the
consuming generator can also use
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
| `device-conflict` | Desk phones that share a line name sit on different servers |

For `extension-duplicate`, the first claimant in data set order (people, then
roles) keeps the extension and each later one is reported. For
`device-conflict`, such as a person with desk phones at two offices, the first
phone in data set order decides the host and each later one is reported.

`Problem` carries a rule name, the record kind and id, and a detail string,
with no secret in any field. The consuming generator decides whether problems block a
commit.

With a single server that is also the hub, every device and extension is local
and the per-server output matches the existing generators (see Testing).

## Trunks

For each peer, a server renders:

- `endpoint` `trunk-<peer>`: `transport=transport-tcp`, `context=from-trunk`,
  `aors=trunk-<peer>`, `auth` and `outbound_auth` = `auth-trunk-<peer>`,
  `direct_media=no`.
- `aor` `trunk-<peer>`: `contact=sip:<peer address>:5060;transport=tcp`, where
  the address is `Server.Address`, a tailnet hostname.
- `identify` `trunk-<peer>`: `endpoint=trunk-<peer>`, `match=<peer address>`.
- One shared `transport` `transport-tcp`, `protocol=tcp`, bound to an address
  the caller passes in (the server's tailnet IP, which `Server.Address` does
  not carry since it is a hostname). The address is a render input, not a
  dataset field.

**Secrets.** astgen never renders a password. Each trunk endpoint names its
auth section, and `Placement` lists the credential as `{kind: "trunk", id:
<pair id>, section: "auth-trunk-<peer>"}`. The node agent writes that section
into its local secrets file, as it does for phone auth (the system's contracts).

**State publication.** For each peer:

- `outbound-publish` `publish-<peer>`: `server_uri=sip:trunk-<self>@<peer
  address>:5060;transport=tcp`, `event=asterisk-devicestate`,
  `outbound_auth=auth-trunk-<peer>`, `@device_state_filter=<regex of this
  server's own devices>`. The user part `trunk-<self>` is there because the
  peer finds the inbound-publication by the request URI's user. The hub adds
  a second one, `publish-mwi-<peer>`, with the same `server_uri` and
  `event=asterisk-mwi`.
- `inbound-publication` `trunk-<peer>`: `endpoint=trunk-<peer>`,
  `event_asterisk-devicestate=<peer>`, and on office servers
  `event_asterisk-mwi=<peer>` from the hub.
- `asterisk-publication` `<peer>`: `devicestate_publish=publish-<peer>`,
  `device_state=yes`, `device_state_filter=<regex of the peer's devices>`.
  On the hub it also sets `mailboxstate_publish=publish-mwi-<peer>` (the
  `asterisk-mwi` outbound-publish). On an office, the `asterisk-publication`
  for the hub sets `mailbox_state=yes`, which permits the incoming mailbox
  state.

In `res_pjsip_publish_asterisk`, `asterisk-publication`'s `device_state`,
`device_state_filter` and `mailbox_state` govern what a server *accepts*; the
outbound-publish's `@device_state_filter` governs what it *sends*. So each
server sends only the devices it hosts and accepts from each peer only that
peer's devices, and an update never echoes back. A server that hosts no
devices uses `^$`, which matches nothing. A busy-lamp hint for a colleague on
another server names the colleague's device, and lights from the published
state.

**astconf additions.** The `pjsip` package gains typed sections for
`identify`, `outbound-publish`, `inbound-publication` and
`asterisk-publication`, each with encode tests in the existing style.

## Dialplan

Every per-server generator names devices `PJSIP/...`. The existing generators
keep `SIP/...` for the current server.

**Home extension.** For each extension homed on this server:

1. `Gosub(pre-answer,s,1(${EXTEN}))`. The subroutine is hand-written and
   must end in `Return`. The target is Asterisk 22, and Asterisk 21 removed
   `app_macro`, so every per-server subroutine call uses `Gosub`; the
   existing generators keep `Macro`.
2. Build one dial string: each local device is added only when
   `DEVICE_STATE(PJSIP/<device>)` is `NOT_INUSE`, and each remote device is
   added as `PJSIP/dev-<device>@trunk-<host>`.
3. `Dial(<string>,20)` when the string is not empty.
4. Voicemail (below) when nobody answered or the string was empty.

**Device on this server for another home.** `from-trunk` has `dev-<device>`
for each device this server hosts: ring it when idle, otherwise `Busy()`. The
caller's ring loses that one leg and keeps the rest.

**Voicemail.** On the hub, `vm-<mailbox>` in `from-trunk` runs `VoiceMail`.
The `no-answer` subroutine stays hand-written and must end in `Return`;
astgen supplies the target as its second argument:
`Gosub(no-answer,s,1(${EXTEN},<target>))`, read as `${ARG1}` and `${ARG2}`.
On an office the target is `PJSIP/vm-<mailbox>@trunk-<hub>`; on the hub it is
`Local/vm-<mailbox>@from-trunk`, so the subroutine can dial `${ARG2}`
everywhere. A person's mailbox is their extension when they have a voicemail
code; with none, the subroutine gets no target. A role's mailbox is its `MailboxNumber`; a
role with none ends in `Congestion()`, as roles do today.

**Routing table.** Every server renders, in the phones' context, each
extension homed elsewhere as `Dial(PJSIP/<ext>@trunk-<home>)`: offices call
each other directly. An office server falls through unknown extensions to the
hub (`_X.` to `trunk-<hub>`). On the hub an unknown extension is invalid.

**Paging.** The server where a page starts pages its local members and adds
one relay leg for each other server that hosts members. The leg is
`Local/page-relay-<group>-<peer>@astgen-page-relay`, whose extension dials
`PJSIP/page-<group>@trunk-<peer>`. It goes through a local channel because
`Page`'s `s` option skips a recipient that is not idle, and a trunk endpoint
reads as in use while it carries any call. That server's `from-trunk` has
`page-<group>`, which pages its local members only and never relays.

`SIPAddHeader` does nothing on a PJSIP channel, so the per-server `Page`
adds Alert-Info with its `b()` option:
`b(astgen-page-alert^s^1(<alert>))`. astgen renders the `astgen-page-alert`
subroutine, which sets `PJSIP_HEADER(add,Alert-Info)` on each PJSIP callee
and skips local relay channels. The existing paging generator keeps
`SIPAddHeader`.

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
func ServerTrunkExtensions(data *astorg.DataSet, p Placement, serverID string, vmContext string) dialplan.Section // from-trunk
func ServerPagingExtensions(data *astorg.DataSet, p Placement, serverID string, context string) dialplan.Section
func ServerPageRelayExtensions(data *astorg.DataSet, p Placement, serverID string) dialplan.Section // astgen-page-relay
func PageAlertSection() dialplan.Section // astgen-page-alert predial subroutine
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
  through the per-server path to byte-identical endpoints, aors, lines, phones
  and mailboxes; to the same paging apart from the `SIP/` to `PJSIP/` rename
  and Alert-Info moving to the predial subroutine; and to phone extensions
  with the same numbers, comments and devices. Extension actions differ on
  purpose (Decisions 1 and 6), so they are not compared byte for byte.

## Out of scope

- Twilio, inbound DIDs, outbound PSTN and E911 caller ID: hand-written on the
  hub.
- Secrets, the node agent and reload policy: the node agent.
- The tree layout, manifest and commit: the consuming generator.
- The refactor `astgen/doc.go` asks for (generators that do not take whole
  datasets). This design adds new functions beside the old ones and does not
  block that refactor.

## Delivery

Two pull requests to `scjalliance/astconf`, both for the maintainer's review:

1. The four new `pjsip` section types.
2. `Place` and the per-server generators, with the golden fixture. The
   `dialplan` package gains what they need: `Busy`, `VoiceMail`, `Subroutine`
   (Gosub with arguments), `Return`, `PJSIP_HEADER`, `CHANNEL`, `Cat`,
   `Substring`, `DialExpr`, and a predial target on `Page` (its `b()`
   option). A `Page` without one renders as before.

## Amendments (2026-10-13 and 2026-10-14)

Made while writing the implementation plan
(`docs/superpowers/plans/2026-10-13-per-server-rendering.md`), after checking
the Asterisk source. The sections above already read this way.

1. **State publication.** `asterisk-publication`'s `device_state_filter`
   filters incoming state, not outgoing. The own-device filter moved to the
   outbound-publish as `@device_state_filter`, and the
   `asterisk-publication` filter now holds the peer's devices.
   `mailbox_state=yes` (permit incoming) moved from the hub to the receiving
   offices. `server_uri` became `sip:trunk-<self>@<peer address>:5060;transport=tcp`,
   because `res_pjsip_pubsub` finds the inbound-publication by the request
   URI's user part.
2. **Call waiting removed** in the per-server path (Decision 6, project
   owner, 2026-10-13).
3. **Voicemail target** is the no-answer subroutine's second argument, with
   `Local/vm-<mailbox>@from-trunk` on the hub.
4. **Equivalence test scope** is spelled out under Testing.
5. **Duplicates.** First claimant wins, and a fourth problem rule,
   `device-conflict`, covers desk phones sharing a line on two servers.
6. **dialplan helpers** are part of the second pull request (Delivery).
7. **Paging on PJSIP.** Relay legs go through `Local/page-relay-...` so a busy
   trunk is never skipped, and Alert-Info comes from a rendered predial
   subroutine instead of `SIPAddHeader`.
8. **Gosub instead of Macro.** Asterisk 21 removed `app_macro` and the target
   is Asterisk 22, so the per-server dialplan calls the hand-written
   `pre-answer` and `no-answer` subroutines with
   `Gosub(<name>,s,1(${EXTEN}[,<target>]))`. Each must end in `Return`. The
   existing generators keep `Macro`. The `dialplan` package gains a
   `Subroutine` application for this.
9. **Placement keying (2026-10-14).** `Placement` is keyed by device and
   extension and read through per-server methods, not keyed by server ID.
10. **Mailboxes (2026-10-14).** Per-server endpoints never carry `VMCODE`.
    Hub mailboxes carry a fixed locked sentinel password
    (`LockedMailboxPassword`, `x`, rendered `-x`), so a mailbox whose spool
    `secret.conf` is missing is locked rather than open. PINs are read from
    the spool directory (`voicemail.conf` `passwordlocation=spooldir`).
11. **Busy-lamp hints (2026-10-14).** Hints are rendered on every server for
    every extension. `dialplan.Extension` gained a `Hint` field.
12. **Legacy paging (2026-10-14).** The legacy paging generator now shares its
    membership logic with the per-server code. Its output is unchanged, pinned
    by `TestPagingGroupExtensionsLegacy`.
