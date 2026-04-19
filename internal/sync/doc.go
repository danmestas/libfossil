// Package sync implements Fossil's multi-round sync and clone protocol.
//
// The client side drives convergence: [Sync] opens a session against a
// [Transport], exchanges xfer cards in rounds until both sides agree on
// the set of blobs, then returns a [SyncResult]. [Clone] creates a new
// repo file and populates it from the remote.
//
// The server side is stateless per round: [HandleSync] processes one
// request and produces one response. Use [XferHandler] to embed sync
// in an [http.ServeMux] alongside operational routes like /healthz.
//
// Observability is injected via the [Observer] interface — pass nil for
// zero-cost no-ops, or supply an implementation (e.g. OTelObserver from
// leaf/telemetry) for traces, metrics, and structured logs.
package sync
