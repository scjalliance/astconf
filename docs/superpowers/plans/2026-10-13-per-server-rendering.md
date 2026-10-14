# Per-Server Rendering Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Let astgen render one Asterisk server's share of a multi-server phone system (office servers plus a hub) while the existing whole-dataset generators keep rendering today's single server unchanged.

**Architecture:** Phase 1 adds four pjsip section types (identify, outbound-publish, inbound-publication, asterisk-publication) in the same pattern as `pjsip.Transport`. Phase 2 adds `astgen.Place`, which decides once per data set which server hosts each device and owns each extension, and a set of `Server*` generators that read that `Placement` and render only what one server needs. Device sections reuse the existing generators and filter them; dialplan and trunk sections are new code.

**Tech Stack:** Go 1.25 (module `github.com/scjalliance/astconf`), the repo's own `astconf` encoder, standard library only.

**Spec:** `docs/superpowers/specs/2026-10-13-per-server-rendering-design.md`. Read it before starting. Where this plan departs from it, the reason is in "Spec decisions" below; those departures are deliberate and reviewed.

## Global Constraints

- Go only. No Python, shell scripts or other languages in the repo or in tooling for this work.
- `gofmt -l .` prints nothing, `go vet ./...` is clean, and `go test ./...` passes at the end of every task.
- No em dashes or en dashes anywhere: code comments, docs, test strings, commit messages.
- Output is deterministic: the same input renders the same bytes. Never range over a map to produce output without sorting the keys first.
- astgen's new code never renders a password. Secrets reach a node only through its agent, from the credentials `Placement` lists.
- The existing generators (`SIP`, `Endpoints`, `Auths`, `AORs`, `PhoneExtensions`, `PagingGroupExtensions`, `Phones`, `Lines`, `Mailboxes` and the rest) are not edited and their output does not change. `SIP/` device names stay in those functions; only the new `Server*` functions name devices `PJSIP/`.
- Existing tests are not edited. New tests go in new files.
- Every commit message ends with these two lines, after a blank line:

  ```
  Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
  Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
  ```

- Never push. Work stays on the local branch `per-server-rendering`. The controller opens each pull request.
- `astconf-eval` is a working copy of `scjalliance/astconf`, which has another maintainer. Keep each change small and in the style of the file next to it.

## Review Focus

Five inputs the spec implies but does not spell out, most likely to bite first. Each has a test in the task that owns the code.

1. **An extension homed on an office whose owner has no device there, only a softphone on the hub.** The ring must still dial the one remote leg, with no local `DEVICE_STATE` check. Test: `TestHomeExtensionOnlyRemoteDevices` in Task 11.
2. **A role (or person) at an office with no server.** Its extension lives on the hub, every office routes it to the hub, and it has no remote devices. Test: `TestPlaceRoleAtOfficeWithoutServer` in Task 7, with the rendered routes in `TestServerExtensionsOffice` and `TestServerExtensionsHub` in Task 11.
3. **A server with no locations, so it hosts no devices.** It still gets a trunk to every peer, and its device state filters match nothing rather than everything. Test: `TestServerTrunksServerWithoutDevices` in Task 9.
4. **A paging group whose members are all on other servers.** A page from the hub reaches each office through one relay leg, each leg is a local channel so a busy trunk cannot drop it, and the hub's `from-trunk` has no `page-` entry for that group. Tests: `TestServerPagingAllMembersRemote` in Task 10 and `TestServerTrunkExtensionsNoPageWithoutMembers` in Task 11.
5. **Server IDs that sort differently than server names.** Peers, trunk pairs, credentials and every per-peer section must follow ID order. The fixture's names (Hub, Quarry, Rock Den) sort differently than its IDs (den, hub, qry). Test: `TestPlaceSortsByID` in Task 7, and the peer order in `TestServerTrunksOffice` in Task 9.

## Spec decisions

These settle what the spec left open or got wrong. Each was checked against the Asterisk source where noted.

