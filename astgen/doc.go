// Package astgen is a highly experimental package that converts astorg
// datasets into DPMA-based phone configuration files.
//
// The package has two families of generators.
//
// The legacy whole-dataset generators, such as Endpoints and
// PhoneExtensions, render one server that hosts every device and names
// devices SIP/.
//
// Place and the Server generators, such as ServerEndpoints and
// ServerExtensions, render one server's share of a deployment with office
// servers and a hub. Place decides once where each device and extension
// lives, and each Server generator renders only what one server needs,
// naming devices PJSIP/. They only read the data set and the placement, so
// several servers may render at once from one of each.
//
// A spare is a desk phone with no location and no owner. Place lists it in
// Placement.Spares and no server hosts it. Every server renders its DPMA
// phone and line and its pjsip endpoint and aor, named by its MAC address,
// as the whole-dataset generators render an unassigned phone, so DPMA can
// provision it at any office with no extension and no user. It gets no
// extension, hint, route or paging membership, and no device state filter
// names it, so no server publishes state for it. A desk phone with no
// location that a person or role owns is reported as RuleOwnedNoLocation.
//
// These generators never render a secret from the data set: no phone or
// softphone secret, voicemail code or mailbox access code, and no trunk
// password, which the node agent writes from the credentials Placement
// lists. Each hub mailbox carries the fixed sentinel LockedMailboxPassword
// instead of its PIN. A value in a base the caller passes in, such as a PIN
// in the dpma.Phone base given to ServerPhones, passes through as given.
//
// Placement.Credentials is the one list of derived secrets a server's node
// agent needs: a trunk credential per peer, and a phone credential, with a
// MAC address as its ID, for each desk phone line the server hosts and
// each spare. Softphone secrets are not derived and are not listed. A
// phone credential's section is the auth section its line's endpoint
// names, "<line>-auth", which is "<MAC>-auth" for a spare or another
// unassigned phone. A consumer must write the section that
// Credential.Section names and never build the name itself.
//
// # What the consumer must provide
//
// The per-server output relies on pieces that this package does not render.
// The consuming generator and the hand-written common configuration must
// provide:
//
//   - voicemail.conf with passwordlocation=spooldir, and each mailbox's
//     secret.conf written into its spool directory by the hub's node agent
//     before voicemail loads. Hub mailboxes carry the locked sentinel
//     password LockedMailboxPassword, so a mailbox with no secret.conf is
//     locked rather than open.
//   - The Gosub subroutines pre-answer and no-answer, each ending in Return.
//     The no-answer subroutine must handle an empty second argument.
//   - A subscribe_context on each phone endpoint that equals the context
//     ServerExtensions renders into.
//   - Includes in the phones' context, in this order: the contexts that
//     hold local extensions, such as paging, feature codes and outbound
//     routes, first, and HubFallthroughContext from ServerHubFallthrough
//     last, on offices only. Asterisk tries a context's own extensions
//     first and then each include in order, so the fallthrough's _X.
//     catches only what nothing else matched.
//   - The TCP transport bind address and a base endpoint that carries the
//     trunks' codec list, both passed to ServerTrunks.
//
// TODO: Refactor the functions in this package to avoid ingestion of whole
// datasets. Taking a dataset as an argument hides the actual requirements
// of each generator. Fixing this will probably require that we redesign
// the generator functions to cumulatively apply the astorg phone, person
// and role objects in separate functions.
package astgen