1. **Own-device filter lives on the outbound-publish.** In `res_pjsip_publish_asterisk.c`, `asterisk-publication`'s `device_state_filter` filters *incoming* events ("Optional regular expression used to filter what devices we accept events for"). The outgoing filter is the outbound-publish extended option `@device_state_filter`, read with `ast_sorcery_object_get_extended(configuration, "device_state_filter")`. So `publish-<peer>` carries `@device_state_filter=<this server's devices>`, and `asterisk-publication <peer>` carries `device_state_filter=<the peer's devices>`. The spec's intent (each server publishes only its own devices) holds; Asterisk also drops events that did not originate locally by EID, so nothing echoes.
2. **`mailbox_state=yes` goes on the receiver.** `mailbox_state` means "permit incoming mailbox state events". The hub sets `mailboxstate_publish=publish-mwi-<peer>` and does not set `mailbox_state`; each office sets `mailbox_state=yes` on its `asterisk-publication hub`.
3. **The outbound-publish URI carries a user part.** `res_pjsip_pubsub.c` finds the inbound-publication by the request URI's user part (`ast_sorcery_retrieve_by_id(..., "inbound-publication", resource_name)`). Server S publishing to peer P sends to `sip:trunk-S@<P address>:5060;transport=tcp`, which matches P's `inbound-publication trunk-S`.
4. **Gosub, not Macro, and the voicemail target is the no-answer subroutine's second argument.** The target is Asterisk 22, and Asterisk 21 removed `app_macro`, so the per-server dialplan calls the hand-written subroutines with `Gosub`: `Gosub(pre-answer,s,1(${EXTEN}))` and `Gosub(no-answer,s,1(${EXTEN},<target>))`. They read `${ARG1}` and `${ARG2}` and must end in `Return`. On an office the target is `PJSIP/vm-<mailbox>@trunk-<hub>`; on the hub it is `Local/vm-<mailbox>@from-trunk`, so the subroutine can `Dial(${ARG2})` everywhere. A person's mailbox is their extension when they have a voicemail code; with none, the subroutine gets no target. A role's mailbox is `MailboxNumber`; a role with none ends in `Congestion()`, as today. The old `PhoneExtensions` keeps `Macro`.
5. **"Same extensions" in the equivalence test.** Decision 1 of the spec changes how a ring treats busy devices, so dialplan actions cannot match the old generator. The test requires byte equality for endpoints, aors, lines, phones and mailboxes; byte equality for paging after renaming `SIP/` to `PJSIP/`; and, for phone extensions, the same numbers, comments and set of devices.
6. **Call waiting is removed in the per-server path** (project owner's decision, 2026-10-13; spec Decisions). The old `PhoneExtensions` rings every device when `CustomPresence` says "with Call Waiting Enabled". The per-server ring adds a local device only when it is `NOT_INUSE`, so a busy device never rings and the override is not rendered. The old generators keep theirs.
7. **Duplicate extensions.** The first claimant in data set order (people, then roles) keeps the extension; each later one is reported as `extension-duplicate` and left out.
8. **One extra problem rule, `device-conflict`.** Two desk phones that share a line name (a person with desk phones at two offices) cannot be hosted on two servers under one endpoint name. The first phone in data set order wins; the later one is reported and left out.
9. **dialplan additions.** The ring needs a dial string built at run time and `Busy()` and `VoiceMail()`; paging needs a predial subroutine. Task 6 adds `Busy`, `VoiceMail`, `Cat`, `Substring`, `DialExpr`, `Subroutine`, `Return`, `PJSIPHeader`, `Channel`, `PredialTarget` and a `Predial` field on `PageApp` to the dialplan package in PR 2. A `PageApp` with no `Predial` renders exactly as before, so the old paging output is unchanged.
10. **No `ServerAuths`.** Phone auth sections carry passwords, so astgen does not render them per server. The node agent writes them from `Placement.HostedDevices`, as the system's contracts already describe for phone auth.
11. **`bindAddr` is a host, not host and port.** `ServerTrunks` appends `:5060` with `net.JoinHostPort`.

12. **Paging Alert-Info uses Page's `b()` predial option, and astgen renders the subroutine.** `SIPAddHeader` does nothing on a PJSIP channel. Every per-server `Page` for a group with an alert gets `b(astgen-page-alert^s^1(<alert>))`, which runs on each callee channel before it is dialed. `PageAlertSection()` renders the `astgen-page-alert` context, so the subroutine travels with the config instead of being hand-written on each node. It sets `PJSIP_HEADER(add,Alert-Info)` only when `CHANNEL(channeltype)` is `PJSIP`: a relay leg is a local channel with no SIP headers, and the receiving server adds the header for its own members. The old `PagingGroupExtensions` keeps `SIPAddHeader` and its output is unchanged.
13. **Each paging relay leg is a local channel.** `Page`'s `s` option skips a recipient that is not idle, and a trunk endpoint reads as in use while it carries any call. The originating `Page` dials `Local/page-relay-<group>-<peer>@astgen-page-relay`, and `ServerPageRelayExtensions` renders that context, where each extension dials `PJSIP/page-<group>@trunk-<peer>`. A local channel reads idle unless the same relay is already running, so `s` still judges every local member device. Dropping `s` from the originating `Page` was the other choice; it was rejected because it would also page busy local phones. The relay context is separate from `from-trunk`, so the no-loop rule holds.

## File structure

Phase 1 (PR 1), all in `pjsip/`:

| File | Responsibility |
|---|---|
| `identify.go`, `identify_test.go` | `Identify` section type |
| `outboundpublish.go`, `outboundpublish_test.go` | `OutboundPublish` section type |
| `inboundpublication.go`, `inboundpublication_test.go` | `InboundPublication`, `PublicationEvent`, `PublicationEvents` |
| `asteriskpublication.go`, `asteriskpublication_test.go` | `AsteriskPublication` section type |
| `example_publish_test.go` | Runnable example of a state-sharing peer |
| `doc.go` (modify) | Package doc lists the new types |

The new types follow `pjsip.Transport`: written once per peer, so no overlay, merge or map types. `Endpoint`, `Auth` and `AOR` have those because they are generated per device and merged with a base.

Phase 2 (PR 2):

| File | Responsibility |
|---|---|
| `dialplan/busy.go`, `dialplan/voicemail.go`, `dialplan/cat.go`, `dialplan/dialexpr.go`, `dialplan/return.go`, `dialplan/pjsipheader.go`, `dialplan/channel.go`, `dialplan/predial.go`, `dialplan/perserver_test.go` | Applications, functions and expressions the per-server dialplan needs |
| `dialplan/page.go` (modify) | `PageApp.Predial`, rendered as Page's `b()` option |
| `astgen/names.go` | Section and context names shared by the trunk and dialplan code |
| `astgen/placement.go`, `astgen/placement_test.go` | `Place`, `Placement`, `Problem`, `TrunkPair`, `Credential` |
| `astgen/multiserver_test.go` | The multi-server fixture every Phase 2 test uses |
| `astgen/server.go`, `astgen/server_test.go` | `ServerEndpoints`, `ServerAORs`, `ServerLines`, `ServerPhones`, `ServerMailboxes` |
| `astgen/trunk.go`, `astgen/trunk_test.go` | `ServerTrunks`, `TrunkConfig` |
| `astgen/serverpaging.go`, `astgen/serverpaging_test.go` | `ServerPagingExtensions`, `ServerPageRelayExtensions`, `PageAlertSection` and the paging helpers |
| `astgen/serverdialplan.go`, `astgen/serverdialplan_test.go` | `ServerExtensions`, `ServerTrunkExtensions` |
| `astgen/golden_test.go`, `astgen/testdata/<server>/*` | Golden tree, `-update` flag, determinism, no-password check |
| `astgen/equivalence_test.go` | Single-server equivalence with the old generators |
| `astgen/doc.go` (modify) | Package doc describes the two families of generators |

## How to check rendered output

Every expected string in this plan was produced by running the code in this plan. If a test fails only on whitespace or ordering, do not paste the actual output into the test until you have checked it means what the step says:

- Each section starts with `[name]` and, for pjsip types, `type = <kind>` on the next line.
- Sections are separated by exactly one blank line; there is no trailing blank line.
- Settings are `name = value` with one space each side of `=`.
- A `;` in a value is written `\;` (the encoder escapes it; Asterisk reads it back as `;`). In Go source that is `\\;`.
- Dialplan lines are `exten => N,1,App(...)` then `same => n,App(...)`, with a blank line and `; Comment` before each extension.

If the only difference is one of those mechanical rules, the code is right and the plan's string has a typo: fix the string. Any other difference is a bug in the code.

---

# Phase 1: pjsip section types (PR 1)

### Task 1: `pjsip.Identify`

**Files:**
- Create: `pjsip/identify.go`
- Test: `pjsip/identify_test.go`

**Interfaces:**
- Consumes: `astconf.Encoder`, `astval.YesNoNone` (existing).
- Produces: `type Identify struct { Name string; Templates []string; Endpoint string; Match []string; SRVLookups astval.YesNoNone; MatchHeader string }` with pointer methods `SectionName() string`, `SectionTemplates() []string`, `MarshalAsteriskPreamble(*astconf.Encoder) error`. Task 9 sets `Name`, `Endpoint`, `Match`.

Fields: `endpoint`, `match`, `srv_lookups`, `match_header`. Left out: `match_request_uri` (newer Asterisk, no use here). `match` is written as one line per value, like `transport.local_net`, which Asterisk accepts.

- [ ] **Step 1: Write the failing test**

Create `pjsip/identify_test.go`:

```go
package pjsip_test

import (
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

func TestIdentifyMinimal(t *testing.T) {
	got, err := astconf.Marshal(&pjsip.Identify{Name: "trunk-den"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[trunk-den]\ntype = identify\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestIdentifyFull(t *testing.T) {
	identify := pjsip.Identify{
		Name:        "trunk-den",
		Endpoint:    "trunk-den",
		Match:       []string{"den-phone", "100.64.0.0/10"},
		SRVLookups:  astval.No,
		MatchHeader: "X-Trunk: den",
	}
	got, err := astconf.Marshal(&identify)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[trunk-den]\n" +
		"type = identify\n" +
		"endpoint = trunk-den\n" +
		"match = den-phone\n" +
		"match = 100.64.0.0/10\n" +
		"srv_lookups = no\n" +
		"match_header = X-Trunk: den\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./pjsip -run TestIdentify -v`
Expected: FAIL to compile with `undefined: pjsip.Identify`.

- [ ] **Step 3: Write the implementation**

Create `pjsip/identify.go`:

```go
package pjsip

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
)

// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip/#identify

// Identify is a pjsip identify section. It matches inbound requests to an
// endpoint by source address or by a header, which is how a trunk with no
// registration is recognized.
//
// Identify sections are written once per peer and are not generated per
// device, so there are no overlay or merge functions for them.
//
// Only the fields that multi-server rendering needs are modeled.
type Identify struct {
	Name        string           `astconf:"-"`
	Templates   []string         `astconf:"-"`
	Endpoint    string           `astconf:"endpoint,omitempty"`
	Match       []string         `astconf:"match,omitempty"` // IP addresses, networks or hostnames
	SRVLookups  astval.YesNoNone `astconf:"srv_lookups"`
	MatchHeader string           `astconf:"match_header,omitempty"` // "Header-Name: value"
}

// SectionName returns the name of the identify section.
func (identify *Identify) SectionName() string {
	return identify.Name
}

// SectionTemplates returns the section templates used by the identify.
func (identify *Identify) SectionTemplates() []string {
	return identify.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (identify *Identify) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "identify")
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./pjsip -run TestIdentify -v && gofmt -l . && go vet ./pjsip`
Expected: both tests PASS; gofmt and vet print nothing.

- [ ] **Step 5: Commit**

```bash
git add pjsip/identify.go pjsip/identify_test.go
git commit -F - <<'EOF'
feat(pjsip): add identify section type

An identify section matches inbound requests to an endpoint by source
address or header, which is how a trunk to another server is recognized.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 2: `pjsip.OutboundPublish`

**Files:**
- Create: `pjsip/outboundpublish.go`
- Test: `pjsip/outboundpublish_test.go`

**Interfaces:**
- Consumes: `astconf.Encoder`, `astval.Seconds`, `astval.Int` (existing).
- Produces: `type OutboundPublish struct { Name string; Templates []string; ServerURI, FromURI, ToURI, Event, OutboundAuth string; Expiration astval.Seconds; Transport, OutboundProxy string; MaxAuthAttempts astval.Int; DeviceStateFilter, MailboxStateFilter string }` with the same three pointer methods. Task 9 sets `Name`, `ServerURI`, `Event`, `OutboundAuth`, `DeviceStateFilter`.

Fields: `server_uri`, `from_uri`, `to_uri`, `event`, `outbound_auth`, `expiration`, `transport`, `outbound_proxy`, `max_auth_attempts`, plus the res_pjsip_publish_asterisk extended options `@device_state_filter` and `@mailbox_state_filter` (Spec decision 1). Left out: `multi_user` and `payload` (no use here). `outbound_auth` is a string like `Endpoint.OutboundAuth`.

- [ ] **Step 1: Write the failing test**

Create `pjsip/outboundpublish_test.go`:

```go
package pjsip_test

import (
	"testing"
	"time"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

func TestOutboundPublishMinimal(t *testing.T) {
	got, err := astconf.Marshal(&pjsip.OutboundPublish{Name: "publish-den"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[publish-den]\ntype = outbound-publish\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestOutboundPublishFull(t *testing.T) {
	publish := pjsip.OutboundPublish{
		Name:               "publish-den",
		ServerURI:          "sip:trunk-lcy@den-phone:5060;transport=tcp",
		FromURI:            "sip:lcy@lcy-phone",
		ToURI:              "sip:den@den-phone",
		Event:              "asterisk-devicestate",
		OutboundAuth:       "auth-trunk-den",
		Expiration:         astval.NewSeconds(time.Hour),
		Transport:          "transport-tcp",
		OutboundProxy:      "sip:proxy.example.com",
		MaxAuthAttempts:    astval.NewInt(3),
		DeviceStateFilter:  "^PJSIP/(fred)$",
		MailboxStateFilter: "^100@default$",
	}
	got, err := astconf.Marshal(&publish)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	// A semicolon starts a comment in asterisk configuration files, so the
	// encoder escapes each one as "\;". Asterisk reads "\;" back as ";".
	want := "[publish-den]\n" +
		"type = outbound-publish\n" +
		"server_uri = sip:trunk-lcy@den-phone:5060\\;transport=tcp\n" +
		"from_uri = sip:lcy@lcy-phone\n" +
		"to_uri = sip:den@den-phone\n" +
		"event = asterisk-devicestate\n" +
		"outbound_auth = auth-trunk-den\n" +
		"expiration = 3600\n" +
		"transport = transport-tcp\n" +
		"outbound_proxy = sip:proxy.example.com\n" +
		"max_auth_attempts = 3\n" +
		"@device_state_filter = ^PJSIP/(fred)$\n" +
		"@mailbox_state_filter = ^100@default$\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./pjsip -run TestOutboundPublish -v`
Expected: FAIL to compile with `undefined: pjsip.OutboundPublish`.

- [ ] **Step 3: Write the implementation**

Create `pjsip/outboundpublish.go`:

```go
package pjsip

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
)

// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip_outbound_publish/#outbound-publish
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip_publish_asterisk/

// OutboundPublish is a pjsip outbound-publish section. It sends PUBLISH
// requests for one event to one server, such as device or mailbox state
// for another Asterisk server.
//
// The user part of ServerURI names the inbound-publication section that
// receives the requests on the other server.
//
// DeviceStateFilter and MailboxStateFilter are the extended options that
// res_pjsip_publish_asterisk reads from an outbound-publish. They limit
// which local devices or mailboxes are published, and are written with a
// leading "@" as Asterisk requires for extended options.
//
// Outbound-publish sections are written once per peer and are not
// generated per device, so there are no overlay or merge functions for
// them.
//
// Only the fields that multi-server rendering needs are modeled.
type OutboundPublish struct {
	Name               string         `astconf:"-"`
	Templates          []string       `astconf:"-"`
	ServerURI          string         `astconf:"server_uri,omitempty"`
	FromURI            string         `astconf:"from_uri,omitempty"`
	ToURI              string         `astconf:"to_uri,omitempty"`
	Event              string         `astconf:"event,omitempty"` // "asterisk-devicestate", "asterisk-mwi"
	OutboundAuth       string         `astconf:"outbound_auth,omitempty"`
	Expiration         astval.Seconds `astconf:"expiration"`
	Transport          string         `astconf:"transport,omitempty"`
	OutboundProxy      string         `astconf:"outbound_proxy,omitempty"`
	MaxAuthAttempts    astval.Int     `astconf:"max_auth_attempts"`
	DeviceStateFilter  string         `astconf:"@device_state_filter,omitempty"`  // Regular expression
	MailboxStateFilter string         `astconf:"@mailbox_state_filter,omitempty"` // Regular expression
}

// SectionName returns the name of the outbound-publish section.
func (publish *OutboundPublish) SectionName() string {
	return publish.Name
}

// SectionTemplates returns the section templates used by the
// outbound-publish.
func (publish *OutboundPublish) SectionTemplates() []string {
	return publish.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (publish *OutboundPublish) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "outbound-publish")
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./pjsip -run TestOutboundPublish -v && gofmt -l . && go vet ./pjsip`
Expected: both tests PASS; gofmt and vet print nothing. Note `server_uri` renders with `\;` (see "How to check rendered output").

- [ ] **Step 5: Commit**

```bash
git add pjsip/outboundpublish.go pjsip/outboundpublish_test.go
git commit -F - <<'EOF'
feat(pjsip): add outbound-publish section type

Includes the @device_state_filter and @mailbox_state_filter extended
options that res_pjsip_publish_asterisk reads to limit what is sent.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 3: `pjsip.InboundPublication`

**Files:**
- Create: `pjsip/inboundpublication.go`
- Test: `pjsip/inboundpublication_test.go`

**Interfaces:**
- Consumes: `astconf.Encoder`, `(*astconf.Printer).Setting` (existing).
- Produces: `type InboundPublication struct { Name string; Templates []string; Endpoint string; Events PublicationEvents }`, `type PublicationEvent struct { Event, Handler string }`, `type PublicationEvents []PublicationEvent` with `MarshalAsterisk(*astconf.Encoder) error`. Task 9 builds `PublicationEvents{{Event: "asterisk-devicestate", Handler: peerID}}`.

Fields: `endpoint` and the `event_<name>` family. The setting name depends on the data, so `PublicationEvents` writes its own settings. Because it implements `astconf.Marshaler`, the encoder treats it as a block and writes it after the plain fields, which is the order Asterisk examples use.

- [ ] **Step 1: Write the failing test**

Create `pjsip/inboundpublication_test.go`:

```go
package pjsip_test

import (
	"errors"
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/pjsip"
)

func TestInboundPublicationMinimal(t *testing.T) {
	got, err := astconf.Marshal(&pjsip.InboundPublication{Name: "trunk-den"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[trunk-den]\ntype = inbound-publication\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestInboundPublicationFull(t *testing.T) {
	publication := pjsip.InboundPublication{
		Name:     "trunk-hub",
		Endpoint: "trunk-hub",
		Events: pjsip.PublicationEvents{
			{Event: "asterisk-devicestate", Handler: "hub"},
			{Event: "asterisk-mwi", Handler: "hub"},
		},
	}
	got, err := astconf.Marshal(&publication)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[trunk-hub]\n" +
		"type = inbound-publication\n" +
		"endpoint = trunk-hub\n" +
		"event_asterisk-devicestate = hub\n" +
		"event_asterisk-mwi = hub\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

// TestInboundPublicationInvalidEvent checks that an event name that would
// break the setting name is rejected instead of written.
func TestInboundPublicationInvalidEvent(t *testing.T) {
	publication := pjsip.InboundPublication{
		Name:   "trunk-den",
		Events: pjsip.PublicationEvents{{Event: "a=b", Handler: "den"}},
	}
	_, err := astconf.Marshal(&publication)
	var ice astconf.InvalidContentError
	if !errors.As(err, &ice) {
		t.Fatalf("Marshal() error = %v, want InvalidContentError", err)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./pjsip -run TestInboundPublication -v`
Expected: FAIL to compile with `undefined: pjsip.InboundPublication`.

- [ ] **Step 3: Write the implementation**

Create `pjsip/inboundpublication.go`:

```go
package pjsip

import (
	"github.com/scjalliance/astconf"
)

// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip_pubsub/#inbound-publication

// InboundPublication is a pjsip inbound-publication section. It accepts
// PUBLISH requests for one resource and hands each event to a handler.
//
// The section name must match the user part of the request URI that the
// publishing server sends to, which is the user part of its outbound-publish
// server_uri.
//
// Inbound-publication sections are written once per peer and are not
// generated per device, so there are no overlay or merge functions for
// them.
type InboundPublication struct {
	Name      string   `astconf:"-"`
	Templates []string `astconf:"-"`
	Endpoint  string   `astconf:"endpoint,omitempty"` // The only endpoint allowed to publish
	Events    PublicationEvents
}

// SectionName returns the name of the inbound-publication section.
func (publication *InboundPublication) SectionName() string {
	return publication.Name
}

// SectionTemplates returns the section templates used by the
// inbound-publication.
func (publication *InboundPublication) SectionTemplates() []string {
	return publication.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (publication *InboundPublication) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "inbound-publication")
}

// PublicationEvent maps one event of an inbound publication to the
// configuration that handles it. For the asterisk-devicestate and
// asterisk-mwi events, Handler is the name of an asterisk-publication
// section.
type PublicationEvent struct {
	Event   string // "asterisk-devicestate", "asterisk-mwi"
	Handler string
}

// PublicationEvents is an ordered list of publication events. Each one is
// written as its own "event_<event> = <handler>" setting, because Asterisk
// reads every setting whose name starts with "event_" as an event mapping.
type PublicationEvents []PublicationEvent

// MarshalAsterisk writes one setting for each event, in order.
func (events PublicationEvents) MarshalAsterisk(e *astconf.Encoder) error {
	for _, event := range events {
		if err := e.Printer().Setting("event_"+event.Event, event.Handler); err != nil {
			return err
		}
	}
	return nil
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./pjsip -run TestInboundPublication -v && gofmt -l . && go vet ./pjsip`
Expected: three tests PASS. The invalid-event test passes because `Printer.Setting` rejects `=` in a name, and `errors.As` finds the `InvalidContentError` through the encoder's `MarshalerError` wrapper.

- [ ] **Step 5: Commit**

```bash
git add pjsip/inboundpublication.go pjsip/inboundpublication_test.go
git commit -F - <<'EOF'
feat(pjsip): add inbound-publication section type

Event mappings are written as event_<name> settings, one per event.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 4: `pjsip.AsteriskPublication`

**Files:**
- Create: `pjsip/asteriskpublication.go`
- Test: `pjsip/asteriskpublication_test.go`

**Interfaces:**
- Consumes: `astconf.Encoder`, `astval.YesNoNone` (existing).
- Produces: `type AsteriskPublication struct { Name string; Templates []string; DeviceStatePublish, MailboxStatePublish string; DeviceState astval.YesNoNone; DeviceStateFilter string; MailboxState astval.YesNoNone; MailboxStateFilter string }` with the same three pointer methods. Task 9 sets every field but `MailboxStateFilter`.

Fields: all six options (`devicestate_publish`, `mailboxstate_publish`, `device_state`, `device_state_filter`, `mailbox_state`, `mailbox_state_filter`). None left out.

- [ ] **Step 1: Write the failing test**

Create `pjsip/asteriskpublication_test.go`:

```go
package pjsip_test

import (
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

func TestAsteriskPublicationMinimal(t *testing.T) {
	got, err := astconf.Marshal(&pjsip.AsteriskPublication{Name: "den"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[den]\ntype = asterisk-publication\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestAsteriskPublicationFull(t *testing.T) {
	publication := pjsip.AsteriskPublication{
		Name:                "den",
		DeviceStatePublish:  "publish-den",
		MailboxStatePublish: "publish-mwi-den",
		DeviceState:         astval.Yes,
		DeviceStateFilter:   "^PJSIP/(barney)$",
		MailboxState:        astval.No,
		MailboxStateFilter:  "^200@default$",
	}
	got, err := astconf.Marshal(&publication)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[den]\n" +
		"type = asterisk-publication\n" +
		"devicestate_publish = publish-den\n" +
		"mailboxstate_publish = publish-mwi-den\n" +
		"device_state = yes\n" +
		"device_state_filter = ^PJSIP/(barney)$\n" +
		"mailbox_state = no\n" +
		"mailbox_state_filter = ^200@default$\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./pjsip -run TestAsteriskPublication -v`
Expected: FAIL to compile with `undefined: pjsip.AsteriskPublication`.

- [ ] **Step 3: Write the implementation**

Create `pjsip/asteriskpublication.go`:

```go
package pjsip

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
)

// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip_publish_asterisk/

// AsteriskPublication is a pjsip asterisk-publication section. It ties
// device and mailbox state sharing with one other Asterisk server together.
//
// DeviceState and MailboxState permit incoming state from that server, and
// the two filters limit which incoming devices and mailboxes are accepted.
//
// Sending does not depend on this section. An outbound-publish with
// event=asterisk-devicestate or event=asterisk-mwi publishes on its own, and
// its @device_state_filter or @mailbox_state_filter limits what it sends.
// DeviceStatePublish and MailboxStatePublish name the outbound-publish that
// is used to request a full state refresh from that server.
//
// Asterisk-publication sections are written once per peer and are not
// generated per device, so there are no overlay or merge functions for
// them.
type AsteriskPublication struct {
	Name                string           `astconf:"-"`
	Templates           []string         `astconf:"-"`
	DeviceStatePublish  string           `astconf:"devicestate_publish,omitempty"`  // Name of an outbound-publish section
	MailboxStatePublish string           `astconf:"mailboxstate_publish,omitempty"` // Name of an outbound-publish section
	DeviceState         astval.YesNoNone `astconf:"device_state"`
	DeviceStateFilter   string           `astconf:"device_state_filter,omitempty"` // Regular expression
	MailboxState        astval.YesNoNone `astconf:"mailbox_state"`
	MailboxStateFilter  string           `astconf:"mailbox_state_filter,omitempty"` // Regular expression
}

// SectionName returns the name of the asterisk-publication section.
func (publication *AsteriskPublication) SectionName() string {
	return publication.Name
}

// SectionTemplates returns the section templates used by the
// asterisk-publication.
func (publication *AsteriskPublication) SectionTemplates() []string {
	return publication.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (publication *AsteriskPublication) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "asterisk-publication")
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./pjsip -run TestAsteriskPublication -v && gofmt -l . && go vet ./pjsip`
Expected: both tests PASS.

- [ ] **Step 5: Commit**

```bash
git add pjsip/asteriskpublication.go pjsip/asteriskpublication_test.go
git commit -F - <<'EOF'
feat(pjsip): add asterisk-publication section type

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 5: Package doc and example, then PR 1 boundary

**Files:**
- Modify: `pjsip/doc.go` (whole file)
- Create: `pjsip/example_publish_test.go`

**Interfaces:**
- Consumes: the four types from Tasks 1 to 4.
- Produces: nothing new for later tasks.

- [ ] **Step 1: Write the example (it is the test)**

Create `pjsip/example_publish_test.go`:

```go
package pjsip_test

import (
	"bytes"
	"fmt"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

// ExampleAsteriskPublication shares device state with one other server,
// den, over a trunk that is recognized by its source address.
func ExampleAsteriskPublication() {
	identify := pjsip.Identify{
		Name:     "trunk-den",
		Endpoint: "trunk-den",
		Match:    []string{"den-phone"},
	}
	publish := pjsip.OutboundPublish{
		Name:              "publish-den",
		ServerURI:         "sip:trunk-lcy@den-phone:5060;transport=tcp",
		Event:             "asterisk-devicestate",
		OutboundAuth:      "auth-trunk-den",
		DeviceStateFilter: "^PJSIP/(fred)$",
	}
	inbound := pjsip.InboundPublication{
		Name:     "trunk-den",
		Endpoint: "trunk-den",
		Events:   pjsip.PublicationEvents{{Event: "asterisk-devicestate", Handler: "den"}},
	}
	publication := pjsip.AsteriskPublication{
		Name:               "den",
		DeviceStatePublish: "publish-den",
		DeviceState:        astval.Yes,
		DeviceStateFilter:  "^PJSIP/(barney)$",
	}

	var buf bytes.Buffer
	e := astconf.NewEncoder(&buf)
	for _, v := range []interface{}{&identify, &publish, &inbound, &publication} {
		if err := e.Encode(v); err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Print(buf.String())
	// Output:
	// [trunk-den]
	// type = identify
	// endpoint = trunk-den
	// match = den-phone
	//
	// [publish-den]
	// type = outbound-publish
	// server_uri = sip:trunk-lcy@den-phone:5060\;transport=tcp
	// event = asterisk-devicestate
	// outbound_auth = auth-trunk-den
	// @device_state_filter = ^PJSIP/(fred)$
	//
	// [trunk-den]
	// type = inbound-publication
	// endpoint = trunk-den
	// event_asterisk-devicestate = den
	//
	// [den]
	// type = asterisk-publication
	// devicestate_publish = publish-den
	// device_state = yes
	// device_state_filter = ^PJSIP/(barney)$
}
```

- [ ] **Step 2: Run it**

Run: `go test ./pjsip -run ExampleAsteriskPublication -v`
Expected: PASS. (It passes as soon as it is written because Tasks 1 to 4 are done; if it fails, the diff points at the type whose output changed.)

- [ ] **Step 3: Update the package doc**

Replace the whole of `pjsip/doc.go` with:

```go
// Package pjsip provides types for writing chan_pjsip configuration
// sections to pjsip.conf.
//
// Each PJSIP object type (endpoint, auth, aor, transport, identify,
// outbound-publish, inbound-publication, asterisk-publication) is a
// separate Go type. Unlike chan_sip, a single device is described by an
// endpoint section, an auth section and an aor section that reference each
// other by name. The astgen package generates all three for each phone.
//
// A trunk to another Asterisk server adds an identify section, which
// recognizes the server by its address, and the three publication types,
// which share device and mailbox state with it through
// res_pjsip_publish_asterisk.
//
// The section name, templates and type preamble are provided by pointer
// receiver methods, so pass a pointer to the encoder, as in
// e.Encode(&endpoint). Encoding a value instead of a pointer writes the
// settings without the section header. This matches the sip and dpma
// packages.
//
// Only a subset of the PJSIP settings is modeled. Fields are added as
// configurations need them.
package pjsip
```

- [ ] **Step 4: Full verification**

Run: `gofmt -l . && go vet ./... && go test ./...`
Expected: gofmt prints nothing; every package `ok`.

Run: `git grep -nP '[\x{2013}\x{2014}]' -- pjsip`
Expected: no output.

- [ ] **Step 5: Commit**

```bash
git add pjsip/doc.go pjsip/example_publish_test.go
git commit -F - <<'EOF'
docs(pjsip): describe the trunk and publication section types

Adds a runnable example of one state-sharing peer.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

- [ ] **Step 6: PR 1 boundary. Stop.**

Phase 1 ends here. Do not start Task 6. Report to the controller: the five commit SHAs and the output of `go test ./...`. The controller opens PR 1 (the four pjsip types) against `scjalliance/astconf` and decides when Phase 2 starts. Do not push.

---

# Phase 2: Placement and per-server generators (PR 2)

Phase 2 builds on Phase 1. Start it from the branch with Tasks 1 to 5 committed.

### Task 6: dialplan applications for the per-server ring and paging

**Files:**
- Create: `dialplan/busy.go`, `dialplan/voicemail.go`, `dialplan/cat.go`, `dialplan/dialexpr.go`, `dialplan/subroutine.go`, `dialplan/return.go`, `dialplan/pjsipheader.go`, `dialplan/channel.go`, `dialplan/predial.go`
- Modify: `dialplan/page.go` (whole file; adds `Predial`)
- Test: `dialplan/perserver_test.go`

**Interfaces:**
- Consumes: `dialplan.AppCall`, `dialplan.ExprDef`, `dialplan.Expression`, `dialplan.Var`, `errorIfAny`, `validate`, `invalidArgChars`, `invalidResourceChars` (existing).
- Produces:
  - `func Busy() BusyApp`
  - `func VoiceMail(mailbox, context, options string) VoiceMailApp`
  - `func Cat(parts ...Expression) CatExpr` (renders the parts joined, as a literal)
  - `func Substring(name Var, offset int) SubstringExpr` (renders `${name:offset}`)
  - `func DialExpr(target Expression, seconds int) DialExprApp` (renders `Dial(<target>[,seconds])`)
  - `func Subroutine(context, extension string, priority int, args ...Expression) SubroutineApp` (renders `Gosub(context,extension,priority(arg1,arg2))`; the existing `GosubApp` writes arguments after the priority as separate arguments and accepts no variables, so it cannot express `${EXTEN}`)
  - `func Return() ReturnApp`
  - `func PJSIPHeader(action, name string) PJSIPHeaderFunc` (a `Reference`, so it can be `Set`'s name)
  - `func Channel(item string) ChannelFunc` (quoted when compared)
  - `type PredialTarget struct { Context, Extension string; Priority int; Args []string }`, rendered `context^extension^priority(arg1^arg2)`
  - `PageApp.Predial *PredialTarget`, rendered as `b(<target>)` after the other options; nil renders as before

- [ ] **Step 1: Write the failing test**

Create `dialplan/perserver_test.go`:

```go
package dialplan_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/dialplan"
)

// TestPerServerApps renders the applications and expressions that the
// per-server generators in astgen use to build a dial string at run time.
func TestPerServerApps(t *testing.T) {
	dialVar := dialplan.Var("ASTGEN_DIAL")
	section := dialplan.Section{
		Context: "test",
		Extensions: []dialplan.Extension{
			{
				Number: "100",
				Actions: []dialplan.Action{
					dialplan.Set(dialVar, dialplan.String("&PJSIP/dev-fred.soft@trunk-hub")),
					dialplan.ExecIf(
						dialplan.Equal(dialplan.DeviceState(dialplan.PJSIP("fred")), dialplan.String("NOT_INUSE")),
						dialplan.Set(dialVar, dialplan.Cat(dialVar, dialplan.String("&PJSIP/fred"))),
					),
					dialplan.ExecIf(
						dialplan.GreaterThan(dialplan.Len(dialVar), dialplan.Int(0)),
						dialplan.DialExpr(dialplan.Substring(dialVar, 1), 20),
					),
					dialplan.Busy(),
				},
			},
			{
				Number: "vm-100",
				Actions: []dialplan.Action{
					dialplan.VoiceMail("100", "default", "u"),
					dialplan.VoiceMail("200", "", ""),
				},
			},
		},
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&section); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}
	want := "[test]\n" +
		"\n" +
		"exten => 100,1,Set(ASTGEN_DIAL=&PJSIP/dev-fred.soft@trunk-hub)\n" +
		"same => n,ExecIf($[${DEVICE_STATE(PJSIP/fred)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/fred))\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Busy()\n" +
		"\n" +
		"exten => vm-100,1,VoiceMail(100@default,u)\n" +
		"same => n,VoiceMail(200)\n"
	if got := buf.String(); got != want {
		t.Errorf("Encode() =\n%s\nwant:\n%s", got, want)
	}
}

// TestPerServerAppsInvalid checks that the new applications reject text
// that would change the structure of a dialplan line.
func TestPerServerAppsInvalid(t *testing.T) {
	tests := []struct {
		name   string
		action dialplan.Action
	}{
		{name: "comma in mailbox", action: dialplan.VoiceMail("1,2", "default", "")},
		{name: "paren in voicemail options", action: dialplan.VoiceMail("100", "default", "u)")},
		{name: "comma in cat part", action: dialplan.Set(dialplan.Var("X"), dialplan.Cat(dialplan.String("a,b")))},
		{name: "colon in substring name", action: dialplan.DialExpr(dialplan.Substring(dialplan.Var("A:B"), 1), 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := dialplan.Section{Context: "test", Extensions: []dialplan.Extension{{Number: "1", Actions: []dialplan.Action{tt.action}}}}
			var buf bytes.Buffer
			err := astconf.NewEncoder(&buf).Encode(&section)
			var ice astconf.InvalidContentError
			if !errors.As(err, &ice) {
				t.Fatalf("Encode() error = %v, want InvalidContentError", err)
			}
		})
	}
}

// TestPagePredialAndHeaders renders Page's predial option and the pieces of
// the subroutine it calls. A Page without a predial target renders as it
// always has.
func TestPagePredialAndHeaders(t *testing.T) {
	plain := dialplan.Page(dialplan.PJSIP("fred"))
	alerted := dialplan.Page(dialplan.PJSIP("fred"), dialplan.Device{Technology: "Local", Resource: "page-relay-500-den@astgen-page-relay"})
	alerted.Predial = &dialplan.PredialTarget{Context: "astgen-page-alert", Extension: "s", Priority: 1, Args: []string{"intercom"}}
	section := dialplan.Section{
		Context: "test",
		Extensions: []dialplan.Extension{
			{Number: "500", Actions: []dialplan.Action{plain, alerted}},
			{
				Number: "s",
				Actions: []dialplan.Action{
					dialplan.ExecIf(
						dialplan.Equal(dialplan.Channel("channeltype"), dialplan.String("PJSIP")),
						dialplan.Set(dialplan.PJSIPHeader("add", "Alert-Info"), dialplan.Var("ARG1")),
					),
					dialplan.Return(),
				},
			},
		},
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&section); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}
	want := "[test]\n" +
		"\n" +
		"exten => 500,1,Page(PJSIP/fred,is)\n" +
		"same => n,Page(PJSIP/fred&Local/page-relay-500-den@astgen-page-relay,isb(astgen-page-alert^s^1(intercom)))\n" +
		"\n" +
		"exten => s,1,ExecIf($[\"${CHANNEL(channeltype)}\"=\"PJSIP\"]?Set(PJSIP_HEADER(add,Alert-Info)=${ARG1}))\n" +
		"same => n,Return()\n"
	if got := buf.String(); got != want {
		t.Errorf("Encode() =\n%s\nwant:\n%s", got, want)
	}
}

// TestPagePredialInvalid checks that a predial target that would break the
// option is rejected.
func TestPagePredialInvalid(t *testing.T) {
	tests := []struct {
		name   string
		target dialplan.PredialTarget
	}{
		{name: "caret in argument", target: dialplan.PredialTarget{Context: "a", Extension: "s", Priority: 1, Args: []string{"x^y"}}},
		{name: "paren in context", target: dialplan.PredialTarget{Context: "a)", Extension: "s", Priority: 1}},
		{name: "priority zero", target: dialplan.PredialTarget{Context: "a", Extension: "s"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := dialplan.Page(dialplan.PJSIP("fred"))
			target := tt.target
			page.Predial = &target
			section := dialplan.Section{Context: "test", Extensions: []dialplan.Extension{{Number: "1", Actions: []dialplan.Action{page}}}}
			var buf bytes.Buffer
			if err := astconf.NewEncoder(&buf).Encode(&section); err == nil {
				t.Fatalf("Encode() succeeded, want an error:\n%s", buf.String())
			}
		})
	}
}

// TestSubroutine renders Gosub calls with arguments in parentheses after the
// priority, which is how Asterisk passes them as ${ARG1}, ${ARG2}.
func TestSubroutine(t *testing.T) {
	section := dialplan.Section{
		Context: "test",
		Extensions: []dialplan.Extension{
			{
				Number: "100",
				Actions: []dialplan.Action{
					dialplan.Subroutine("pre-answer", "s", 1, dialplan.Var("EXTEN")),
					dialplan.Subroutine("no-answer", "s", 1, dialplan.Var("EXTEN"), dialplan.String("PJSIP/vm-100@trunk-hub")),
					dialplan.Subroutine("cleanup", "s", 1),
				},
			},
		},
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&section); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}
	want := "[test]\n" +
		"\n" +
		"exten => 100,1,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-100@trunk-hub))\n" +
		"same => n,Gosub(cleanup,s,1)\n"
	if got := buf.String(); got != want {
		t.Errorf("Encode() =\n%s\nwant:\n%s", got, want)
	}
}

// TestSubroutineInvalid checks that a Gosub call that would break the line
// is rejected.
func TestSubroutineInvalid(t *testing.T) {
	tests := []struct {
		name   string
		action dialplan.Action
	}{
		{name: "comma in context", action: dialplan.Subroutine("a,b", "s", 1)},
		{name: "paren in argument", action: dialplan.Subroutine("a", "s", 1, dialplan.String("x)"))},
		{name: "priority zero", action: dialplan.Subroutine("a", "s", 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := dialplan.Section{Context: "test", Extensions: []dialplan.Extension{{Number: "1", Actions: []dialplan.Action{tt.action}}}}
			var buf bytes.Buffer
			if err := astconf.NewEncoder(&buf).Encode(&section); err == nil {
				t.Fatalf("Encode() succeeded, want an error:\n%s", buf.String())
			}
		})
	}
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `go test ./dialplan -run 'TestPerServer|TestPagePredial|TestSubroutine' -v`
Expected: FAIL to compile with `undefined: dialplan.Cat` (and the other new names).

- [ ] **Step 3: Write the implementation**

Create `dialplan/busy.go`:

```go
package dialplan

// Busy returns a dialplan application that signals busy to the caller.
func Busy() BusyApp {
	return BusyApp{}
}

// BusyApp is a dialplan application that signals busy to the caller.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Applications/Busy/
type BusyApp struct{}

// App returns the assembled application call.
func (app BusyApp) App() AppCall {
	return AppCall{Name: "Busy"}
}
```

Create `dialplan/voicemail.go`:

```go
package dialplan

// VoiceMail returns a dialplan application that leaves a message in a
// mailbox. Options are passed through, such as "u" for the unavailable
// greeting.
func VoiceMail(mailbox, context, options string) VoiceMailApp {
	return VoiceMailApp{
		Mailbox: mailbox,
		Context: context,
		Options: options,
	}
}

// VoiceMailApp is a dialplan application that leaves a message in a
// mailbox.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Applications/VoiceMail/
type VoiceMailApp struct {
	Mailbox string
	Context string
	Options string
}

// App returns the assembled application call.
func (app VoiceMailApp) App() AppCall {
	box := app.Mailbox
	if app.Context != "" {
		box += "@" + app.Context
	}
	call := AppCall{Name: "VoiceMail", Args: []string{box}}
	if app.Options != "" {
		call.Args = append(call.Args, app.Options)
	}
	return call
}

// Validate returns an error if the mailbox, context or options contain
// invalid characters.
func (app VoiceMailApp) Validate() error {
	if err := errorIfAny("voicemail mailbox", app.Mailbox, invalidResourceChars); err != nil {
		return err
	}
	if err := errorIfAny("voicemail context", app.Context, invalidResourceChars); err != nil {
		return err
	}
	return errorIfAny("voicemail options", app.Options, invalidArgChars)
}
```

Create `dialplan/cat.go`:

```go
package dialplan

import (
	"strconv"
	"strings"
)

// Cat returns an expression that joins the rendered form of each part, such
// as a variable followed by literal text.
func Cat(parts ...Expression) CatExpr {
	return CatExpr{Parts: parts}
}

// CatExpr is an expression that joins the rendered form of its parts.
type CatExpr struct {
	Parts []Expression
}

// Expr returns the joined parts as a literal expression.
func (cat CatExpr) Expr() ExprDef {
	var builder strings.Builder
	for _, part := range cat.Parts {
		builder.WriteString(part.Expr().String())
	}
	return ExprDef{Content: builder.String(), Kind: StringLit}
}

// Validate returns an error if any part is invalid.
func (cat CatExpr) Validate() error {
	for _, part := range cat.Parts {
		if err := validate(part); err != nil {
			return err
		}
	}
	return nil
}

// Substring returns an expression that reads a variable from the given
// offset onward, as ${name:offset}.
func Substring(name Var, offset int) SubstringExpr {
	return SubstringExpr{Name: name, Offset: offset}
}

// SubstringExpr is an expression that reads part of a variable.
type SubstringExpr struct {
	Name   Var
	Offset int
}

// Expr returns the substring as a variable reference.
func (sub SubstringExpr) Expr() ExprDef {
	return ExprDef{Content: string(sub.Name) + ":" + strconv.Itoa(sub.Offset)}
}

// Validate returns an error if the variable name contains invalid
// characters.
func (sub SubstringExpr) Validate() error {
	return errorIfAny("variable name", string(sub.Name), invalidArgChars+":")
}
```

Create `dialplan/dialexpr.go`:

```go
package dialplan

import "strconv"

// DialExpr returns a dial application whose dial string is computed when
// the call runs, such as a variable that earlier actions built up.
func DialExpr(target Expression, seconds int) DialExprApp {
	return DialExprApp{
		Target:  target,
		Timeout: seconds,
	}
}

// DialExprApp is a dialplan application that dials a computed dial string.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Applications/Dial/
type DialExprApp struct {
	Target  Expression
	Timeout int // Seconds
}

// App returns the assembled application call.
func (dial DialExprApp) App() AppCall {
	app := AppCall{
		Name: "Dial",
		Args: []string{dial.Target.Expr().String()},
	}
	if dial.Timeout > 0 {
		app.Args = append(app.Args, strconv.Itoa(dial.Timeout))
	}
	return app
}

// Validate returns an error if the target is invalid.
func (dial DialExprApp) Validate() error {
	return validate(dial.Target)
}
```

Create `dialplan/subroutine.go`:

```go
package dialplan

import (
	"fmt"
	"strconv"
	"strings"
)

// Subroutine returns a dialplan application that calls a subroutine with
// Gosub, passing args in parentheses after the priority. The subroutine
// reads them as ${ARG1}, ${ARG2} and so on, and must end in Return.
//
// Unlike GosubApp, which writes each argument after the priority as a
// separate application argument, Subroutine writes them the way Asterisk
// passes them, and each argument is an expression, so it may be a variable.
func Subroutine(context, extension string, priority int, args ...Expression) SubroutineApp {
	return SubroutineApp{
		Context:   context,
		Extension: extension,
		Priority:  priority,
		Args:      args,
	}
}

// SubroutineApp is a dialplan application that calls a subroutine with
// arguments.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Applications/Gosub/
type SubroutineApp struct {
	Context   string
	Extension string
	Priority  int // Should be >= 1
	Args      []Expression
}

// App returns the assembled application call.
func (sub SubroutineApp) App() AppCall {
	target := strconv.Itoa(sub.Priority)
	if len(sub.Args) > 0 {
		args := make([]string, 0, len(sub.Args))
		for _, arg := range sub.Args {
			args = append(args, arg.Expr().String())
		}
		target += "(" + strings.Join(args, ",") + ")"
	}
	return AppCall{
		Name: "Gosub",
		Args: []string{sub.Context, sub.Extension, target},
	}
}

// Validate returns an error if the priority is below one, or if the context,
// extension or any argument is invalid.
func (sub SubroutineApp) Validate() error {
	if sub.Priority < 1 {
		return fmt.Errorf("gosub priority must be at least 1, got %d", sub.Priority)
	}
	if err := errorIfAny("gosub context", sub.Context, invalidArgChars); err != nil {
		return err
	}
	if err := errorIfAny("gosub extension", sub.Extension, invalidArgChars); err != nil {
		return err
	}
	for _, arg := range sub.Args {
		if err := validate(arg); err != nil {
			return err
		}
	}
	return nil
}
```

Create `dialplan/return.go`:

```go
package dialplan

// Return returns a dialplan application that returns from a subroutine.
func Return() ReturnApp {
	return ReturnApp{}
}

// ReturnApp is a dialplan application that returns from a subroutine
// entered with Gosub or a predial option.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Applications/Return/
type ReturnApp struct{}

// App returns the assembled application call.
func (app ReturnApp) App() AppCall {
	return AppCall{Name: "Return"}
}
```

Create `dialplan/pjsipheader.go`:

```go
package dialplan

// PJSIPHeader returns a dialplan function that reads or changes a SIP
// header on a PJSIP channel. Action is "read", "add", "update" or
// "remove".
func PJSIPHeader(action, name string) PJSIPHeaderFunc {
	return PJSIPHeaderFunc{
		Action: action,
		Name:   name,
	}
}

// PJSIPHeaderFunc is a dialplan function that reads or changes a SIP header
// on a PJSIP channel.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Functions/PJSIP_HEADER/
type PJSIPHeaderFunc struct {
	Action string
	Name   string
}

// Func returns the assembled function call.
func (f PJSIPHeaderFunc) Func() FuncCall {
	return FuncCall{
		Name: "PJSIP_HEADER",
		Args: []string{f.Action, f.Name},
	}
}

// Expr returns the function as an expression.
func (f PJSIPHeaderFunc) Expr() ExprDef {
	return ExprDef{Content: f.Func().String()}
}

// Ref returns a reference to the function, for use as the name in Set.
func (f PJSIPHeaderFunc) Ref() string {
	return f.Func().String()
}

// Validate returns an error if the action or header name contains invalid
// characters.
func (f PJSIPHeaderFunc) Validate() error {
	if err := errorIfAny("pjsip header action", f.Action, invalidArgChars); err != nil {
		return err
	}
	return errorIfAny("pjsip header name", f.Name, invalidArgChars)
}
```

Create `dialplan/channel.go`:

```go
package dialplan

// Channel returns a dialplan function that reads an item of information
// about the current channel, such as "channeltype".
func Channel(item string) ChannelFunc {
	return ChannelFunc{Item: item}
}

// ChannelFunc is a dialplan function that reads information about the
// current channel.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Functions/CHANNEL/
type ChannelFunc struct {
	Item string
}

// Func returns the assembled function call.
func (f ChannelFunc) Func() FuncCall {
	return FuncCall{
		Name: "CHANNEL",
		Args: []string{f.Item},
	}
}

// Expr returns the function as an expression.
func (f ChannelFunc) Expr() ExprDef {
	return ExprDef{Content: f.Func().String()}
}

// QuotedContent returns true, which indicates that the function call should
// be wrapped in double quotes when compared.
func (f ChannelFunc) QuotedContent() bool {
	return true
}

// Validate returns an error if the item contains invalid characters.
func (f ChannelFunc) Validate() error {
	return errorIfAny("channel item", f.Item, invalidArgChars)
}
```

Create `dialplan/predial.go`:

```go
package dialplan

import (
	"fmt"
	"strconv"
	"strings"
)

// PredialTarget is a subroutine that an application runs on each callee
// channel before dialing it, such as the target of Page's b() option.
//
// It is written as context^extension^priority(arg1^arg2), the form Asterisk
// uses inside an option, where commas would end the option early.
type PredialTarget struct {
	Context   string
	Extension string
	Priority  int // Should be >= 1
	Args      []string
}

// String returns the target in option form.
func (t PredialTarget) String() string {
	var b strings.Builder
	b.WriteString(t.Context)
	b.WriteString("^")
	b.WriteString(t.Extension)
	b.WriteString("^")
	b.WriteString(strconv.Itoa(t.Priority))
	if len(t.Args) > 0 {
		b.WriteString("(")
		b.WriteString(strings.Join(t.Args, "^"))
		b.WriteString(")")
	}
	return b.String()
}

// Validate returns an error if the priority is below one, or if any part
// contains a character that would change the meaning of the option.
func (t PredialTarget) Validate() error {
	if t.Priority < 1 {
		return fmt.Errorf("predial priority must be at least 1, got %d", t.Priority)
	}
	if err := errorIfAny("predial context", t.Context, invalidArgChars+"^"); err != nil {
		return err
	}
	if err := errorIfAny("predial extension", t.Extension, invalidArgChars+"^"); err != nil {
		return err
	}
	for _, arg := range t.Args {
		if err := errorIfAny("predial argument", arg, invalidArgChars+"^"); err != nil {
			return err
		}
	}
	return nil
}
```

Replace the whole of `dialplan/page.go` with (the struct's comment column moves because gofmt realigns it for the longer `*PredialTarget` type; the only behavior change is the `Predial` field):

```go
package dialplan

import (
	"strconv"
	"strings"
)

// http://www.asteriskdocs.org/en/3rd_Edition/asterisk-book-html-chunk/AdditionalConfig_id257169.html

// Page returns a dialplan application that pages a set of devices.
func Page(recipients ...Device) PageApp {
	return PageApp{
		Recipients: recipients,
	}
}

// PageApp is a dialplan application that pages a set of devices.
//
// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Application_Page
type PageApp struct {
	Recipients       []Device
	FullDuplex       bool           //  d: Enable two-way audio for all recipients
	Forward          bool           // !i: Forward the page to forwarded recipients
	Quiet            bool           //  q: Don't present a paging beep to recipients
	Record           bool           //  r: Record the page to a file
	InUse            bool           // !s: Include recipients with a phone state other than NOT_INUSE
	Announcement     string         //  A: Play an announcement to all recipients
	AnnounceToCaller bool           //  x: Play an announcement to the caller as well
	Timeout          int            // Number of seconds to attempt a call before giving up
	Predial          *PredialTarget // b: Run a subroutine on each callee channel before dialing it
}

// App returns the assembled application call.
func (app PageApp) App() AppCall {
	devices := make([]string, 0, len(app.Recipients))
	for _, device := range app.Recipients {
		devices = append(devices, device.String())
	}
	call := AppCall{
		Name: "Page",
		Args: []string{strings.Join(devices, "&")},
	}
	options := app.Options()
	if len(options) == 0 && app.Timeout <= 0 {
		return call
	}
	call.Args = append(call.Args, options)
	if app.Timeout > 0 {
		call.Args = append(call.Args, strconv.Itoa(app.Timeout))
	}
	return call
}

// Options returns the options string for the page.
func (app PageApp) Options() string {
	var opts string
	if app.FullDuplex {
		opts += "d"
	}
	if !app.Forward {
		opts += "i"
	}
	if app.Quiet {
		opts += "q"
	}
	if app.Record {
		opts += "r"
	}
	if !app.InUse {
		opts += "s"
	}
	if app.Announcement != "" {
		opts += "A(" + app.Announcement + ")"
		if app.AnnounceToCaller {
			opts += "n"
		}
	}
	if app.Predial != nil {
		opts += "b(" + app.Predial.String() + ")"
	}
	return opts
}

// Validate returns an error if any recipient or the announcement is invalid.
func (app PageApp) Validate() error {
	for _, device := range app.Recipients {
		if err := device.Validate(); err != nil {
			return err
		}
	}
	if err := errorIfAny("page announcement", app.Announcement, invalidArgChars); err != nil {
		return err
	}
	if app.Predial != nil {
		return app.Predial.Validate()
	}
	return nil
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./dialplan -v -run 'TestPerServer|TestPagePredial|TestSubroutine' && gofmt -l . && go vet ./dialplan && go test ./dialplan`
Expected: `TestPerServerApps`, `TestPerServerAppsInvalid`, `TestPagePredialAndHeaders`, `TestPagePredialInvalid`, `TestSubroutine` and `TestSubroutineInvalid` PASS. The existing dialplan tests still pass, including `Example`, and `TestPagePredialAndHeaders` shows a `Page` with no predial still renders `Page(PJSIP/fred,is)`.

- [ ] **Step 5: Commit**

```bash
git add dialplan/busy.go dialplan/voicemail.go dialplan/cat.go dialplan/dialexpr.go dialplan/subroutine.go dialplan/return.go dialplan/pjsipheader.go dialplan/channel.go dialplan/predial.go dialplan/page.go dialplan/perserver_test.go
git commit -F - <<'EOF'
feat(dialplan): add the pieces the per-server dialplan needs

The per-server ring builds its dial string in a channel variable at run
time, so it needs to append to a variable, read part of it, and dial an
expression rather than a fixed device list. It calls its hand-written
subroutines with Gosub and arguments, since Asterisk 21 removed Macro.
Paging on PJSIP adds its
Alert-Info header in a predial subroutine, so Page gains a predial target
and the dialplan gains Return, PJSIP_HEADER and CHANNEL.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 7: `astgen.Place`

**Files:**
- Create: `astgen/names.go`, `astgen/placement.go`
- Create (test fixture): `astgen/multiserver_test.go`
- Test: `astgen/placement_test.go`

**Interfaces:**
- Consumes: `astorg.DataSet`, `astorg.Lookup.LocationServer`, `astorg.Lookup.ServerByID`, `lineUsername` (existing, `astgen/username.go`).
- Produces (later tasks rely on every name here):
  - `func Place(data *astorg.DataSet) (Placement, []Problem, error)`; errors `ErrNoHub`, `ErrMultipleHubs`
  - `type Placement struct { Hub string; Servers map[string]astorg.Server; Phones, Hosts, Homes, Owners map[string]string; Devices map[string][]string; Remote map[string]map[string][]string; TrunkPairs []TrunkPair }`
  - methods `ServerIDs() []string`, `Peers(serverID string) []astorg.Server`, `HostedDevices(serverID string) []string`, `HomedExtensions(serverID string) []string`, `Credentials(serverID string) []Credential`
  - `func PairID(a, b string) string`
  - `type Problem struct { Rule, Kind, ID, Detail string }`, rules `RuleLocationNoServer`, `RuleServerMissing`, `RuleExtensionDuplicate`, `RuleDeviceConflict`
  - `type TrunkPair struct { ID, A, B string }`, `type Credential struct { Kind, ID, Section string }`
  - unexported, in `names.go`: `TrunkContext = "from-trunk"`, `TrunkTransport = "transport-tcp"`, `trunkPort`, `trunkName(peer)`, `trunkAuthName(peer)`, `publishName(peer)`, `publishMWIName(peer)`, and the exported `PageAlertContext = "astgen-page-alert"`, `PageRelayContext = "astgen-page-relay"` (used from Task 10); in `placement.go`: `sortedUnique([]string) []string`
  - test helpers: `multiServerData() *astorg.DataSet` and `place(t, data) (astgen.Placement, []astgen.Problem)`

Device names are pjsip endpoint names: a desk phone's line name (`lineUsername`, the assigned person's or role's username, else the MAC) and a softphone's username. An extension rings the device named after its owner's username, if one is hosted, and its owner's hosted softphones. This matches the device list the old `PhoneExtensions` builds.

- [ ] **Step 1: Write the fixture**

Create `astgen/multiserver_test.go`:

```go
package astgen_test

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astorg/astorgvm"
)

// multiServerData returns a data set with two office servers and a hub.
//
// Server IDs sort differently than server names, so a test that orders by
// name instead of ID fails. Quarry and Den have servers; Remote is an
// office with no server, so extensions there live on the hub.
//
//   - fred (100) has a desk phone at Quarry and a softphone, so his
//     extension lives on qry and rings one device on the hub.
//   - barney (101) has a desk phone at Den.
//   - wilma (102) works at Remote and has only a softphone, so her
//     extension and her one device both live on the hub.
//   - The lobby role (200) has a desk phone at Den and a mailbox.
//   - The yard role (201) is at Remote with only a softphone.
//   - A spare phone at Quarry is assigned to nobody.
//   - The All Call paging group (500) covers Quarry and Den, so none of its
//     members are on the hub.
func multiServerData() *astorg.DataSet {
	return &astorg.DataSet{
		Servers: astorg.ServerList{
			{ID: "qry", Name: "Quarry", Address: "qry-phone"},
			{ID: "den", Name: "Rock Den", Address: "den-phone"},
			{ID: "hub", Name: "Hub", Hub: true, Address: "hub-phone"},
		},
		Locations: astorg.LocationList{
			{Name: "Quarry", Abbreviation: "QRY", Server: "qry", CallerID: "5555550100", AreaCode: "555", PagingGroups: []string{"500"}},
			{Name: "Den", Abbreviation: "DEN", Server: "den", CallerID: "5555550200", AreaCode: "555", PagingGroups: []string{"500"}},
			{Name: "Remote", Abbreviation: "RMT", CallerID: "5555550300", AreaCode: "555"},
		},
		People: astorg.PersonList{
			{
				Username:        "fred",
				FullName:        "Fred Flintstone",
				Extension:       "100",
				Location:        "Quarry",
				VoicemailCode:   "1234",
				VoicemailAccess: astorgvm.PhoneAndEmail,
				EmailAddresses:  []astorg.Email{{Address: "fred@example.com", Primary: true}},
				Phones:          []string{"000fd3000001"},
				Softphones:      []string{"fred.soft"},
			},
			{
				Username:      "barney",
				FullName:      "Barney Rubble",
				Extension:     "101",
				Location:      "Den",
				VoicemailCode: "4321",
				Phones:        []string{"000fd3000004"},
			},
			{
				Username:   "wilma",
				FullName:   "Wilma Flintstone",
				Extension:  "102",
				Location:   "Remote",
				Softphones: []string{"wilma.soft"},
			},
		},
		PhoneRoles: astorg.PhoneRoleList{
			{
				Username:      "lobby",
				DisplayName:   "Lobby",
				Extension:     "200",
				Location:      "Den",
				MailboxNumber: "300",
				Phones:        []string{"000fd3000002"},
			},
			{
				Username:    "yard",
				DisplayName: "Yard",
				Extension:   "201",
				Location:    "Remote",
				Softphones:  []string{"yard.soft"},
			},
		},
		Phones: astorg.PhoneList{
			{MAC: "000fd3000001", Location: "Quarry", Secret: "fred-secret"},
			{MAC: "000fd3000002", Location: "Den", Secret: "lobby-secret"},
			{MAC: "000fd3000003", Location: "Quarry", Secret: "spare-secret"},
			{MAC: "000fd3000004", Location: "Den", Secret: "barney-secret"},
		},
		Softphones: astorg.SoftphoneList{
			{Username: "fred.soft", Location: "Quarry", Secret: "fred-soft-secret"},
			{Username: "wilma.soft", Location: "Remote", Secret: "wilma-soft-secret"},
			{Username: "yard.soft", Location: "Remote", Secret: "yard-soft-secret"},
		},
		PagingGroups: astorg.PagingGroupList{
			{Name: "All Call", Extension: "500", Alert: "intercom"},
		},
		Mailboxes: astorg.MailboxList{
			{Number: "300", Name: "Sales", AccessCode: "3000", AccessMode: astorgvm.Email, Email: "sales@example.com"},
		},
	}
}
```

- [ ] **Step 2: Write the failing tests**

Create `astgen/placement_test.go`:

```go
package astgen_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
)

// place runs Place on data and fails the test on an error.
func place(t *testing.T, data *astorg.DataSet) (astgen.Placement, []astgen.Problem) {
	t.Helper()
	p, problems, err := astgen.Place(data)
	if err != nil {
		t.Fatalf("Place() error: %v", err)
	}
	return p, problems
}

func TestPlaceHubCount(t *testing.T) {
	t.Run("no hub", func(t *testing.T) {
		data := multiServerData()
		data.Servers[2].Hub = false
		if _, _, err := astgen.Place(data); !errors.Is(err, astgen.ErrNoHub) {
			t.Errorf("Place() error = %v, want ErrNoHub", err)
		}
	})
	t.Run("two hubs", func(t *testing.T) {
		data := multiServerData()
		data.Servers[0].Hub = true
		if _, _, err := astgen.Place(data); !errors.Is(err, astgen.ErrMultipleHubs) {
			t.Errorf("Place() error = %v, want ErrMultipleHubs", err)
		}
	})
	t.Run("no servers", func(t *testing.T) {
		if _, _, err := astgen.Place(testData()); !errors.Is(err, astgen.ErrNoHub) {
			t.Errorf("Place() error = %v, want ErrNoHub", err)
		}
	})
}

func TestPlaceHosts(t *testing.T) {
	p, problems := place(t, multiServerData())
	if len(problems) != 0 {
		t.Fatalf("Place() problems: %v", problems)
	}
	if p.Hub != "hub" {
		t.Errorf("Hub = %q, want hub", p.Hub)
	}
	wantHosts := map[string]string{
		"fred":         "qry",
		"lobby":        "den",
		"000fd3000003": "qry",
		"barney":       "den",
		"fred.soft":    "hub", // A softphone is hosted by the hub even when its location has a server
		"wilma.soft":   "hub",
		"yard.soft":    "hub",
	}
	if !reflect.DeepEqual(p.Hosts, wantHosts) {
		t.Errorf("Hosts = %v, want %v", p.Hosts, wantHosts)
	}
	wantPhones := map[string]string{
		"000fd3000001": "qry",
		"000fd3000002": "den",
		"000fd3000003": "qry",
		"000fd3000004": "den",
	}
	if !reflect.DeepEqual(p.Phones, wantPhones) {
		t.Errorf("Phones = %v, want %v", p.Phones, wantPhones)
	}
	if got, want := p.HostedDevices("qry"), []string{"000fd3000003", "fred"}; !reflect.DeepEqual(got, want) {
		t.Errorf("HostedDevices(qry) = %v, want %v", got, want)
	}
}

func TestPlaceHomes(t *testing.T) {
	p, _ := place(t, multiServerData())
	wantHomes := map[string]string{
		"100": "qry",
		"101": "den",
		"102": "hub", // Remote has no server
		"200": "den",
		"201": "hub",
	}
	if !reflect.DeepEqual(p.Homes, wantHomes) {
		t.Errorf("Homes = %v, want %v", p.Homes, wantHomes)
	}
	wantOwners := map[string]string{
		"100": "fred",
		"101": "barney",
		"102": "wilma",
		"200": "lobby",
		"201": "yard",
	}
	if !reflect.DeepEqual(p.Owners, wantOwners) {
		t.Errorf("Owners = %v, want %v", p.Owners, wantOwners)
	}
	if got, want := p.HomedExtensions("hub"), []string{"102", "201"}; !reflect.DeepEqual(got, want) {
		t.Errorf("HomedExtensions(hub) = %v, want %v", got, want)
	}
}

func TestPlaceRemoteDevices(t *testing.T) {
	p, _ := place(t, multiServerData())
	wantDevices := map[string][]string{
		"100": {"fred", "fred.soft"},
		"101": {"barney"},
		"102": {"wilma.soft"},
		"200": {"lobby"},
		"201": {"yard.soft"},
	}
	if !reflect.DeepEqual(p.Devices, wantDevices) {
		t.Errorf("Devices = %v, want %v", p.Devices, wantDevices)
	}
	wantRemote := map[string]map[string][]string{
		"100": {"hub": {"fred.soft"}},
	}
	if !reflect.DeepEqual(p.Remote, wantRemote) {
		t.Errorf("Remote = %v, want %v", p.Remote, wantRemote)
	}
}

// TestPlaceRoleAtOfficeWithoutServer checks that a role at an office with no
// server lives on the hub with its softphone, so it has no remote devices.
func TestPlaceRoleAtOfficeWithoutServer(t *testing.T) {
	p, _ := place(t, multiServerData())
	if got := p.Homes["201"]; got != "hub" {
		t.Errorf("Homes[201] = %q, want hub", got)
	}
	if got, want := p.Devices["201"], []string{"yard.soft"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Devices[201] = %v, want %v", got, want)
	}
	if got := p.Remote["201"]; got != nil {
		t.Errorf("Remote[201] = %v, want none", got)
	}
}

// TestPlaceSortsByID checks that servers, peers and trunk pairs are ordered
// by server ID. The fixture's names sort Hub, Quarry, Rock Den, which is not
// the ID order den, hub, qry.
func TestPlaceSortsByID(t *testing.T) {
	p, _ := place(t, multiServerData())
	if got, want := p.ServerIDs(), []string{"den", "hub", "qry"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ServerIDs() = %v, want %v", got, want)
	}
	var peers []string
	for _, peer := range p.Peers("qry") {
		peers = append(peers, peer.ID)
	}
	if want := []string{"den", "hub"}; !reflect.DeepEqual(peers, want) {
		t.Errorf("Peers(qry) = %v, want %v", peers, want)
	}
	wantPairs := []astgen.TrunkPair{
		{ID: "den~hub", A: "den", B: "hub"},
		{ID: "den~qry", A: "den", B: "qry"},
		{ID: "hub~qry", A: "hub", B: "qry"},
	}
	if !reflect.DeepEqual(p.TrunkPairs, wantPairs) {
		t.Errorf("TrunkPairs = %v, want %v", p.TrunkPairs, wantPairs)
	}
	if got := astgen.PairID("qry", "den"); got != "den~qry" {
		t.Errorf("PairID(qry, den) = %q, want den~qry", got)
	}
}

func TestPlaceCredentials(t *testing.T) {
	p, _ := place(t, multiServerData())
	want := []astgen.Credential{
		{Kind: "trunk", ID: "den~hub", Section: "auth-trunk-den"},
		{Kind: "trunk", ID: "hub~qry", Section: "auth-trunk-qry"},
	}
	if got := p.Credentials("hub"); !reflect.DeepEqual(got, want) {
		t.Errorf("Credentials(hub) = %v, want %v", got, want)
	}
}

func TestPlaceProblems(t *testing.T) {
	tests := []struct {
		name   string
		modify func(data *astorg.DataSet)
		want   []astgen.Problem
		check  func(t *testing.T, p astgen.Placement)
	}{
		{
			name: "desk phone at an office with no server",
			modify: func(data *astorg.DataSet) {
				data.Phones[2].Location = "Remote"
			},
			want: []astgen.Problem{
				{Rule: astgen.RuleLocationNoServer, Kind: "phone", ID: "000fd3000003", Detail: `location "Remote" has no server`},
			},
			check: func(t *testing.T, p astgen.Placement) {
				if _, ok := p.Hosts["000fd3000003"]; ok {
					t.Error("the phone is still hosted")
				}
				if _, ok := p.Phones["000fd3000003"]; ok {
					t.Error("the phone is still placed")
				}
			},
		},
		{
			name: "location names a missing server",
			modify: func(data *astorg.DataSet) {
				data.Locations[1].Server = "gone"
			},
			want: []astgen.Problem{
				{Rule: astgen.RuleServerMissing, Kind: "phone", ID: "000fd3000002", Detail: `location "Den" names server "gone", which is not in the data set`},
				{Rule: astgen.RuleServerMissing, Kind: "phone", ID: "000fd3000004", Detail: `location "Den" names server "gone", which is not in the data set`},
				{Rule: astgen.RuleServerMissing, Kind: "person", ID: "barney", Detail: `location "Den" names server "gone", which is not in the data set`},
				{Rule: astgen.RuleServerMissing, Kind: "role", ID: "lobby", Detail: `location "Den" names server "gone", which is not in the data set`},
			},
			check: func(t *testing.T, p astgen.Placement) {
				if _, ok := p.Homes["101"]; ok {
					t.Error("extension 101 is still placed")
				}
			},
		},
		{
			name: "two people claim one extension",
			modify: func(data *astorg.DataSet) {
				data.People = append(data.People, astorg.Person{Username: "dino", FullName: "Dino", Extension: "100", Location: "Den"})
			},
			want: []astgen.Problem{
				{Rule: astgen.RuleExtensionDuplicate, Kind: "person", ID: "dino", Detail: `extension 100 is already held by "fred"`},
			},
			check: func(t *testing.T, p astgen.Placement) {
				if p.Owners["100"] != "fred" || p.Homes["100"] != "qry" {
					t.Errorf("extension 100 = %q on %q, want fred on qry", p.Owners["100"], p.Homes["100"])
				}
			},
		},
		{
			name: "desk phones sharing a line on two servers",
			modify: func(data *astorg.DataSet) {
				data.Phones = append(data.Phones, astorg.Phone{MAC: "000fd3000005", Location: "Den", Secret: "fred-den-secret"})
				data.People[0].Phones = append(data.People[0].Phones, "000fd3000005")
			},
			want: []astgen.Problem{
				{Rule: astgen.RuleDeviceConflict, Kind: "phone", ID: "000fd3000005", Detail: `line "fred" is already hosted on server "qry", not "den"`},
			},
			check: func(t *testing.T, p astgen.Placement) {
				if p.Hosts["fred"] != "qry" {
					t.Errorf("Hosts[fred] = %q, want qry", p.Hosts["fred"])
				}
				if _, ok := p.Phones["000fd3000005"]; ok {
					t.Error("the conflicting phone is still placed")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := multiServerData()
			tt.modify(data)
			p, problems := place(t, data)
			if !reflect.DeepEqual(problems, tt.want) {
				t.Errorf("Place() problems =\n%v\nwant:\n%v", problems, tt.want)
			}
			tt.check(t, p)
		})
	}
}
```

- [ ] **Step 3: Run the tests to verify they fail**

Run: `go test ./astgen -run TestPlace -v`
Expected: FAIL to compile with `undefined: astgen.Place`.

- [ ] **Step 4: Write the shared names**

Create `astgen/names.go`:

```go
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

// PageAlertContext is the dialplan context of the subroutine that adds a
// paging group's Alert-Info header to each paged PJSIP device.
// PageAlertSection renders it.
const PageAlertContext = "astgen-page-alert"

// PageRelayContext is the dialplan context that holds one relay extension
// for each paging group and each other server with members of it.
// ServerPageRelayExtensions renders it.
const PageRelayContext = "astgen-page-relay"
```

- [ ] **Step 5: Write `Place`**

Create `astgen/placement.go`:

```go
package astgen

import (
	"errors"
	"fmt"
	"sort"

	"github.com/scjalliance/astconf/astorg"
)

// Errors returned by Place. Nothing can route without exactly one hub.
var (
	ErrNoHub        = errors.New("astgen: the data set has no hub server")
	ErrMultipleHubs = errors.New("astgen: the data set has more than one hub server")
)

// Problem rules reported by Place.
const (
	RuleLocationNoServer   = "location-no-server"  // A desk phone at a location with no server
	RuleServerMissing      = "server-missing"      // A location names a server that is not in the data set
	RuleExtensionDuplicate = "extension-duplicate" // Two people or roles claim the same extension
	RuleDeviceConflict     = "device-conflict"     // Desk phones that share a line name sit on different servers
)

// Problem describes a record that Place left out of the placement. No
// field ever carries a secret.
type Problem struct {
	Rule   string // One of the Rule constants
	Kind   string // "phone", "person" or "role"
	ID     string // The MAC address, or the person's or role's username
	Detail string
}

// String returns a one-line description of the problem.
func (p Problem) String() string {
	return fmt.Sprintf("%s: %s %s: %s", p.Rule, p.Kind, p.ID, p.Detail)
}

// TrunkPair is an unordered pair of servers joined by a trunk.
type TrunkPair struct {
	ID string // The two server IDs, sorted, joined by "~"
	A  string // The server ID that sorts first
	B  string // The server ID that sorts second
}

// Credential names a secret that the node agent writes into a server's
// local secrets file. astgen never renders the secret itself.
type Credential struct {
	Kind    string // "trunk"
	ID      string // The trunk pair ID, from which both ends derive the same secret
	Section string // The pjsip auth section the secret is written to
}

// Placement records which server hosts each device and which server owns
// each extension. It is computed once per data set by Place and shared by
// every per-server generator.
//
// Device names are pjsip endpoint names: the line name of a desk phone and
// the username of a softphone.
type Placement struct {
	Hub        string                         // ID of the hub server
	Servers    map[string]astorg.Server       // Every server, keyed by ID
	Phones     map[string]string              // Desk phone MAC address to host server ID
	Hosts      map[string]string              // Device name to host server ID
	Homes      map[string]string              // Extension to home server ID
	Owners     map[string]string              // Extension to the username of the person or role that holds it
	Devices    map[string][]string            // Extension to the names of the devices it rings, sorted
	Remote     map[string]map[string][]string // Extension to host server ID to the devices it rings there, sorted, for every host other than its home
	TrunkPairs []TrunkPair                    // Every unordered pair of servers, sorted by ID
}

// Place decides which server hosts each device and which server owns each
// extension.
//
// A desk phone is hosted by its location's server. Every softphone is
// hosted by the hub. A person's or role's extension lives on its location's
// server, or on the hub when that location has no server.
//
// Place returns an error only when the data set does not have exactly one
// hub. Any other record that cannot be placed is reported as a Problem and
// left out of the placement.
func Place(data *astorg.DataSet) (Placement, []Problem, error) {
	lookup := data.Lookup()

	var hubs []string
	for id, server := range lookup.ServerByID {
		if server.Hub {
			hubs = append(hubs, id)
		}
	}
	switch {
	case len(hubs) == 0:
		return Placement{}, nil, ErrNoHub
	case len(hubs) > 1:
		return Placement{}, nil, ErrMultipleHubs
	}

	p := Placement{
		Hub:     hubs[0],
		Servers: make(map[string]astorg.Server, len(lookup.ServerByID)),
		Phones:  make(map[string]string),
		Hosts:   make(map[string]string),
		Homes:   make(map[string]string),
		Owners:  make(map[string]string),
		Devices: make(map[string][]string),
		Remote:  make(map[string]map[string][]string),
	}
	for id, server := range lookup.ServerByID {
		p.Servers[id] = server
	}

	var problems []Problem

	// Step 1: Host each desk phone on its location's server
	for _, phone := range data.Phones {
		name := lineUsername(phone.MAC, lookup)
		if name == "" {
			continue
		}
		id, _, ok := lookup.LocationServer(phone.Location)
		switch {
		case id == "":
			problems = append(problems, Problem{
				Rule:   RuleLocationNoServer,
				Kind:   "phone",
				ID:     phone.MAC,
				Detail: fmt.Sprintf("location %q has no server", phone.Location),
			})
			continue
		case !ok:
			problems = append(problems, Problem{
				Rule:   RuleServerMissing,
				Kind:   "phone",
				ID:     phone.MAC,
				Detail: fmt.Sprintf("location %q names server %q, which is not in the data set", phone.Location, id),
			})
			continue
		}
		if host, found := p.Hosts[name]; found && host != id {
			problems = append(problems, Problem{
				Rule:   RuleDeviceConflict,
				Kind:   "phone",
				ID:     phone.MAC,
				Detail: fmt.Sprintf("line %q is already hosted on server %q, not %q", name, host, id),
			})
			continue
		}
		p.Phones[phone.MAC] = id
		p.Hosts[name] = id
	}

	// Step 2: Host every softphone on the hub
	for _, softphone := range data.Softphones {
		if softphone.Username == "" {
			continue
		}
		if _, found := p.Hosts[softphone.Username]; found {
			continue
		}
		p.Hosts[softphone.Username] = p.Hub
	}

	// Step 3: Home each extension, people first and then roles
	type owner struct {
		kind       string
		username   string
		extension  string
		location   string
		softphones []string
	}
	var owners []owner
	for _, person := range data.People {
		owners = append(owners, owner{"person", person.Username, person.Extension, person.Location, person.Softphones})
	}
	for _, role := range data.PhoneRoles {
		owners = append(owners, owner{"role", role.Username, role.Extension, role.Location, role.Softphones})
	}

	for _, o := range owners {
		if o.extension == "" || o.username == "" {
			continue
		}
		if holder, found := p.Owners[o.extension]; found {
			problems = append(problems, Problem{
				Rule:   RuleExtensionDuplicate,
				Kind:   o.kind,
				ID:     o.username,
				Detail: fmt.Sprintf("extension %s is already held by %q", o.extension, holder),
			})
			continue
		}
		home := p.Hub
		id, _, ok := lookup.LocationServer(o.location)
		switch {
		case id != "" && !ok:
			problems = append(problems, Problem{
				Rule:   RuleServerMissing,
				Kind:   o.kind,
				ID:     o.username,
				Detail: fmt.Sprintf("location %q names server %q, which is not in the data set", o.location, id),
			})
			continue
		case id != "":
			home = id
		}
		p.Homes[o.extension] = home
		p.Owners[o.extension] = o.username

		// The devices an extension rings are its desk phone line, which is
		// named after the username, and its softphones.
		var devices []string
		for _, name := range append([]string{o.username}, o.softphones...) {
			if _, hosted := p.Hosts[name]; hosted {
				devices = append(devices, name)
			}
		}
		devices = sortedUnique(devices)
		p.Devices[o.extension] = devices

		for _, name := range devices {
			host := p.Hosts[name]
			if host == home {
				continue
			}
			if p.Remote[o.extension] == nil {
				p.Remote[o.extension] = make(map[string][]string)
			}
			p.Remote[o.extension][host] = append(p.Remote[o.extension][host], name)
		}
	}

	// Step 4: Pair every two servers
	ids := p.ServerIDs()
	for i := range ids {
		for j := i + 1; j < len(ids); j++ {
			p.TrunkPairs = append(p.TrunkPairs, TrunkPair{
				ID: PairID(ids[i], ids[j]),
				A:  ids[i],
				B:  ids[j],
			})
		}
	}

	return p, problems, nil
}

// ServerIDs returns the ID of every server, sorted.
func (p Placement) ServerIDs() []string {
	ids := make([]string, 0, len(p.Servers))
	for id := range p.Servers {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Peers returns every server other than serverID, sorted by ID.
func (p Placement) Peers(serverID string) []astorg.Server {
	var peers []astorg.Server
	for _, id := range p.ServerIDs() {
		if id != serverID {
			peers = append(peers, p.Servers[id])
		}
	}
	return peers
}

// HostedDevices returns the names of the devices that serverID hosts,
// sorted.
func (p Placement) HostedDevices(serverID string) []string {
	var names []string
	for name, host := range p.Hosts {
		if host == serverID {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// HomedExtensions returns the extensions that live on serverID, sorted.
func (p Placement) HomedExtensions(serverID string) []string {
	var extensions []string
	for extension, home := range p.Homes {
		if home == serverID {
			extensions = append(extensions, extension)
		}
	}
	sort.Strings(extensions)
	return extensions
}

// Credentials returns the trunk credentials that serverID's node agent
// writes, one per peer, sorted by pair ID.
func (p Placement) Credentials(serverID string) []Credential {
	var creds []Credential
	for _, pair := range p.TrunkPairs {
		var peer string
		switch serverID {
		case pair.A:
			peer = pair.B
		case pair.B:
			peer = pair.A
		default:
			continue
		}
		creds = append(creds, Credential{
			Kind:    "trunk",
			ID:      pair.ID,
			Section: trunkAuthName(peer),
		})
	}
	return creds
}

// PairID returns the trunk pair ID for two servers: their IDs, sorted,
// joined by "~".
func PairID(a, b string) string {
	if b < a {
		a, b = b, a
	}
	return a + "~" + b
}

func sortedUnique(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	sorted := append([]string(nil), values...)
	sort.Strings(sorted)
	out := sorted[:1]
	for _, v := range sorted[1:] {
		if v != out[len(out)-1] {
			out = append(out, v)
		}
	}
	return out
}
```

- [ ] **Step 6: Run the tests to verify they pass**

Run: `go test ./astgen -run TestPlace -v && gofmt -l . && go vet ./astgen && go test ./astgen`
Expected: every `TestPlace*` test and subtest PASSes, including `TestPlaceRoleAtOfficeWithoutServer` and `TestPlaceSortsByID` (Review Focus 2 and 5). The existing astgen tests still pass.

- [ ] **Step 7: Commit**

```bash
git add astgen/names.go astgen/placement.go astgen/multiserver_test.go astgen/placement_test.go
git commit -F - <<'EOF'
feat(astgen): add Place for multi-server rendering

Place decides once per data set which server hosts each device and owns
each extension, lists the trunk pairs and their credentials, and reports
records it cannot place as problems instead of failing.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 8: Per-server device sections

**Files:**
- Create: `astgen/server.go`
- Test: `astgen/server_test.go`

**Interfaces:**
- Consumes: `Placement.Hosts`, `Placement.Phones`, `Placement.Hub` (Task 7); `Endpoints`, `AORs`, `Lines`, `Phones`, `Mailboxes`, `aorName` (existing); `place`, `multiServerData`, `render` (test helpers).
- Produces:
  - `func ServerEndpoints(data *astorg.DataSet, p Placement, serverID string, base pjsip.Endpoint, context string) []pjsip.Endpoint`
  - `func ServerAORs(data *astorg.DataSet, p Placement, serverID string, base pjsip.AOR) []pjsip.AOR`
  - `func ServerLines(data *astorg.DataSet, p Placement, serverID string) []dpma.Line`
  - `func ServerPhones(data *astorg.DataSet, p Placement, serverID string, base dpma.Phone, contactsURL string) []dpma.Phone`
  - `func ServerMailboxes(data *astorg.DataSet, p Placement, serverID string, context string) voicemail.Section`

Each filters the old generator's output, so a device's section is byte for byte what it is today. That is what makes Task 13's equivalence hold.

- [ ] **Step 1: Write the failing tests**

Create `astgen/server_test.go`:

```go
package astgen_test

import (
	"reflect"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/dpma"
	"github.com/scjalliance/astconf/pjsip"
)

// TestServerDeviceSections checks which devices each server renders
// endpoints, aors, lines and phones for. The rendered text of each entry is
// covered by the golden files.
func TestServerDeviceSections(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)

	tests := []struct {
		server    string
		endpoints []string
		aors      []string
		lines     []string
		phones    []string
	}{
		{
			server:    "den",
			endpoints: []string{"lobby", "barney"},
			aors:      []string{"lobby-aor", "barney-aor"},
			lines:     []string{"lobby", "barney"},
			phones:    []string{"000fd3000002", "000fd3000004"},
		},
		{
			server:    "hub",
			endpoints: []string{"fred.soft", "wilma.soft", "yard.soft"},
			aors:      []string{"fred.soft-aor", "wilma.soft-aor", "yard.soft-aor"},
		},
		{
			server:    "qry",
			endpoints: []string{"fred", "000fd3000003"},
			aors:      []string{"fred-aor", "000fd3000003-aor"},
			lines:     []string{"fred", "000fd3000003"},
			phones:    []string{"000fd3000001", "000fd3000003"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.server, func(t *testing.T) {
			var endpoints, aors, lines, phones []string
			for _, e := range astgen.ServerEndpoints(data, p, tt.server, pjsip.Endpoint{}, "default") {
				endpoints = append(endpoints, e.Name)
			}
			for _, a := range astgen.ServerAORs(data, p, tt.server, pjsip.AOR{}) {
				aors = append(aors, a.Name)
			}
			for _, l := range astgen.ServerLines(data, p, tt.server) {
				lines = append(lines, l.Name)
			}
			for _, ph := range astgen.ServerPhones(data, p, tt.server, dpma.Phone{}, "") {
				phones = append(phones, ph.MAC)
			}
			if !reflect.DeepEqual(endpoints, tt.endpoints) {
				t.Errorf("ServerEndpoints() = %v, want %v", endpoints, tt.endpoints)
			}
			if !reflect.DeepEqual(aors, tt.aors) {
				t.Errorf("ServerAORs() = %v, want %v", aors, tt.aors)
			}
			if !reflect.DeepEqual(lines, tt.lines) {
				t.Errorf("ServerLines() = %v, want %v", lines, tt.lines)
			}
			if !reflect.DeepEqual(phones, tt.phones) {
				t.Errorf("ServerPhones() = %v, want %v", phones, tt.phones)
			}
		})
	}
}

// TestServerMailboxes checks that voicemail renders on the hub alone.
func TestServerMailboxes(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)

	office := astgen.ServerMailboxes(data, p, "qry", "default")
	if got, want := render(t, &office), "[default]\n"; got != want {
		t.Errorf("ServerMailboxes(qry) rendered:\n%s\nwant:\n%s", got, want)
	}

	hub := astgen.ServerMailboxes(data, p, "hub", "default")
	got := render(t, &hub)
	want := "[default]\n" +
		"300 => -3000,Sales,sales@example.com,,attach=yes|saycid=yes|delete=yes\n" +
		"100 => -1234,fred,fred@example.com,,attach=yes|saycid=yes|delete=no\n" +
		"101 => -4321,barney,,,attach=yes|saycid=yes|delete=no\n"
	if got != want {
		t.Errorf("ServerMailboxes(hub) rendered:\n%s\nwant:\n%s", got, want)
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./astgen -run 'TestServerDeviceSections|TestServerMailboxes' -v`
Expected: FAIL to compile with `undefined: astgen.ServerEndpoints`.

- [ ] **Step 3: Write the implementation**

Create `astgen/server.go`:

```go
package astgen

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dpma"
	"github.com/scjalliance/astconf/pjsip"
	"github.com/scjalliance/astconf/voicemail"
)

// hostedSet returns the set of device names that serverID hosts.
func hostedSet(p Placement, serverID string) map[string]bool {
	set := make(map[string]bool)
	for name, host := range p.Hosts {
		if host == serverID {
			set[name] = true
		}
	}
	return set
}

// ServerEndpoints generates PJSIP endpoint entries for the devices that
// serverID hosts. Each entry is the one Endpoints generates for the same
// device, in the same order.
func ServerEndpoints(data *astorg.DataSet, p Placement, serverID string, base pjsip.Endpoint, context string) []pjsip.Endpoint {
	hosted := hostedSet(p, serverID)
	var out []pjsip.Endpoint
	for _, endpoint := range Endpoints(data, base, context) {
		if hosted[endpoint.Name] {
			out = append(out, endpoint)
		}
	}
	return out
}

// ServerAORs generates PJSIP aor entries for the devices that serverID
// hosts. Each entry is the one AORs generates for the same device, in the
// same order.
func ServerAORs(data *astorg.DataSet, p Placement, serverID string, base pjsip.AOR) []pjsip.AOR {
	hosted := make(map[string]bool)
	for name := range hostedSet(p, serverID) {
		hosted[aorName(name)] = true
	}
	var out []pjsip.AOR
	for _, aor := range AORs(data, base) {
		if hosted[aor.Name] {
			out = append(out, aor)
		}
	}
	return out
}

// ServerLines generates DPMA line entries for the desk phones that serverID
// hosts. Each entry is the one Lines generates for the same line, in the
// same order.
func ServerLines(data *astorg.DataSet, p Placement, serverID string) []dpma.Line {
	hosted := hostedSet(p, serverID)
	var out []dpma.Line
	for _, line := range Lines(data) {
		if hosted[line.Name] {
			out = append(out, line)
		}
	}
	return out
}

// ServerPhones generates DPMA phone entries for the desk phones that
// serverID hosts. Each entry is the one Phones generates for the same
// phone, in the same order.
func ServerPhones(data *astorg.DataSet, p Placement, serverID string, base dpma.Phone, contactsURL string) []dpma.Phone {
	var out []dpma.Phone
	for _, phone := range Phones(data, base, contactsURL) {
		if host, ok := p.Phones[phone.MAC]; ok && host == serverID {
			out = append(out, phone)
		}
	}
	return out
}

// ServerMailboxes generates the voicemail section for serverID. Voicemail
// lives on the hub alone, so every other server gets an empty section with
// the same context.
func ServerMailboxes(data *astorg.DataSet, p Placement, serverID string, context string) voicemail.Section {
	if serverID != p.Hub {
		return voicemail.Section{Context: context}
	}
	return Mailboxes(data, context)
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./astgen -run 'TestServerDeviceSections|TestServerMailboxes' -v && gofmt -l . && go vet ./astgen`
Expected: PASS. The hub row has no lines or phones: it hosts only softphones.

- [ ] **Step 5: Commit**

```bash
git add astgen/server.go astgen/server_test.go
git commit -F - <<'EOF'
feat(astgen): render device sections per server

ServerEndpoints, ServerAORs, ServerLines and ServerPhones keep the
sections for the devices a server hosts, and ServerMailboxes renders
voicemail on the hub alone.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 9: Trunks and state publication

**Files:**
- Create: `astgen/trunk.go`
- Test: `astgen/trunk_test.go`

**Interfaces:**
- Consumes: `Placement.Peers`, `Placement.HostedDevices`, `Placement.Hub` (Task 7); names from `astgen/names.go` (Task 7); `pjsip.Identify`, `pjsip.OutboundPublish`, `pjsip.InboundPublication`, `pjsip.PublicationEvents`, `pjsip.PublicationEvent`, `pjsip.AsteriskPublication` (Tasks 1 to 4).
- Produces:
  - `type TrunkConfig struct { Transport pjsip.Transport; Endpoints []pjsip.Endpoint; AORs []pjsip.AOR; Identifies []pjsip.Identify; OutboundPublishes []pjsip.OutboundPublish; InboundPublications []pjsip.InboundPublication; AsteriskPublications []pjsip.AsteriskPublication }`
  - `func ServerTrunks(p Placement, serverID, bindAddr string) TrunkConfig`
  - unexported `deviceStateFilter(devices []string) string`

Per peer, in peer ID order: endpoint, aor, identify, outbound-publish (`publish-<peer>`; the hub adds `publish-mwi-<peer>`), inbound-publication, asterisk-publication. See Spec decisions 1 to 3 for the filter placement, `mailbox_state` and the URI user part. No section carries a password; the auth sections are written by the node agent from `Placement.Credentials`.

- [ ] **Step 1: Write the failing tests**

Create `astgen/trunk_test.go`:

```go
package astgen_test

import (
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
)

// renderTrunks renders every section of a trunk configuration in the order
// a pjsip.conf would hold them.
func renderTrunks(t *testing.T, cfg astgen.TrunkConfig) string {
	t.Helper()
	values := []interface{}{&cfg.Transport}
	for i := range cfg.Endpoints {
		values = append(values, &cfg.Endpoints[i])
	}
	for i := range cfg.AORs {
		values = append(values, &cfg.AORs[i])
	}
	for i := range cfg.Identifies {
		values = append(values, &cfg.Identifies[i])
	}
	for i := range cfg.OutboundPublishes {
		values = append(values, &cfg.OutboundPublishes[i])
	}
	for i := range cfg.InboundPublications {
		values = append(values, &cfg.InboundPublications[i])
	}
	for i := range cfg.AsteriskPublications {
		values = append(values, &cfg.AsteriskPublications[i])
	}
	return render(t, values...)
}

// TestServerTrunksOffice renders the trunks of an office server. Its peers
// come in ID order, it publishes only its own devices, it accepts only each
// peer's devices, and it accepts mailbox state from the hub alone.
func TestServerTrunksOffice(t *testing.T) {
	p, _ := place(t, multiServerData())
	got := renderTrunks(t, astgen.ServerTrunks(p, "qry", "100.64.0.3"))
	want := "[transport-tcp]\n" +
		"type = transport\n" +
		"protocol = tcp\n" +
		"bind = 100.64.0.3:5060\n" +
		"\n" +
		"[trunk-den]\n" +
		"type = endpoint\n" +
		"transport = transport-tcp\n" +
		"context = from-trunk\n" +
		"auth = auth-trunk-den\n" +
		"outbound_auth = auth-trunk-den\n" +
		"aors = trunk-den\n" +
		"direct_media = no\n" +
		"\n" +
		"[trunk-hub]\n" +
		"type = endpoint\n" +
		"transport = transport-tcp\n" +
		"context = from-trunk\n" +
		"auth = auth-trunk-hub\n" +
		"outbound_auth = auth-trunk-hub\n" +
		"aors = trunk-hub\n" +
		"direct_media = no\n" +
		"\n" +
		"[trunk-den]\n" +
		"type = aor\n" +
		"contact = sip:den-phone:5060\\;transport=tcp\n" +
		"\n" +
		"[trunk-hub]\n" +
		"type = aor\n" +
		"contact = sip:hub-phone:5060\\;transport=tcp\n" +
		"\n" +
		"[trunk-den]\n" +
		"type = identify\n" +
		"endpoint = trunk-den\n" +
		"match = den-phone\n" +
		"\n" +
		"[trunk-hub]\n" +
		"type = identify\n" +
		"endpoint = trunk-hub\n" +
		"match = hub-phone\n" +
		"\n" +
		"[publish-den]\n" +
		"type = outbound-publish\n" +
		"server_uri = sip:trunk-qry@den-phone:5060\\;transport=tcp\n" +
		"event = asterisk-devicestate\n" +
		"outbound_auth = auth-trunk-den\n" +
		"@device_state_filter = ^PJSIP/(000fd3000003|fred)$\n" +
		"\n" +
		"[publish-hub]\n" +
		"type = outbound-publish\n" +
		"server_uri = sip:trunk-qry@hub-phone:5060\\;transport=tcp\n" +
		"event = asterisk-devicestate\n" +
		"outbound_auth = auth-trunk-hub\n" +
		"@device_state_filter = ^PJSIP/(000fd3000003|fred)$\n" +
		"\n" +
		"[trunk-den]\n" +
		"type = inbound-publication\n" +
		"endpoint = trunk-den\n" +
		"event_asterisk-devicestate = den\n" +
		"\n" +
		"[trunk-hub]\n" +
		"type = inbound-publication\n" +
		"endpoint = trunk-hub\n" +
		"event_asterisk-devicestate = hub\n" +
		"event_asterisk-mwi = hub\n" +
		"\n" +
		"[den]\n" +
		"type = asterisk-publication\n" +
		"devicestate_publish = publish-den\n" +
		"device_state = yes\n" +
		"device_state_filter = ^PJSIP/(barney|lobby)$\n" +
		"\n" +
		"[hub]\n" +
		"type = asterisk-publication\n" +
		"devicestate_publish = publish-hub\n" +
		"device_state = yes\n" +
		"device_state_filter = ^PJSIP/(fred\\.soft|wilma\\.soft|yard\\.soft)$\n" +
		"mailbox_state = yes\n"
	if got != want {
		t.Errorf("ServerTrunks(qry) rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestServerTrunksHubPublishesMailboxes checks that the hub sends mailbox
// state to each peer through a second outbound-publish.
func TestServerTrunksHubPublishesMailboxes(t *testing.T) {
	p, _ := place(t, multiServerData())
	cfg := astgen.ServerTrunks(p, "hub", "100.64.0.1")
	var names []string
	for _, publish := range cfg.OutboundPublishes {
		names = append(names, publish.Name+" "+publish.Event)
	}
	want := []string{
		"publish-den asterisk-devicestate",
		"publish-mwi-den asterisk-mwi",
		"publish-qry asterisk-devicestate",
		"publish-mwi-qry asterisk-mwi",
	}
	if len(names) != len(want) {
		t.Fatalf("OutboundPublishes = %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("OutboundPublishes[%d] = %q, want %q", i, names[i], want[i])
		}
	}
	for _, publication := range cfg.AsteriskPublications {
		if publication.MailboxStatePublish != "publish-mwi-"+publication.Name {
			t.Errorf("asterisk-publication %s mailboxstate_publish = %q", publication.Name, publication.MailboxStatePublish)
		}
		if publication.MailboxState.Specified() {
			t.Errorf("asterisk-publication %s accepts mailbox state on the hub", publication.Name)
		}
	}
}

// TestServerTrunksServerWithoutDevices checks a server with no locations.
// It still gets a trunk to every peer, and its device state filters match
// nothing, both on what it sends and on what its peers accept from it.
func TestServerTrunksServerWithoutDevices(t *testing.T) {
	data := multiServerData()
	data.Servers = append(data.Servers, astorg.Server{ID: "lcy", Name: "Lacey", Address: "lcy-phone"})
	p, _ := place(t, data)

	cfg := astgen.ServerTrunks(p, "lcy", "100.64.0.4")
	if got := len(cfg.Endpoints); got != 3 {
		t.Errorf("lcy has %d trunk endpoints, want 3", got)
	}
	for _, publish := range cfg.OutboundPublishes {
		if publish.DeviceStateFilter != "^$" {
			t.Errorf("%s @device_state_filter = %q, want ^$", publish.Name, publish.DeviceStateFilter)
		}
	}
	found := false
	for _, publication := range astgen.ServerTrunks(p, "qry", "100.64.0.3").AsteriskPublications {
		if publication.Name != "lcy" {
			continue
		}
		found = true
		if publication.DeviceStateFilter != "^$" {
			t.Errorf("qry accepts %q from lcy, want ^$", publication.DeviceStateFilter)
		}
	}
	if !found {
		t.Error("qry has no asterisk-publication for lcy")
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./astgen -run TestServerTrunks -v`
Expected: FAIL to compile with `undefined: astgen.ServerTrunks`.

- [ ] **Step 3: Write the implementation**

Create `astgen/trunk.go`:

```go
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
// bindAddr is the address the trunk transport listens on, which is the
// server's tailnet IP. It is a render input rather than a data set field.
func ServerTrunks(p Placement, serverID, bindAddr string) TrunkConfig {
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
		peerURI := "sip:" + peer.Address + ":" + trunkPort + ";transport=tcp"

		cfg.Endpoints = append(cfg.Endpoints, pjsip.Endpoint{
			Name:         name,
			Transport:    TrunkTransport,
			Context:      TrunkContext,
			Auth:         auth,
			OutboundAuth: auth,
			AORs:         []string{name},
			DirectMedia:  astval.No,
		})
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
		publishURI := "sip:" + trunkName(serverID) + "@" + peer.Address + ":" + trunkPort + ";transport=tcp"
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
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./astgen -run TestServerTrunks -v && gofmt -l . && go vet ./astgen`
Expected: three tests PASS, including `TestServerTrunksServerWithoutDevices` (Review Focus 3). The `\\.` in the hub filter is `regexp.QuoteMeta` escaping the dot in `fred.soft`.

- [ ] **Step 5: Commit**

```bash
git add astgen/trunk.go astgen/trunk_test.go
git commit -F - <<'EOF'
feat(astgen): render trunks and state publication per server

Each server gets a TCP trunk to every peer, recognized by address, and
shares device state over it with res_pjsip_publish_asterisk. The hub also
publishes mailbox state. Filters keep each server publishing only its own
devices and accepting only each peer's.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 10: Per-server paging

**Files:**
- Create: `astgen/serverpaging.go`
- Test: `astgen/serverpaging_test.go`

**Interfaces:**
- Consumes: `Placement.Hosts` (Task 7), `trunkName`, `sortedUnique` (Task 7), `lineUsername` (existing).
- Consumes also: `PageAlertContext`, `PageRelayContext` (Task 7); `dialplan.PredialTarget`, `PageApp.Predial`, `dialplan.Channel`, `dialplan.PJSIPHeader`, `dialplan.Return` (Task 6).
- Produces:
  - `func ServerPagingExtensions(data *astorg.DataSet, p Placement, serverID string, context string) dialplan.Section`
  - `func ServerPageRelayExtensions(data *astorg.DataSet, p Placement, serverID string) dialplan.Section` (context `astgen-page-relay`)
  - `func PageAlertSection() dialplan.Section` (context `astgen-page-alert`)
  - unexported `pageRelayName(group, peer string) string`
  - unexported, used by Task 11: `pagingMembers(data *astorg.DataSet) map[string][]string`, `splitMembers(p Placement, serverID string, names []string) (local []dialplan.Device, peers []string)`, `pageActions(group astorg.PagingGroup, devices []dialplan.Device) []dialplan.Action`

`pagingMembers` repeats the three membership steps of `PagingGroupExtensions` with device names instead of `SIP/` devices. It is a copy rather than a refactor so the old generator stays untouched; Task 13 proves the two agree. `pageActions` differs from the old actions in one way: the Alert-Info header moves from `SIPAddHeader` to the predial subroutine (Spec decision 12). Relay legs go through `astgen-page-relay` (Spec decision 13).

- [ ] **Step 1: Write the failing tests**

Create `astgen/serverpaging_test.go`:

```go
package astgen_test

import (
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
)

// TestServerPagingAllMembersRemote checks a paging group with no members on
// the hub. A page from the hub reaches each office through one relay leg,
// and each leg is a local channel, not the trunk.
func TestServerPagingAllMembersRemote(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerPagingExtensions(data, p, "hub", "paging")
	got := render(t, &section)
	want := "[paging]\n" +
		"\n" +
		"; All Call\n" +
		"exten => 500,1,Noop(Paging All Call)\n" +
		"same => n,Playback(beep)\n" +
		"same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})\n" +
		"same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})\n" +
		"same => n,Page(Local/page-relay-500-den@astgen-page-relay&Local/page-relay-500-qry@astgen-page-relay,isb(astgen-page-alert^s^1(intercom)))\n" +
		"same => n,Hangup()\n"
	if got != want {
		t.Errorf("ServerPagingExtensions(hub) rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestServerPagingOffice checks that an office pages its local members
// directly and relays to the one other office that hosts members.
func TestServerPagingOffice(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerPagingExtensions(data, p, "den", "paging")
	got := render(t, &section)
	want := "same => n,Page(PJSIP/barney&PJSIP/lobby&Local/page-relay-500-qry@astgen-page-relay,isb(astgen-page-alert^s^1(intercom)))\n"
	if !strings.Contains(got, want) {
		t.Errorf("ServerPagingExtensions(den) rendered:\n%s\nwant it to contain:\n%s", got, want)
	}
	if strings.Contains(got, "SIPAddHeader") {
		t.Errorf("ServerPagingExtensions(den) uses SIPAddHeader, which does nothing on PJSIP:\n%s", got)
	}
}

// TestServerPagingNoAlert checks that a paging group with no alert pages
// without the predial subroutine.
func TestServerPagingNoAlert(t *testing.T) {
	data := multiServerData()
	data.PagingGroups[0].Alert = ""
	p, _ := place(t, data)
	section := astgen.ServerPagingExtensions(data, p, "den", "paging")
	want := "same => n,Page(PJSIP/barney&PJSIP/lobby&Local/page-relay-500-qry@astgen-page-relay,is)\n"
	if got := render(t, &section); !strings.Contains(got, want) {
		t.Errorf("ServerPagingExtensions(den) rendered:\n%s\nwant it to contain:\n%s", got, want)
	}
}

// TestServerPageRelayExtensions checks the relay context. Each relay
// extension dials one server's from-trunk page- entry over the trunk; a
// server with no relays renders an empty context.
func TestServerPageRelayExtensions(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)

	hub := astgen.ServerPageRelayExtensions(data, p, "hub")
	got := render(t, &hub)
	want := "[astgen-page-relay]\n" +
		"\n" +
		"; Relay All Call to den\n" +
		"exten => page-relay-500-den,1,Dial(PJSIP/page-500@trunk-den)\n" +
		"\n" +
		"; Relay All Call to qry\n" +
		"exten => page-relay-500-qry,1,Dial(PJSIP/page-500@trunk-qry)\n"
	if got != want {
		t.Errorf("ServerPageRelayExtensions(hub) rendered:\n%s\nwant:\n%s", got, want)
	}

	data.Locations[1].PagingGroups = nil // Den's phones leave the group
	p, _ = place(t, data)
	qry := astgen.ServerPageRelayExtensions(data, p, "qry")
	if got, want := render(t, &qry), "[astgen-page-relay]\n"; got != want {
		t.Errorf("ServerPageRelayExtensions(qry) with no remote members rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestPageAlertSection checks the predial subroutine that adds Alert-Info
// on PJSIP callee channels only.
func TestPageAlertSection(t *testing.T) {
	section := astgen.PageAlertSection()
	got := render(t, &section)
	want := "[astgen-page-alert]\n" +
		"\n" +
		"; Add a paging group's Alert-Info header to a paged phone\n" +
		"exten => s,1,ExecIf($[\"${CHANNEL(channeltype)}\"=\"PJSIP\"]?Set(PJSIP_HEADER(add,Alert-Info)=${ARG1}))\n" +
		"same => n,Return()\n"
	if got != want {
		t.Errorf("PageAlertSection() rendered:\n%s\nwant:\n%s", got, want)
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./astgen -run 'TestServerPag|TestPageAlert' -v`
Expected: FAIL to compile with `undefined: astgen.ServerPagingExtensions`.

- [ ] **Step 3: Write the implementation**

Create `astgen/serverpaging.go`:

```go
package astgen

import (
	"fmt"
	"sort"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dialplan"
)

// ServerPagingExtensions generates the paging group extensions for
// serverID. A page reaches the members that serverID hosts directly, and
// reaches the members on each other server through one relay leg per
// server.
//
// Each relay leg is a local channel into PageRelayContext rather than the
// trunk itself. Page's "s" option skips a recipient that is not idle, and a
// trunk endpoint reads as in use while it carries any call, so a direct
// trunk leg could be skipped. A local channel is idle unless the same relay
// is already running, so "s" judges only the local member devices.
func ServerPagingExtensions(data *astorg.DataSet, p Placement, serverID string, context string) dialplan.Section {
	section := dialplan.Section{Context: context}
	members := pagingMembers(data)
	for _, group := range data.PagingGroups {
		local, peers := splitMembers(p, serverID, members[group.Extension])
		devices := local
		for _, peer := range peers {
			devices = append(devices, dialplan.Device{
				Technology: "Local",
				Resource:   pageRelayName(group.Extension, peer) + "@" + PageRelayContext,
			})
		}
		section.Extensions = append(section.Extensions, dialplan.Extension{
			Comment: group.Name,
			Number:  group.Extension,
			Actions: pageActions(group, devices),
		})
	}
	return section
}

// ServerPageRelayExtensions generates PageRelayContext for serverID: one
// extension for each paging group and each other server that hosts members
// of it, which carries the page over the trunk to that server's from-trunk
// page- entry.
func ServerPageRelayExtensions(data *astorg.DataSet, p Placement, serverID string) dialplan.Section {
	section := dialplan.Section{Context: PageRelayContext}
	members := pagingMembers(data)
	for _, group := range data.PagingGroups {
		_, peers := splitMembers(p, serverID, members[group.Extension])
		for _, peer := range peers {
			section.Extensions = append(section.Extensions, dialplan.Extension{
				Comment: fmt.Sprintf("Relay %s to %s", group.Name, peer),
				Number:  pageRelayName(group.Extension, peer),
				Actions: []dialplan.Action{
					dialplan.Dial(dialplan.PJSIP("page-"+group.Extension+"@"+trunkName(peer)), 0),
				},
			})
		}
	}
	return section
}

// PageAlertSection generates the subroutine in PageAlertContext that Page
// runs on each callee channel before dialing it. It adds the paging
// group's Alert-Info header, passed as the first argument, so the phone
// answers the page on its own. A relay leg is a local channel with no SIP
// headers, so the header is added only on a PJSIP channel; the receiving
// server adds it for its own members.
func PageAlertSection() dialplan.Section {
	return dialplan.Section{
		Context: PageAlertContext,
		Extensions: []dialplan.Extension{
			{
				Comment: "Add a paging group's Alert-Info header to a paged phone",
				Number:  "s",
				Actions: []dialplan.Action{
					dialplan.ExecIf(
						dialplan.Equal(dialplan.Channel("channeltype"), dialplan.String("PJSIP")),
						dialplan.Set(dialplan.PJSIPHeader("add", "Alert-Info"), dialplan.Var("ARG1")),
					),
					dialplan.Return(),
				},
			},
		},
	}
}

// pageRelayName returns the name of the relay extension that carries a
// page for the group with the given extension to peer.
func pageRelayName(group, peer string) string {
	return "page-relay-" + group + "-" + peer
}

// splitMembers divides the members of a paging group into the devices that
// serverID hosts, sorted, and the IDs of the other servers that host at
// least one member, sorted. A member that no server hosts is dropped.
func splitMembers(p Placement, serverID string, names []string) (local []dialplan.Device, peers []string) {
	seenPeer := make(map[string]bool)
	for _, name := range sortedUnique(names) {
		host, ok := p.Hosts[name]
		switch {
		case !ok:
			continue
		case host == serverID:
			local = append(local, dialplan.PJSIP(name))
		case !seenPeer[host]:
			seenPeer[host] = true
			peers = append(peers, host)
		}
	}
	sort.Strings(peers)
	return local, peers
}

// pagingMembers returns the device names in each paging group, keyed by
// the group's extension. It follows the same steps as
// PagingGroupExtensions.
func pagingMembers(data *astorg.DataSet) map[string][]string {
	lookup := data.Lookup()
	membership := make(map[string][]string)

	addDevice := func(name string, groups []string) {
		for _, group := range groups {
			membership[group] = append(membership[group], name)
		}
	}

	assigned := make(map[string]bool)
	assignDevice := func(name string, groups []string) {
		if assigned[name] {
			return
		}
		assigned[name] = true
		addDevice(name, groups)
	}

	// Step 1: Add phones in locations with paging groups
	for _, phone := range data.Phones {
		if loc, ok := lookup.LocationByName[phone.Location]; ok {
			addDevice(lineUsername(phone.MAC, lookup), loc.PagingGroups)
		}
	}

	// Step 2: Add phones assigned to people with paging groups
	for _, person := range data.People {
		groups := person.PagingGroups
		for _, t := range person.Tags {
			if tag, ok := lookup.TagsByName[t]; ok {
				groups = append(groups, tag.PagingGroups...)
			}
		}
		for _, mac := range person.Phones {
			assignDevice(lineUsername(mac, lookup), groups)
		}
		for _, softphone := range person.Softphones {
			assignDevice(softphone, groups)
		}
	}

	// Step 3: Add phones assigned to roles with paging groups
	for _, role := range data.PhoneRoles {
		groups := role.PagingGroups
		for _, t := range role.Tags {
			if tag, ok := lookup.TagsByName[t]; ok {
				groups = append(groups, tag.PagingGroups...)
			}
		}
		for _, mac := range role.Phones {
			assignDevice(lineUsername(mac, lookup), groups)
		}
		for _, softphone := range role.Softphones {
			assignDevice(softphone, groups)
		}
	}

	return membership
}

// pageActions returns the actions of a paging group extension that pages
// the given devices. They follow PagingGroupExtensions, except that the
// Alert-Info header is added on each callee channel through
// PageAlertContext, because SIPAddHeader has no effect on a PJSIP channel.
func pageActions(group astorg.PagingGroup, devices []dialplan.Device) []dialplan.Action {
	name := dialplan.CallerID("name")
	number := dialplan.CallerID("num")
	intercom := dialplan.String("Intercom")
	zero := dialplan.Int(0)

	actions := []dialplan.Action{
		dialplan.Noop(fmt.Sprintf("Paging %s", group.Name)),
		dialplan.Playback("beep"),
		dialplan.Set(name, dialplan.IfElse(
			dialplan.GreaterThan(dialplan.Len(number), zero),
			name,
			intercom,
		)),
		dialplan.Set(number, dialplan.IfElse(
			dialplan.GreaterThan(dialplan.Len(number), zero),
			number,
			intercom,
		)),
	}
	if len(devices) > 0 {
		page := dialplan.Page(devices...)
		if group.Alert != "" {
			page.Predial = &dialplan.PredialTarget{
				Context:   PageAlertContext,
				Extension: "s",
				Priority:  1,
				Args:      []string{group.Alert},
			}
		}
		actions = append(actions, page)
	}
	actions = append(actions, dialplan.Hangup())
	return actions
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./astgen -run 'TestServerPag|TestPageAlert' -v && gofmt -l . && go vet ./astgen`
Expected: five tests PASS, including `TestServerPagingAllMembersRemote` (Review Focus 4). In the expected strings, the predial argument separator is `^` and the option closes with three parentheses: one for `1(...)`, one for `b(...)`, one for `Page(...)`.

- [ ] **Step 5: Commit**

```bash
git add astgen/serverpaging.go astgen/serverpaging_test.go
git commit -F - <<'EOF'
feat(astgen): page across servers

A page reaches the members a server hosts directly and the members on
each other server through one relay leg to that server's from-trunk.
Relay legs are local channels, so a busy trunk is never skipped, and the
Alert-Info header is added in a rendered predial subroutine, since
SIPAddHeader has no effect on PJSIP.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 11: Per-server extensions and `from-trunk`

**Files:**
- Create: `astgen/serverdialplan.go`
- Test: `astgen/serverdialplan_test.go`

**Interfaces:**
- Consumes: `Placement` fields and methods (Task 7); `TrunkContext`, `trunkName` (Task 7); `pagingMembers`, `splitMembers`, `pageActions` (Task 10); `Mailboxes` (existing); `dialplan.Busy`, `dialplan.VoiceMail`, `dialplan.Cat`, `dialplan.Substring`, `dialplan.DialExpr`, `dialplan.Subroutine` (Task 6).
- Produces:
  - `func ServerExtensions(data *astorg.DataSet, p Placement, serverID string, context string) dialplan.Section`
  - `func ServerTrunkExtensions(data *astorg.DataSet, p Placement, serverID string, vmContext string) dialplan.Section` (its context is always `from-trunk`)

A home extension: `Noop`, `Gosub(pre-answer,s,1(${EXTEN}))`, `Set(ASTGEN_DIAL=<remote legs>)`, one `ExecIf(DEVICE_STATE=NOT_INUSE?Set(...&PJSIP/<device>))` per local device, `ExecIf(LEN>0?Dial(${ASTGEN_DIAL:1},20))`, then the voicemail step from Spec decision 4. The leading `&` that every append adds is dropped by `${ASTGEN_DIAL:1}`.

- [ ] **Step 1: Write the failing tests**

Create `astgen/serverdialplan_test.go`:

```go
package astgen_test

import (
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
)

// TestServerExtensionsOffice renders the phones' context of an office. Its
// own extension rings its local desk phone when idle and the owner's
// softphone on the hub, every other extension routes to its home, and
// unknown extensions go to the hub.
func TestServerExtensionsOffice(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerExtensions(data, p, "qry", "from-internal")
	got := render(t, &section)
	want := "[from-internal]\n" +
		"\n" +
		"; Fred Flintstone\n" +
		"exten => 100,1,Noop(Call fred)\n" +
		"same => n,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Set(ASTGEN_DIAL=&PJSIP/dev-fred.soft@trunk-hub)\n" +
		"same => n,ExecIf($[${DEVICE_STATE(PJSIP/fred)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/fred))\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-100@trunk-hub))\n" +
		"\n" +
		"; Barney Rubble\n" +
		"exten => 101,1,Dial(PJSIP/101@trunk-den)\n" +
		"\n" +
		"; Wilma Flintstone\n" +
		"exten => 102,1,Dial(PJSIP/102@trunk-hub)\n" +
		"\n" +
		"; Lobby\n" +
		"exten => 200,1,Dial(PJSIP/200@trunk-den)\n" +
		"\n" +
		"; Yard\n" +
		"exten => 201,1,Dial(PJSIP/201@trunk-hub)\n" +
		"\n" +
		"; Unknown extensions go to the hub\n" +
		"exten => _X.,1,Dial(PJSIP/${EXTEN}@trunk-hub)\n"
	if got != want {
		t.Errorf("ServerExtensions(qry) rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestServerExtensionsHub checks the hub's own extensions. Wilma has no
// mailbox, so the no-answer subroutine gets no target. The yard role has no
// mailbox, so it ends in congestion as a role does today. The hub has no
// fallthrough, so an unknown extension is invalid.
func TestServerExtensionsHub(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerExtensions(data, p, "hub", "from-internal")
	got := render(t, &section)
	want := "[from-internal]\n" +
		"\n" +
		"; Fred Flintstone\n" +
		"exten => 100,1,Dial(PJSIP/100@trunk-qry)\n" +
		"\n" +
		"; Barney Rubble\n" +
		"exten => 101,1,Dial(PJSIP/101@trunk-den)\n" +
		"\n" +
		"; Wilma Flintstone\n" +
		"exten => 102,1,Noop(Call wilma)\n" +
		"same => n,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Set(ASTGEN_DIAL=)\n" +
		"same => n,ExecIf($[${DEVICE_STATE(PJSIP/wilma.soft)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/wilma.soft))\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Gosub(no-answer,s,1(${EXTEN}))\n" +
		"\n" +
		"; Lobby\n" +
		"exten => 200,1,Dial(PJSIP/200@trunk-den)\n" +
		"\n" +
		"; Yard\n" +
		"exten => 201,1,Noop(Call yard)\n" +
		"same => n,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Set(ASTGEN_DIAL=)\n" +
		"same => n,ExecIf($[${DEVICE_STATE(PJSIP/yard.soft)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/yard.soft))\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Congestion()\n"
	if got != want {
		t.Errorf("ServerExtensions(hub) rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestHomeExtensionOnlyRemoteDevices checks an extension that lives on an
// office whose owner has no desk phone there, only a softphone on the hub.
// The ring has one remote leg, no local leg, and still dials.
func TestHomeExtensionOnlyRemoteDevices(t *testing.T) {
	data := multiServerData()
	data.Phones = data.Phones[1:] // Drop fred's desk phone
	data.People[0].Phones = nil
	p, _ := place(t, data)
	section := astgen.ServerExtensions(data, p, "qry", "from-internal")
	got := render(t, &section)
	want := "; Fred Flintstone\n" +
		"exten => 100,1,Noop(Call fred)\n" +
		"same => n,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Set(ASTGEN_DIAL=&PJSIP/dev-fred.soft@trunk-hub)\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-100@trunk-hub))\n" +
		"\n"
	if !strings.Contains(got, want) {
		t.Errorf("ServerExtensions(qry) rendered:\n%s\nwant it to contain:\n%s", got, want)
	}
}

// TestServerTrunkExtensionsOffice renders an office's from-trunk context:
// its home extension, a dev- entry per hosted device, and the page- entry
// for its paging members.
func TestServerTrunkExtensionsOffice(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerTrunkExtensions(data, p, "qry", "default")
	got := render(t, &section)
	want := "[from-trunk]\n" +
		"\n" +
		"; Fred Flintstone\n" +
		"exten => 100,1,Noop(Call fred)\n" +
		"same => n,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Set(ASTGEN_DIAL=&PJSIP/dev-fred.soft@trunk-hub)\n" +
		"same => n,ExecIf($[${DEVICE_STATE(PJSIP/fred)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/fred))\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-100@trunk-hub))\n" +
		"\n" +
		"; Ring 000fd3000003 for another server\n" +
		"exten => dev-000fd3000003,1,ExecIf($[${DEVICE_STATE(PJSIP/000fd3000003)}!=NOT_INUSE]?Busy())\n" +
		"same => n,Dial(PJSIP/000fd3000003)\n" +
		"same => n,Hangup()\n" +
		"\n" +
		"; Ring fred for another server\n" +
		"exten => dev-fred,1,ExecIf($[${DEVICE_STATE(PJSIP/fred)}!=NOT_INUSE]?Busy())\n" +
		"same => n,Dial(PJSIP/fred)\n" +
		"same => n,Hangup()\n" +
		"\n" +
		"; All Call\n" +
		"exten => page-500,1,Noop(Paging All Call)\n" +
		"same => n,Playback(beep)\n" +
		"same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})\n" +
		"same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})\n" +
		"same => n,Page(PJSIP/000fd3000003&PJSIP/fred,isb(astgen-page-alert^s^1(intercom)))\n" +
		"same => n,Hangup()\n"
	if got != want {
		t.Errorf("ServerTrunkExtensions(qry) rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestServerTrunkExtensionsHubVoicemail checks the hub's vm- entries, one
// per mailbox, in the order Mailboxes writes them.
func TestServerTrunkExtensionsHubVoicemail(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerTrunkExtensions(data, p, "hub", "default")
	got := render(t, &section)
	want := "; Voicemail for Sales\n" +
		"exten => vm-300,1,VoiceMail(300@default,u)\n" +
		"same => n,Hangup()\n" +
		"\n" +
		"; Voicemail for fred\n" +
		"exten => vm-100,1,VoiceMail(100@default,u)\n" +
		"same => n,Hangup()\n" +
		"\n" +
		"; Voicemail for barney\n" +
		"exten => vm-101,1,VoiceMail(101@default,u)\n" +
		"same => n,Hangup()\n"
	if !strings.HasSuffix(got, want) {
		t.Errorf("ServerTrunkExtensions(hub) rendered:\n%s\nwant it to end with:\n%s", got, want)
	}
	office := astgen.ServerTrunkExtensions(data, p, "qry", "default")
	if text := render(t, &office); strings.Contains(text, "exten => vm-") {
		t.Errorf("ServerTrunkExtensions(qry) has a vm- entry:\n%s", text)
	}
}

// TestServerTrunkExtensionsNoLoops checks every server's from-trunk context.
// It may hold only extensions homed on that server and dev-, vm- and page-
// entries, and nothing in it may dial another trunk except a home
// extension's remote ring legs and voicemail target.
func TestServerTrunkExtensionsNoLoops(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	for _, serverID := range p.ServerIDs() {
		section := astgen.ServerTrunkExtensions(data, p, serverID, "default")
		for _, ext := range section.Extensions {
			switch {
			case strings.HasPrefix(ext.Number, "dev-"),
				strings.HasPrefix(ext.Number, "vm-"),
				strings.HasPrefix(ext.Number, "page-"):
			case p.Homes[ext.Number] == serverID:
			default:
				t.Errorf("%s from-trunk holds %s, which is not homed there", serverID, ext.Number)
			}
		}
		text := render(t, &section)
		for _, line := range strings.Split(text, "\n") {
			if strings.Contains(line, "Dial(PJSIP/") && strings.Contains(line, "@trunk-") {
				t.Errorf("%s from-trunk routes to another server: %s", serverID, line)
			}
		}
	}
}

// TestServerTrunkExtensionsNoPageWithoutMembers checks that a server with
// no members of a paging group has no page- entry for it, so a relayed page
// never arrives where nobody would hear it.
func TestServerTrunkExtensionsNoPageWithoutMembers(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	section := astgen.ServerTrunkExtensions(data, p, "hub", "default")
	if text := render(t, &section); strings.Contains(text, "page-500,") {
		t.Errorf("hub from-trunk has a page- entry with no local members:\n%s", text)
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./astgen -run 'TestServerExtensions|TestHomeExtension|TestServerTrunkExtensions|TestServerPaging' -v`
Expected: FAIL to compile with `undefined: astgen.ServerExtensions`.

- [ ] **Step 3: Write the implementation**

Create `astgen/serverdialplan.go`:

```go
package astgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dialplan"
)

// dialVar is the channel variable that a home extension builds its dial
// string in.
const dialVar = dialplan.Var("ASTGEN_DIAL")

// extensionOwner is a person or role that holds an extension.
type extensionOwner struct {
	Extension string
	Username  string
	Comment   string
	Mailbox   string // Empty when the owner has no mailbox
	Role      bool
}

// extensionOwners returns every person and then every role that holds an
// extension in p, in data set order.
func extensionOwners(data *astorg.DataSet, p Placement) []extensionOwner {
	var owners []extensionOwner
	for _, person := range data.People {
		if p.Owners[person.Extension] != person.Username || person.Username == "" {
			continue
		}
		owner := extensionOwner{
			Extension: person.Extension,
			Username:  person.Username,
			Comment:   person.FullName,
		}
		if person.VoicemailCode != "" {
			owner.Mailbox = person.Extension
		}
		owners = append(owners, owner)
	}
	for _, role := range data.PhoneRoles {
		if p.Owners[role.Extension] != role.Username || role.Username == "" {
			continue
		}
		owners = append(owners, extensionOwner{
			Extension: role.Extension,
			Username:  role.Username,
			Comment:   role.DisplayName,
			Mailbox:   role.MailboxNumber,
			Role:      true,
		})
	}
	return owners
}

// remoteDevice returns the dial device that rings a device hosted on
// another server through the trunk to that server.
func remoteDevice(device, host string) string {
	return "PJSIP/dev-" + device + "@" + trunkName(host)
}

// homeExtension builds the extension that rings every device of an owner
// whose extension lives on this server.
//
// Each local device joins the ring only when it is idle. Each remote device
// always joins, and its own server answers busy for it when it is in use.
//
// The pre-answer and no-answer subroutines are hand-written on each server
// and called with Gosub, because Asterisk 21 removed Macro. Each reads its
// arguments as ${ARG1} and ${ARG2} and must end in Return.
func homeExtension(p Placement, owner extensionOwner) dialplan.Extension {
	home := p.Homes[owner.Extension]

	var remote strings.Builder
	hosts := make([]string, 0, len(p.Remote[owner.Extension]))
	for host := range p.Remote[owner.Extension] {
		hosts = append(hosts, host)
	}
	sort.Strings(hosts)
	for _, host := range hosts {
		for _, device := range p.Remote[owner.Extension][host] {
			remote.WriteString("&" + remoteDevice(device, host))
		}
	}

	actions := []dialplan.Action{
		dialplan.Noop(fmt.Sprintf("Call %s", owner.Username)),
		dialplan.Subroutine("pre-answer", "s", 1, dialplan.Var("EXTEN")),
		dialplan.Set(dialVar, dialplan.String(remote.String())),
	}
	for _, device := range p.Devices[owner.Extension] {
		if p.Hosts[device] != home {
			continue
		}
		local := dialplan.PJSIP(device)
		actions = append(actions, dialplan.ExecIf(
			dialplan.Equal(dialplan.DeviceState(local), dialplan.String("NOT_INUSE")),
			dialplan.Set(dialVar, dialplan.Cat(dialVar, dialplan.String("&"+local.String()))),
		))
	}
	actions = append(actions, dialplan.ExecIf(
		dialplan.GreaterThan(dialplan.Len(dialVar), dialplan.Int(0)),
		dialplan.DialExpr(dialplan.Substring(dialVar, 1), 20),
	))

	switch {
	case owner.Mailbox != "":
		actions = append(actions, dialplan.Subroutine("no-answer", "s", 1, dialplan.Var("EXTEN"), dialplan.String(voicemailTarget(p, home, owner.Mailbox))))
	case owner.Role:
		actions = append(actions, dialplan.Congestion())
	default:
		actions = append(actions, dialplan.Subroutine("no-answer", "s", 1, dialplan.Var("EXTEN")))
	}

	return dialplan.Extension{
		Comment: owner.Comment,
		Number:  owner.Extension,
		Actions: actions,
	}
}

// voicemailTarget returns the dial target that the no-answer subroutine on
// home sends a caller to for mailbox. Voicemail lives on the hub, so the hub
// reaches its own from-trunk entry through a local channel and every other
// server reaches it over the trunk.
func voicemailTarget(p Placement, home, mailbox string) string {
	if home == p.Hub {
		return "Local/vm-" + mailbox + "@" + TrunkContext
	}
	return "PJSIP/vm-" + mailbox + "@" + trunkName(p.Hub)
}

// ServerExtensions generates the phones' dialplan section for serverID.
//
// It holds a home extension for each extension that lives on serverID and
// a route over the trunk for each extension that lives elsewhere, in data
// set order. An office server also sends every unknown extension to the
// hub. On the hub an unknown extension is invalid.
func ServerExtensions(data *astorg.DataSet, p Placement, serverID string, context string) dialplan.Section {
	section := dialplan.Section{Context: context}

	for _, owner := range extensionOwners(data, p) {
		home := p.Homes[owner.Extension]
		if home == serverID {
			section.Extensions = append(section.Extensions, homeExtension(p, owner))
			continue
		}
		section.Extensions = append(section.Extensions, dialplan.Extension{
			Comment: owner.Comment,
			Number:  owner.Extension,
			Actions: []dialplan.Action{
				dialplan.Dial(dialplan.PJSIP(owner.Extension+"@"+trunkName(home)), 0),
			},
		})
	}

	if serverID != p.Hub {
		section.Extensions = append(section.Extensions, dialplan.Extension{
			Comment: "Unknown extensions go to the hub",
			Number:  "_X.",
			Actions: []dialplan.Action{
				dialplan.DialExpr(dialplan.Cat(dialplan.String("PJSIP/"), dialplan.Var("EXTEN"), dialplan.String("@"+trunkName(p.Hub))), 0),
			},
		})
	}

	return section
}

// ServerTrunkExtensions generates the from-trunk dialplan section for
// serverID, which every call from another server enters.
//
// It holds the home extensions of serverID, a dev- entry for each device
// serverID hosts, a vm- entry for each mailbox on the hub, and a page- entry
// for each paging group with members on serverID. It never holds a route to
// another server, so a call cannot bounce between servers.
//
// vmContext is the voicemail context that the hub's mailboxes live in.
func ServerTrunkExtensions(data *astorg.DataSet, p Placement, serverID string, vmContext string) dialplan.Section {
	section := dialplan.Section{Context: TrunkContext}

	// Step 1: Home extensions
	for _, owner := range extensionOwners(data, p) {
		if p.Homes[owner.Extension] == serverID {
			section.Extensions = append(section.Extensions, homeExtension(p, owner))
		}
	}

	// Step 2: One leg of another server's ring for each hosted device
	for _, name := range p.HostedDevices(serverID) {
		device := dialplan.PJSIP(name)
		section.Extensions = append(section.Extensions, dialplan.Extension{
			Comment: fmt.Sprintf("Ring %s for another server", name),
			Number:  "dev-" + name,
			Actions: []dialplan.Action{
				dialplan.ExecIf(dialplan.NotEqual(dialplan.DeviceState(device), dialplan.String("NOT_INUSE")), dialplan.Busy()),
				dialplan.Dial(device, 0),
				dialplan.Hangup(),
			},
		})
	}

	// Step 3: Voicemail, on the hub alone
	if serverID == p.Hub {
		seen := make(map[string]bool)
		for _, box := range Mailboxes(data, vmContext).Mailboxes {
			if seen[box.Extension] {
				continue
			}
			seen[box.Extension] = true
			section.Extensions = append(section.Extensions, dialplan.Extension{
				Comment: fmt.Sprintf("Voicemail for %s", box.Name),
				Number:  "vm-" + box.Extension,
				Actions: []dialplan.Action{
					dialplan.VoiceMail(box.Extension, vmContext, "u"),
					dialplan.Hangup(),
				},
			})
		}
	}

	// Step 4: Pages relayed from another server, to local members only
	members := pagingMembers(data)
	for _, group := range data.PagingGroups {
		local, _ := splitMembers(p, serverID, members[group.Extension])
		if len(local) == 0 {
			continue
		}
		section.Extensions = append(section.Extensions, dialplan.Extension{
			Comment: group.Name,
			Number:  "page-" + group.Extension,
			Actions: pageActions(group, local),
		})
	}

	return section
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `go test ./astgen -run 'TestServerExtensions|TestHomeExtension|TestServerTrunkExtensions|TestServerPag' -v && gofmt -l . && go vet ./astgen && go test ./astgen`
Expected: every listed test PASSes, including `TestHomeExtensionOnlyRemoteDevices` (Review Focus 1) and `TestServerTrunkExtensionsNoPageWithoutMembers` (Review Focus 4). `TestServerTrunkExtensionsNoLoops` passes because no `from-trunk` line dials `PJSIP/<x>@trunk-`.

- [ ] **Step 5: Commit**

```bash
git add astgen/serverdialplan.go astgen/serverdialplan_test.go
git commit -F - <<'EOF'
feat(astgen): render extensions and from-trunk per server

An extension rings every device its owner has, on whichever server each
one registers to; a busy device drops out of that one ring. Other
servers route the extension to its home, offices send unknown numbers to
the hub, and from-trunk holds nothing that routes onward.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 12: Golden tree, determinism and no passwords

**Files:**
- Create: `astgen/golden_test.go`
- Create (generated, then reviewed): `astgen/testdata/den/{extensions.conf,pjsip.conf,res_digium_phone.conf,voicemail.conf}`, the same four under `astgen/testdata/hub/` and `astgen/testdata/qry/`

**Interfaces:**
- Consumes: every `Server*` function, `PageAlertSection` and `Place` (Tasks 7 to 11), `multiServerData`, `render`.
- Produces: the `-update` flag (`go test ./astgen -run TestGolden -update`) and `renderServer(t, data, p, serverID) map[string]string`.

- [ ] **Step 1: Write the tests**

Create `astgen/golden_test.go`:

```go
package astgen_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/dpma"
	"github.com/scjalliance/astconf/pjsip"
)

// update rewrites the golden files under testdata from the current output.
// Run "go test ./astgen -run TestGolden -update", then review the diff
// before committing it.
var update = flag.Bool("update", false, "rewrite the golden files in testdata")

// bindAddrs holds the tailnet IP each fixture server binds its trunk
// transport to.
var bindAddrs = map[string]string{
	"den": "100.64.0.2",
	"hub": "100.64.0.1",
	"qry": "100.64.0.3",
}

// renderServer renders every per-server generator for serverID and returns
// the output keyed by file name.
func renderServer(t *testing.T, data *astorg.DataSet, p astgen.Placement, serverID string) map[string]string {
	t.Helper()

	endpoints := astgen.ServerEndpoints(data, p, serverID, pjsip.Endpoint{Context: "from-internal"}, "default")
	aors := astgen.ServerAORs(data, p, serverID, pjsip.AOR{MaxContacts: astval.NewInt(5), RemoveExisting: astval.Yes})
	trunks := astgen.ServerTrunks(p, serverID, bindAddrs[serverID])
	var pjsipValues []interface{}
	pjsipValues = append(pjsipValues, &trunks.Transport)
	for i := range endpoints {
		pjsipValues = append(pjsipValues, &endpoints[i])
	}
	for i := range aors {
		pjsipValues = append(pjsipValues, &aors[i])
	}
	for i := range trunks.Endpoints {
		pjsipValues = append(pjsipValues, &trunks.Endpoints[i])
	}
	for i := range trunks.AORs {
		pjsipValues = append(pjsipValues, &trunks.AORs[i])
	}
	for i := range trunks.Identifies {
		pjsipValues = append(pjsipValues, &trunks.Identifies[i])
	}
	for i := range trunks.OutboundPublishes {
		pjsipValues = append(pjsipValues, &trunks.OutboundPublishes[i])
	}
	for i := range trunks.InboundPublications {
		pjsipValues = append(pjsipValues, &trunks.InboundPublications[i])
	}
	for i := range trunks.AsteriskPublications {
		pjsipValues = append(pjsipValues, &trunks.AsteriskPublications[i])
	}

	extensions := astgen.ServerExtensions(data, p, serverID, "from-internal")
	paging := astgen.ServerPagingExtensions(data, p, serverID, "paging")
	relay := astgen.ServerPageRelayExtensions(data, p, serverID)
	alert := astgen.PageAlertSection()
	fromTrunk := astgen.ServerTrunkExtensions(data, p, serverID, "default")

	phones := astgen.ServerPhones(data, p, serverID, dpma.Phone{}, "https://contacts.example.com/")
	lines := astgen.ServerLines(data, p, serverID)
	var dpmaValues []interface{}
	for i := range phones {
		dpmaValues = append(dpmaValues, &phones[i])
	}
	for i := range lines {
		dpmaValues = append(dpmaValues, &lines[i])
	}

	mailboxes := astgen.ServerMailboxes(data, p, serverID, "default")

	return map[string]string{
		"pjsip.conf":            render(t, pjsipValues...),
		"extensions.conf":       render(t, &extensions, &paging, &relay, &alert, &fromTrunk),
		"res_digium_phone.conf": render(t, dpmaValues...),
		"voicemail.conf":        render(t, &mailboxes),
	}
}

// TestGolden renders every server in the multi-server fixture and compares
// each file with its golden copy in testdata/<server-id>/.
func TestGolden(t *testing.T) {
	data := multiServerData()
	p, problems, err := astgen.Place(data)
	if err != nil {
		t.Fatalf("Place() error: %v", err)
	}
	if len(problems) != 0 {
		t.Fatalf("Place() problems: %v", problems)
	}

	for _, serverID := range p.ServerIDs() {
		for name, got := range renderServer(t, data, p, serverID) {
			path := filepath.Join("testdata", serverID, name)
			if *update {
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatalf("MkdirAll() error: %v", err)
				}
				if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
					t.Fatalf("WriteFile() error: %v", err)
				}
				continue
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile() error: %v (run with -update to create it)", err)
			}
			if got != string(want) {
				t.Errorf("%s does not match its golden file.\ngot:\n%s\nwant:\n%s", path, got, want)
			}
		}
	}
}

// TestGoldenDeterministic renders the fixture ten times, with a fresh
// placement each time, and checks that every file is the same byte for byte.
func TestGoldenDeterministic(t *testing.T) {
	first := renderAll(t)
	// Go randomizes map iteration order, so repeated renders catch a map
	// that is ranged over without sorting.
	for run := 1; run < 10; run++ {
		again := renderAll(t)
		if len(again) != len(first) {
			t.Fatalf("render %d produced %d files, want %d", run, len(again), len(first))
		}
		for path, want := range first {
			if got := again[path]; got != want {
				t.Fatalf("render %d of %s differs:\n%s\nfirst render:\n%s", run, path, got, want)
			}
		}
	}
}

// renderAll places the fixture and renders every server, keyed by
// "<server-id>/<file>".
func renderAll(t *testing.T) map[string]string {
	t.Helper()
	data := multiServerData()
	p, _, err := astgen.Place(data)
	if err != nil {
		t.Fatalf("Place() error: %v", err)
	}
	out := make(map[string]string)
	for _, serverID := range p.ServerIDs() {
		for name, text := range renderServer(t, data, p, serverID) {
			out[serverID+"/"+name] = text
		}
	}
	return out
}

// TestGoldenNoPasswords checks that no per-server output carries a secret
// from the fixture.
func TestGoldenNoPasswords(t *testing.T) {
	data := multiServerData()
	var secrets []string
	for _, phone := range data.Phones {
		secrets = append(secrets, phone.Secret)
	}
	for _, softphone := range data.Softphones {
		secrets = append(secrets, softphone.Secret)
	}
	for path, text := range renderAll(t) {
		if strings.Contains(text, "password") {
			t.Errorf("%s contains a password setting", path)
		}
		for _, secret := range secrets {
			if strings.Contains(text, secret) {
				t.Errorf("%s contains the secret %q", path, secret)
			}
		}
	}
}
```

- [ ] **Step 2: Run to verify the golden test fails without files**

Run: `go test ./astgen -run TestGolden$ -v`
Expected: FAIL with `ReadFile() error: open testdata/den/extensions.conf: no such file or directory (run with -update to create it)` (the first file name may differ; map order is random).

- [ ] **Step 3: Generate the golden files**

Run: `go test ./astgen -run TestGolden -update`
Expected: `ok`. Twelve files now exist under `astgen/testdata/`.

- [ ] **Step 4: Review the generated files against this plan**

Each file must match the listing below byte for byte. Check with `diff` against a copy you paste from here, or read them side by side. A difference is a bug in Tasks 7 to 11, not a reason to edit the listing (see "How to check rendered output").

`astgen/testdata/den/extensions.conf`:

```text
[from-internal]

; Fred Flintstone
exten => 100,1,Dial(PJSIP/100@trunk-qry)

; Barney Rubble
exten => 101,1,Noop(Call barney)
same => n,Gosub(pre-answer,s,1(${EXTEN}))
same => n,Set(ASTGEN_DIAL=)
same => n,ExecIf($[${DEVICE_STATE(PJSIP/barney)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/barney))
same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))
same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-101@trunk-hub))

; Wilma Flintstone
exten => 102,1,Dial(PJSIP/102@trunk-hub)

; Lobby
exten => 200,1,Noop(Call lobby)
same => n,Gosub(pre-answer,s,1(${EXTEN}))
same => n,Set(ASTGEN_DIAL=)
same => n,ExecIf($[${DEVICE_STATE(PJSIP/lobby)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/lobby))
same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))
same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-300@trunk-hub))

; Yard
exten => 201,1,Dial(PJSIP/201@trunk-hub)

; Unknown extensions go to the hub
exten => _X.,1,Dial(PJSIP/${EXTEN}@trunk-hub)

[paging]

; All Call
exten => 500,1,Noop(Paging All Call)
same => n,Playback(beep)
same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})
same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})
same => n,Page(PJSIP/barney&PJSIP/lobby&Local/page-relay-500-qry@astgen-page-relay,isb(astgen-page-alert^s^1(intercom)))
same => n,Hangup()

[astgen-page-relay]

; Relay All Call to qry
exten => page-relay-500-qry,1,Dial(PJSIP/page-500@trunk-qry)

[astgen-page-alert]

; Add a paging group's Alert-Info header to a paged phone
exten => s,1,ExecIf($["${CHANNEL(channeltype)}"="PJSIP"]?Set(PJSIP_HEADER(add,Alert-Info)=${ARG1}))
same => n,Return()

[from-trunk]

; Barney Rubble
exten => 101,1,Noop(Call barney)
same => n,Gosub(pre-answer,s,1(${EXTEN}))
same => n,Set(ASTGEN_DIAL=)
same => n,ExecIf($[${DEVICE_STATE(PJSIP/barney)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/barney))
same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))
same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-101@trunk-hub))

; Lobby
exten => 200,1,Noop(Call lobby)
same => n,Gosub(pre-answer,s,1(${EXTEN}))
same => n,Set(ASTGEN_DIAL=)
same => n,ExecIf($[${DEVICE_STATE(PJSIP/lobby)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/lobby))
same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))
same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-300@trunk-hub))

; Ring barney for another server
exten => dev-barney,1,ExecIf($[${DEVICE_STATE(PJSIP/barney)}!=NOT_INUSE]?Busy())
same => n,Dial(PJSIP/barney)
same => n,Hangup()

; Ring lobby for another server
exten => dev-lobby,1,ExecIf($[${DEVICE_STATE(PJSIP/lobby)}!=NOT_INUSE]?Busy())
same => n,Dial(PJSIP/lobby)
same => n,Hangup()

; All Call
exten => page-500,1,Noop(Paging All Call)
same => n,Playback(beep)
same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})
same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})
same => n,Page(PJSIP/barney&PJSIP/lobby,isb(astgen-page-alert^s^1(intercom)))
same => n,Hangup()
```

`astgen/testdata/den/pjsip.conf`:

```text
[transport-tcp]
type = transport
protocol = tcp
bind = 100.64.0.2:5060

[lobby]
type = endpoint
context = from-internal
auth = lobby-auth
aors = lobby-aor
callerid = "Lobby" <200>
mailboxes = 300@default
set_var = USER_LOCATION=Den
set_var = OUTBOUND_CALLERID=5555550200
set_var = AREACODE=555

[barney]
type = endpoint
context = from-internal
auth = barney-auth
aors = barney-aor
callerid = "Barney Rubble" <101>
mailboxes = 101@default
set_var = VMCODE=4321
set_var = USERNAME=barney
set_var = USER_LOCATION=Den
set_var = OUTBOUND_CALLERID=5555550200
set_var = AREACODE=555

[lobby-aor]
type = aor
max_contacts = 5
remove_existing = yes

[barney-aor]
type = aor
max_contacts = 5
remove_existing = yes

[trunk-hub]
type = endpoint
transport = transport-tcp
context = from-trunk
auth = auth-trunk-hub
outbound_auth = auth-trunk-hub
aors = trunk-hub
direct_media = no

[trunk-qry]
type = endpoint
transport = transport-tcp
context = from-trunk
auth = auth-trunk-qry
outbound_auth = auth-trunk-qry
aors = trunk-qry
direct_media = no

[trunk-hub]
type = aor
contact = sip:hub-phone:5060\;transport=tcp

[trunk-qry]
type = aor
contact = sip:qry-phone:5060\;transport=tcp

[trunk-hub]
type = identify
endpoint = trunk-hub
match = hub-phone

[trunk-qry]
type = identify
endpoint = trunk-qry
match = qry-phone

[publish-hub]
type = outbound-publish
server_uri = sip:trunk-den@hub-phone:5060\;transport=tcp
event = asterisk-devicestate
outbound_auth = auth-trunk-hub
@device_state_filter = ^PJSIP/(barney|lobby)$

[publish-qry]
type = outbound-publish
server_uri = sip:trunk-den@qry-phone:5060\;transport=tcp
event = asterisk-devicestate
outbound_auth = auth-trunk-qry
@device_state_filter = ^PJSIP/(barney|lobby)$

[trunk-hub]
type = inbound-publication
endpoint = trunk-hub
event_asterisk-devicestate = hub
event_asterisk-mwi = hub

[trunk-qry]
type = inbound-publication
endpoint = trunk-qry
event_asterisk-devicestate = qry

[hub]
type = asterisk-publication
devicestate_publish = publish-hub
device_state = yes
device_state_filter = ^PJSIP/(fred\.soft|wilma\.soft|yard\.soft)$
mailbox_state = yes

[qry]
type = asterisk-publication
devicestate_publish = publish-qry
device_state = yes
device_state_filter = ^PJSIP/(000fd3000003|fred)$
```

`astgen/testdata/den/res_digium_phone.conf`:

```text
[000fd3000002]
type = phone
mac = 000fd3000002
line = lobby
full_name = Lobby
contact = https://contacts.example.com/global.xml
contact = https://contacts.example.com/lobby.xml
blf_items = https://contacts.example.com/lobby.blf.xml
alert = intercom

[000fd3000004]
type = phone
mac = 000fd3000004
line = barney
full_name = Barney Rubble
contact = https://contacts.example.com/global.xml
contact = https://contacts.example.com/barney.xml
blf_items = https://contacts.example.com/barney.blf.xml
alert = intercom

[lobby]
type = line
line_label = 200 Lobby

[barney]
type = line
line_label = 101 Barney Rubble
```

`astgen/testdata/den/voicemail.conf`:

```text
[default]
```

`astgen/testdata/hub/extensions.conf`:

```text
[from-internal]

; Fred Flintstone
exten => 100,1,Dial(PJSIP/100@trunk-qry)

; Barney Rubble
exten => 101,1,Dial(PJSIP/101@trunk-den)

; Wilma Flintstone
exten => 102,1,Noop(Call wilma)
same => n,Gosub(pre-answer,s,1(${EXTEN}))
same => n,Set(ASTGEN_DIAL=)
same => n,ExecIf($[${DEVICE_STATE(PJSIP/wilma.soft)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/wilma.soft))
same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))
same => n,Gosub(no-answer,s,1(${EXTEN}))

; Lobby
exten => 200,1,Dial(PJSIP/200@trunk-den)

; Yard
exten => 201,1,Noop(Call yard)
same => n,Gosub(pre-answer,s,1(${EXTEN}))
same => n,Set(ASTGEN_DIAL=)
same => n,ExecIf($[${DEVICE_STATE(PJSIP/yard.soft)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/yard.soft))
same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))
same => n,Congestion()

[paging]

; All Call
exten => 500,1,Noop(Paging All Call)
same => n,Playback(beep)
same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})
same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})
same => n,Page(Local/page-relay-500-den@astgen-page-relay&Local/page-relay-500-qry@astgen-page-relay,isb(astgen-page-alert^s^1(intercom)))
same => n,Hangup()

[astgen-page-relay]

; Relay All Call to den
exten => page-relay-500-den,1,Dial(PJSIP/page-500@trunk-den)

; Relay All Call to qry
exten => page-relay-500-qry,1,Dial(PJSIP/page-500@trunk-qry)

[astgen-page-alert]

; Add a paging group's Alert-Info header to a paged phone
exten => s,1,ExecIf($["${CHANNEL(channeltype)}"="PJSIP"]?Set(PJSIP_HEADER(add,Alert-Info)=${ARG1}))
same => n,Return()

[from-trunk]

; Wilma Flintstone
exten => 102,1,Noop(Call wilma)
same => n,Gosub(pre-answer,s,1(${EXTEN}))
same => n,Set(ASTGEN_DIAL=)
same => n,ExecIf($[${DEVICE_STATE(PJSIP/wilma.soft)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/wilma.soft))
same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))
same => n,Gosub(no-answer,s,1(${EXTEN}))

; Yard
exten => 201,1,Noop(Call yard)
same => n,Gosub(pre-answer,s,1(${EXTEN}))
same => n,Set(ASTGEN_DIAL=)
same => n,ExecIf($[${DEVICE_STATE(PJSIP/yard.soft)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/yard.soft))
same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))
same => n,Congestion()

; Ring fred.soft for another server
exten => dev-fred.soft,1,ExecIf($[${DEVICE_STATE(PJSIP/fred.soft)}!=NOT_INUSE]?Busy())
same => n,Dial(PJSIP/fred.soft)
same => n,Hangup()

; Ring wilma.soft for another server
exten => dev-wilma.soft,1,ExecIf($[${DEVICE_STATE(PJSIP/wilma.soft)}!=NOT_INUSE]?Busy())
same => n,Dial(PJSIP/wilma.soft)
same => n,Hangup()

; Ring yard.soft for another server
exten => dev-yard.soft,1,ExecIf($[${DEVICE_STATE(PJSIP/yard.soft)}!=NOT_INUSE]?Busy())
same => n,Dial(PJSIP/yard.soft)
same => n,Hangup()

; Voicemail for Sales
exten => vm-300,1,VoiceMail(300@default,u)
same => n,Hangup()

; Voicemail for fred
exten => vm-100,1,VoiceMail(100@default,u)
same => n,Hangup()

; Voicemail for barney
exten => vm-101,1,VoiceMail(101@default,u)
same => n,Hangup()
```

`astgen/testdata/hub/pjsip.conf`:

```text
[transport-tcp]
type = transport
protocol = tcp
bind = 100.64.0.1:5060

[fred.soft]
type = endpoint
context = from-internal
auth = fred.soft-auth
aors = fred.soft-aor
callerid = "Fred Flintstone" <100>
set_var = VMCODE=1234
set_var = USERNAME=fred
set_var = USER_LOCATION=Quarry
set_var = OUTBOUND_CALLERID=5555550100
set_var = AREACODE=555

[wilma.soft]
type = endpoint
context = from-internal
auth = wilma.soft-auth
aors = wilma.soft-aor
callerid = "Wilma Flintstone" <102>
set_var = USERNAME=wilma
set_var = USER_LOCATION=Remote
set_var = OUTBOUND_CALLERID=5555550300
set_var = AREACODE=555

[yard.soft]
type = endpoint
context = from-internal
auth = yard.soft-auth
aors = yard.soft-aor
callerid = "Yard" <201>
set_var = USER_LOCATION=Remote
set_var = OUTBOUND_CALLERID=5555550300
set_var = AREACODE=555

[fred.soft-aor]
type = aor
max_contacts = 5
remove_existing = yes

[wilma.soft-aor]
type = aor
max_contacts = 5
remove_existing = yes

[yard.soft-aor]
type = aor
max_contacts = 5
remove_existing = yes

[trunk-den]
type = endpoint
transport = transport-tcp
context = from-trunk
auth = auth-trunk-den
outbound_auth = auth-trunk-den
aors = trunk-den
direct_media = no

[trunk-qry]
type = endpoint
transport = transport-tcp
context = from-trunk
auth = auth-trunk-qry
outbound_auth = auth-trunk-qry
aors = trunk-qry
direct_media = no

[trunk-den]
type = aor
contact = sip:den-phone:5060\;transport=tcp

[trunk-qry]
type = aor
contact = sip:qry-phone:5060\;transport=tcp

[trunk-den]
type = identify
endpoint = trunk-den
match = den-phone

[trunk-qry]
type = identify
endpoint = trunk-qry
match = qry-phone

[publish-den]
type = outbound-publish
server_uri = sip:trunk-hub@den-phone:5060\;transport=tcp
event = asterisk-devicestate
outbound_auth = auth-trunk-den
@device_state_filter = ^PJSIP/(fred\.soft|wilma\.soft|yard\.soft)$

[publish-mwi-den]
type = outbound-publish
server_uri = sip:trunk-hub@den-phone:5060\;transport=tcp
event = asterisk-mwi
outbound_auth = auth-trunk-den

[publish-qry]
type = outbound-publish
server_uri = sip:trunk-hub@qry-phone:5060\;transport=tcp
event = asterisk-devicestate
outbound_auth = auth-trunk-qry
@device_state_filter = ^PJSIP/(fred\.soft|wilma\.soft|yard\.soft)$

[publish-mwi-qry]
type = outbound-publish
server_uri = sip:trunk-hub@qry-phone:5060\;transport=tcp
event = asterisk-mwi
outbound_auth = auth-trunk-qry

[trunk-den]
type = inbound-publication
endpoint = trunk-den
event_asterisk-devicestate = den

[trunk-qry]
type = inbound-publication
endpoint = trunk-qry
event_asterisk-devicestate = qry

[den]
type = asterisk-publication
devicestate_publish = publish-den
mailboxstate_publish = publish-mwi-den
device_state = yes
device_state_filter = ^PJSIP/(barney|lobby)$

[qry]
type = asterisk-publication
devicestate_publish = publish-qry
mailboxstate_publish = publish-mwi-qry
device_state = yes
device_state_filter = ^PJSIP/(000fd3000003|fred)$
```

`astgen/testdata/hub/res_digium_phone.conf` is empty (zero bytes): the hub hosts no desk phones.

`astgen/testdata/hub/voicemail.conf`:

```text
[default]
300 => -3000,Sales,sales@example.com,,attach=yes|saycid=yes|delete=yes
100 => -1234,fred,fred@example.com,,attach=yes|saycid=yes|delete=no
101 => -4321,barney,,,attach=yes|saycid=yes|delete=no
```

`astgen/testdata/qry/extensions.conf`:

```text
[from-internal]

; Fred Flintstone
exten => 100,1,Noop(Call fred)
same => n,Gosub(pre-answer,s,1(${EXTEN}))
same => n,Set(ASTGEN_DIAL=&PJSIP/dev-fred.soft@trunk-hub)
same => n,ExecIf($[${DEVICE_STATE(PJSIP/fred)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/fred))
same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))
same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-100@trunk-hub))

; Barney Rubble
exten => 101,1,Dial(PJSIP/101@trunk-den)

; Wilma Flintstone
exten => 102,1,Dial(PJSIP/102@trunk-hub)

; Lobby
exten => 200,1,Dial(PJSIP/200@trunk-den)

; Yard
exten => 201,1,Dial(PJSIP/201@trunk-hub)

; Unknown extensions go to the hub
exten => _X.,1,Dial(PJSIP/${EXTEN}@trunk-hub)

[paging]

; All Call
exten => 500,1,Noop(Paging All Call)
same => n,Playback(beep)
same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})
same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})
same => n,Page(PJSIP/000fd3000003&PJSIP/fred&Local/page-relay-500-den@astgen-page-relay,isb(astgen-page-alert^s^1(intercom)))
same => n,Hangup()

[astgen-page-relay]

; Relay All Call to den
exten => page-relay-500-den,1,Dial(PJSIP/page-500@trunk-den)

[astgen-page-alert]

; Add a paging group's Alert-Info header to a paged phone
exten => s,1,ExecIf($["${CHANNEL(channeltype)}"="PJSIP"]?Set(PJSIP_HEADER(add,Alert-Info)=${ARG1}))
same => n,Return()

[from-trunk]

; Fred Flintstone
exten => 100,1,Noop(Call fred)
same => n,Gosub(pre-answer,s,1(${EXTEN}))
same => n,Set(ASTGEN_DIAL=&PJSIP/dev-fred.soft@trunk-hub)
same => n,ExecIf($[${DEVICE_STATE(PJSIP/fred)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/fred))
same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))
same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-100@trunk-hub))

; Ring 000fd3000003 for another server
exten => dev-000fd3000003,1,ExecIf($[${DEVICE_STATE(PJSIP/000fd3000003)}!=NOT_INUSE]?Busy())
same => n,Dial(PJSIP/000fd3000003)
same => n,Hangup()

; Ring fred for another server
exten => dev-fred,1,ExecIf($[${DEVICE_STATE(PJSIP/fred)}!=NOT_INUSE]?Busy())
same => n,Dial(PJSIP/fred)
same => n,Hangup()

; All Call
exten => page-500,1,Noop(Paging All Call)
same => n,Playback(beep)
same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})
same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})
same => n,Page(PJSIP/000fd3000003&PJSIP/fred,isb(astgen-page-alert^s^1(intercom)))
same => n,Hangup()
```

`astgen/testdata/qry/pjsip.conf`:

```text
[transport-tcp]
type = transport
protocol = tcp
bind = 100.64.0.3:5060

[fred]
type = endpoint
context = from-internal
auth = fred-auth
aors = fred-aor
callerid = "Fred Flintstone" <100>
mailboxes = 100@default
set_var = VMCODE=1234
set_var = USERNAME=fred
set_var = USER_LOCATION=Quarry
set_var = OUTBOUND_CALLERID=5555550100
set_var = AREACODE=555

[000fd3000003]
type = endpoint
context = from-internal
auth = 000fd3000003-auth
aors = 000fd3000003-aor
callerid = "QRY-000FD3000003" <UNAVAILABLE>
set_var = USER_LOCATION=Quarry
set_var = OUTBOUND_CALLERID=5555550100
set_var = AREACODE=555

[fred-aor]
type = aor
max_contacts = 5
remove_existing = yes

[000fd3000003-aor]
type = aor
max_contacts = 5
remove_existing = yes

[trunk-den]
type = endpoint
transport = transport-tcp
context = from-trunk
auth = auth-trunk-den
outbound_auth = auth-trunk-den
aors = trunk-den
direct_media = no

[trunk-hub]
type = endpoint
transport = transport-tcp
context = from-trunk
auth = auth-trunk-hub
outbound_auth = auth-trunk-hub
aors = trunk-hub
direct_media = no

[trunk-den]
type = aor
contact = sip:den-phone:5060\;transport=tcp

[trunk-hub]
type = aor
contact = sip:hub-phone:5060\;transport=tcp

[trunk-den]
type = identify
endpoint = trunk-den
match = den-phone

[trunk-hub]
type = identify
endpoint = trunk-hub
match = hub-phone

[publish-den]
type = outbound-publish
server_uri = sip:trunk-qry@den-phone:5060\;transport=tcp
event = asterisk-devicestate
outbound_auth = auth-trunk-den
@device_state_filter = ^PJSIP/(000fd3000003|fred)$

[publish-hub]
type = outbound-publish
server_uri = sip:trunk-qry@hub-phone:5060\;transport=tcp
event = asterisk-devicestate
outbound_auth = auth-trunk-hub
@device_state_filter = ^PJSIP/(000fd3000003|fred)$

[trunk-den]
type = inbound-publication
endpoint = trunk-den
event_asterisk-devicestate = den

[trunk-hub]
type = inbound-publication
endpoint = trunk-hub
event_asterisk-devicestate = hub
event_asterisk-mwi = hub

[den]
type = asterisk-publication
devicestate_publish = publish-den
device_state = yes
device_state_filter = ^PJSIP/(barney|lobby)$

[hub]
type = asterisk-publication
devicestate_publish = publish-hub
device_state = yes
device_state_filter = ^PJSIP/(fred\.soft|wilma\.soft|yard\.soft)$
mailbox_state = yes
```

`astgen/testdata/qry/res_digium_phone.conf`:

```text
[000fd3000001]
type = phone
mac = 000fd3000001
line = fred
full_name = Fred Flintstone
contact = https://contacts.example.com/global.xml
contact = https://contacts.example.com/fred.xml
blf_items = https://contacts.example.com/fred.blf.xml
alert = intercom

[000fd3000003]
type = phone
mac = 000fd3000003
line = 000fd3000003
full_name = QRY-000fd3000003
alert = intercom

[fred]
type = line
line_label = 100 Fred Flintstone

[000fd3000003]
type = line
line_label = QRY-000fd3000003
```

`astgen/testdata/qry/voicemail.conf`:

```text
[default]
```

- [ ] **Step 5: Run all three tests**

Run: `go test ./astgen -run 'TestGolden' -v -count=1`
Expected: `TestGolden`, `TestGoldenDeterministic` and `TestGoldenNoPasswords` PASS.

- [ ] **Step 6: Commit**

```bash
git add astgen/golden_test.go astgen/testdata
git commit -F - <<'EOF'
test(astgen): golden tree for a three-server fixture

Renders every server's pjsip, dialplan, DPMA and voicemail output to
astgen/testdata/<server-id>/. Run with -update to rewrite the files. Also
checks that ten renders are identical and that no output holds a secret.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 13: Single-server equivalence

**Files:**
- Create: `astgen/equivalence_test.go`

**Interfaces:**
- Consumes: every `Server*` function and `Place`; the old `Endpoints`, `AORs`, `Lines`, `Phones`, `Mailboxes`, `PagingGroupExtensions`, `PhoneExtensions`; `testData`, `render`.
- Produces: `singleServerData()`, `renderEach`, `devicesIn` (test helpers only).

This is a characterization test, so it passes as soon as it is written. It is the spec's "no regression for the current server" requirement; see Spec decision 5 for what "same" means.

- [ ] **Step 1: Write the test**

Create `astgen/equivalence_test.go`:

```go
package astgen_test

import (
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/dialplan"
	"github.com/scjalliance/astconf/dpma"
	"github.com/scjalliance/astconf/pjsip"
)

// singleServerData returns testData with one server, qry, that is also the
// hub, and a paging group at Quarry. This is the current deployment: every
// device and extension is local.
func singleServerData() *astorg.DataSet {
	data := testData()
	data.Servers = astorg.ServerList{{ID: "qry", Name: "Quarry", Hub: true, Address: "qry-phone"}}
	data.Locations[0].Server = "qry"
	data.Locations[0].PagingGroups = []string{"500"}
	data.PagingGroups = astorg.PagingGroupList{{Name: "All Call", Extension: "500", Alert: "intercom"}}
	return data
}

// TestSingleServerEquivalence renders the single-server fixture through the
// per-server generators and checks the result against the existing
// generators. Endpoints, aors, lines, phones and mailboxes match byte for
// byte. Paging matches once SIP/ is renamed PJSIP/ and the Alert-Info
// header moves to the predial subroutine. Phone extensions keep the same
// numbers, comments and devices, but their actions differ on purpose: each
// device now drops out of the ring alone when it is busy.
func TestSingleServerEquivalence(t *testing.T) {
	data := singleServerData()
	p, problems, err := astgen.Place(data)
	if err != nil {
		t.Fatalf("Place() error: %v", err)
	}
	if len(problems) != 0 {
		t.Fatalf("Place() problems: %v", problems)
	}
	const id = "qry"

	t.Run("endpoints", func(t *testing.T) {
		base := pjsip.Endpoint{Context: "from-internal"}
		got := astgen.ServerEndpoints(data, p, id, base, "default")
		want := astgen.Endpoints(data, base, "default")
		if a, b := renderEach(t, got), renderEach(t, want); a != b {
			t.Errorf("ServerEndpoints() rendered:\n%s\nEndpoints() rendered:\n%s", a, b)
		}
	})

	t.Run("aors", func(t *testing.T) {
		base := pjsip.AOR{MaxContacts: astval.NewInt(5)}
		got := astgen.ServerAORs(data, p, id, base)
		want := astgen.AORs(data, base)
		if a, b := renderEach(t, got), renderEach(t, want); a != b {
			t.Errorf("ServerAORs() rendered:\n%s\nAORs() rendered:\n%s", a, b)
		}
	})

	t.Run("lines", func(t *testing.T) {
		got := astgen.ServerLines(data, p, id)
		want := astgen.Lines(data)
		if a, b := renderEach(t, got), renderEach(t, want); a != b {
			t.Errorf("ServerLines() rendered:\n%s\nLines() rendered:\n%s", a, b)
		}
	})

	t.Run("phones", func(t *testing.T) {
		got := astgen.ServerPhones(data, p, id, dpma.Phone{}, "https://contacts.example.com/")
		want := astgen.Phones(data, dpma.Phone{}, "https://contacts.example.com/")
		if a, b := renderEach(t, got), renderEach(t, want); a != b {
			t.Errorf("ServerPhones() rendered:\n%s\nPhones() rendered:\n%s", a, b)
		}
	})

	t.Run("mailboxes", func(t *testing.T) {
		got := astgen.ServerMailboxes(data, p, id, "default")
		want := astgen.Mailboxes(data, "default")
		if a, b := render(t, &got), render(t, &want); a != b {
			t.Errorf("ServerMailboxes() rendered:\n%s\nMailboxes() rendered:\n%s", a, b)
		}
	})

	t.Run("paging", func(t *testing.T) {
		got := astgen.ServerPagingExtensions(data, p, id, "paging")
		want := astgen.PagingGroupExtensions(data, "paging")
		a := render(t, &got)
		// Apply the two deliberate changes to the old output: devices are
		// PJSIP/, and the Alert-Info header moves from SIPAddHeader, which
		// does nothing on PJSIP, to Page's predial subroutine.
		b := render(t, &want)
		b = strings.ReplaceAll(b, "SIP/", "PJSIP/")
		b = strings.ReplaceAll(b, "same => n,SIPAddHeader(Alert-Info:intercom)\n", "")
		b = strings.ReplaceAll(b, ",is)\n", ",isb(astgen-page-alert^s^1(intercom)))\n")
		if a != b {
			t.Errorf("ServerPagingExtensions() rendered:\n%s\nPagingGroupExtensions() changed:\n%s", a, b)
		}
		if !strings.Contains(b, "Page(PJSIP/") {
			t.Errorf("the old paging output has no Page line to compare:\n%s", b)
		}
	})

	t.Run("extensions", func(t *testing.T) {
		got := astgen.ServerExtensions(data, p, id, "from-internal")
		want := astgen.PhoneExtensions(data, "from-internal")
		if len(got.Extensions) != len(want.Extensions) {
			t.Fatalf("ServerExtensions() returned %d extensions, PhoneExtensions() returned %d", len(got.Extensions), len(want.Extensions))
		}
		for i := range want.Extensions {
			g, w := got.Extensions[i], want.Extensions[i]
			if g.Number != w.Number || g.Comment != w.Comment {
				t.Errorf("extension %d = %s %q, want %s %q", i, g.Number, g.Comment, w.Number, w.Comment)
			}
			gotDevices := devicesIn(t, g, "PJSIP/")
			wantDevices := devicesIn(t, w, "SIP/")
			if strings.Join(gotDevices, " ") != strings.Join(wantDevices, " ") {
				t.Errorf("extension %s rings %v, want %v", g.Number, gotDevices, wantDevices)
			}
		}
		if text := render(t, &got); strings.Contains(text, "@trunk-") {
			t.Errorf("single-server extensions reach a trunk:\n%s", text)
		}
	})
}

// renderEach renders every element of a slice of section values in order.
func renderEach[T any](t *testing.T, items []T) string {
	t.Helper()
	values := make([]interface{}, 0, len(items))
	for i := range items {
		values = append(values, &items[i])
	}
	return render(t, values...)
}

// devicePattern matches a device in a dial string, such as SIP/fred.
var devicePattern = regexp.MustCompile(`(?:PJ)?SIP/[^,&)\s]+`)

// devicesIn returns the sorted, distinct resource names of the devices with
// the given technology prefix that an extension refers to.
func devicesIn(t *testing.T, ext dialplan.Extension, prefix string) []string {
	t.Helper()
	section := dialplan.Section{Context: "x", Extensions: []dialplan.Extension{ext}}
	seen := make(map[string]bool)
	for _, match := range devicePattern.FindAllString(render(t, &section), -1) {
		if strings.HasPrefix(match, prefix) {
			seen[strings.TrimPrefix(match, prefix)] = true
		}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
```

- [ ] **Step 2: Run it**

Run: `go test ./astgen -run TestSingleServerEquivalence -v`
Expected: PASS for all seven subtests.

- [ ] **Step 3: Prove it can fail**

Temporarily edit `ServerEndpoints` in `astgen/server.go` to skip the first endpoint (`out = out[1:]` before `return out`, guarded by `if len(out) > 0`). Run the test again.
Expected: FAIL in the `endpoints` subtest. Revert the edit and confirm `git diff astgen/server.go` is empty.

- [ ] **Step 4: Commit**

```bash
git add astgen/equivalence_test.go
git commit -F - <<'EOF'
test(astgen): per-server output matches today's single server

With one server that is also the hub, the per-server generators render
the same endpoints, aors, lines, phones and mailboxes as the existing
ones, the same paging apart from SIP/ becoming PJSIP/ and Alert-Info moving
to the predial subroutine, and the same extensions with the same devices.

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

### Task 14: astgen doc, final verification, PR 2 boundary

**Files:**
- Modify: `astgen/doc.go` (whole file)

**Interfaces:**
- Consumes: everything above.
- Produces: nothing new.

- [ ] **Step 1: Update the package doc**

Replace the whole of `astgen/doc.go` with:

```go
// Package astgen is a highly experimental package that converts astorg
// datasets into DPMA-based phone configuration files.
//
// The whole-dataset generators, such as Endpoints and PhoneExtensions,
// render one server that hosts every device and names devices SIP/.
//
// Place and the Server generators, such as ServerEndpoints and
// ServerExtensions, render one server's share of a deployment with office
// servers and a hub. Place decides once where each device and extension
// lives, and each Server generator renders only what one server needs,
// naming devices PJSIP/. These generators never render a password.
//
// TODO: Refactor the functions in this package to avoid ingestion of whole
// datasets. Taking a dataset as an argument hides the actual requirements
// of each generator. Fixing this will probably require that we redesign
// the generator functions to cumulatively apply the astorg phone, person
// and role objects in separate functions.
package astgen
```

- [ ] **Step 2: Full verification**

Run: `gofmt -l . && go vet ./... && go test -count=1 ./...`
Expected: gofmt prints nothing; vet is clean; every package `ok`.

Run: `git grep -nP '[\x{2013}\x{2014}]'`
Expected: no output.

Run: `go test ./dialplan -run Example -v`
Expected: PASS. The dialplan package example still renders `SIP/` pages and dials unchanged after the `page.go` edit.

Run: `git diff main --stat -- astgen/alert.go astgen/callerid.go astgen/dedup.go astgen/devices.go astgen/line.go astgen/mailbox.go astgen/paginggroupextension.go astgen/phone.go astgen/phoneextension.go astgen/pjsip.go astgen/ringtone.go astgen/sip.go astgen/username.go`
Expected: no output. The old generators are untouched.

Run: `git grep -n 'Password\|Secret' -- astgen/placement.go astgen/server.go astgen/trunk.go astgen/serverpaging.go astgen/serverdialplan.go astgen/names.go`
Expected: no output.

- [ ] **Step 3: Commit**

```bash
git add astgen/doc.go
git commit -F - <<'EOF'
docs(astgen): describe the per-server generators

Co-Authored-By: Claude Opus 5.5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01SVDvjArCwFD9e1ndJktmVq
EOF
```

- [ ] **Step 4: PR 2 boundary. Stop.**

Report to the controller: the commit SHAs for Tasks 6 to 14 and the output of `go test ./...`. Mention that paging relies on `Page`'s `b()` option and `PJSIP_HEADER`, which should be tried once on a test server before cutover. The controller opens PR 2 against `scjalliance/astconf`. Do not push.
