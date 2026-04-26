---
title: Testing
weight: 50
---

# Testing

This guide describes libfossil's test strategy for contributors adding tests or
investigating failures. For the module layout referenced throughout, see
[`architecture.md`](./architecture).

## Overview

libfossil's test suite has three tiers:

1. **Unit tests** — standard `go test ./...` coverage for individual packages
   under the root module and `internal/...`.
2. **DST simulation tests** — deterministic simulation scenarios under
   `./dst/...` that drive multiple repos through a virtual clock with fault
   injection and invariant checks.
3. **Dual-driver matrix** — the whole suite runs against two SQLite drivers
   (`modernc` and `ncruces`) to catch driver-specific behavior.

## Running tests

The `Makefile` wraps the common targets:

- `make test` — runs `go test ./... -count=1 -timeout=120s` with the default
  `modernc` driver.
- `make test-drivers` — runs the full suite twice:
  `go test ./... -count=1 -timeout=120s` followed by
  `go test -tags test_ncruces ./... -count=1 -timeout=120s`.
- `make test-otel` — runs the out-of-workspace OTel submodule with
  `cd observer/otel && GOWORK=off go test ./... -count=1`.
- `make test-all` — chains `test-drivers`, `test-otel`, and `go build
  ./cmd/libfossil/`. This is the target CI uses and what you should run before
  opening a PR.

`make setup-hooks` points `core.hooksPath` at `.githooks/`. The installed
pre-commit hook runs both drivers with `-short`, `go vet`, the OTel submodule,
and the CLI build in roughly 45 seconds. Skip with `git commit --no-verify`
only in emergencies.

## Driver matrix

libfossil abstracts SQLite behind `db/` and two driver submodules. Both must
pass identical tests:

- **`modernc`** (default) — pure-Go port of SQLite. Used by every build
  that doesn't set a tag. Registered via `internal/testdriver/modernc.go`
  under `//go:build !test_ncruces && !test_mattn`.
- **`ncruces`** (wasm-capable) — selected with `-tags test_ncruces`, suitable
  for WASM targets. Registered via `internal/testdriver/ncruces.go` under
  `//go:build test_ncruces`.

GitHub Actions runs `test` and `test-ncruces` as parallel jobs (see
`.github/workflows/test.yml`). The ncruces job excludes `cmd/libfossil` because
the shipped binary only needs one driver built in; DST and unit tests are
executed under both. For the canonical invocation, use `make test-drivers`
rather than hand-rolling tag flags.

## DST (Deterministic Simulation Testing)

DST is inspired by FoundationDB and TigerBeetle: run many repos in a single
process under a simulated clock, inject faults, and replay a failing seed
deterministically. The harness lives in `dst/simulator.go` and
`dst/invariants.go`.

### Seeds and reproducibility

`dst.SimConfig.Seed` drives every random source used by the simulator: the
event queue, network fault decisions, and a `SeededBuggify` PRNG that gates
fault-injection sites. Any failure on seed N reproduces exactly on seed N.

`dst/scenario_test.go` exposes three flags driving `TestDST`:

```
go test ./dst -run TestDST -seed=42 -level=hostile -steps=10000
```

- `-seed=<int64>` — seed passed to `SimConfig.Seed` (0 uses test-specific
  defaults).
- `-level={normal,adversarial,hostile}` — picks a `severity` struct setting
  `DropRate` and `Buggify`. `normal` disables faults; `hostile` uses 20% drop
  and BUGGIFY.
- `-steps=<int>` — caps the number of events processed.

Seed-parameterized scenarios (e.g. `TestCloneDSTSeedSweep` in
`dst/clone_test.go`) loop over seeds via `t.Run(fmt.Sprintf("seed_%d", seed),
...)`. To reproduce a single sub-seed: `go test ./dst -run
TestCloneDSTSeedSweep/seed_7 -v`.

### Invariants

`dst/invariants.go` defines per-node and cross-node checks. `Simulator.CheckSafety`
runs the per-node set after every `SafetyCheckInterval` steps; convergence
checks run once the sim has quiesced:

- **Safety (anytime):** `CheckBlobIntegrity` (UUID matches hash of expanded
  content), `CheckDeltaChains` (`delta.srcid` resolves to a blob),
  `CheckNoOrphanPhantoms`, `CheckUVIntegrity`, `CheckTagxrefIntegrity`,
  `CheckTableSyncIntegrity`.
- **Convergence (after fault-free sync):** `CheckConvergence` (leaves match
  master blob-for-blob), `CheckSubsetOf`, `CheckUVConvergence`,
  `CheckTableSyncConvergence`, `CheckTombstoneConvergence`.

Failures return a `*InvariantError{Invariant, NodeID, Detail}` so the failing
invariant, the offending node, and the seed are all in the log line.

## BUGGIFY

BUGGIFY marks code points where faults can be injected — truncated batches,
dropped cards, retries, corrupted nonces, byte-flipped content. In production
the check is a single boolean read; under DST the seeded PRNG decides whether
to fire.

Two styles appear in the tree:

1. **Process-global, used outside the sync path.** `simio.Buggify(probability)`
   consults a process-global state toggled by `simio.EnableBuggify(seed)` /
   `DisableBuggify()` (see `simio/buggify.go`). Example from
   `internal/content/content.go`:

   ```go
   // BUGGIFY: flip a byte in expanded content to exercise UUID-mismatch detection.
   if simio.Buggify(0.01) && len(content) > 0 {
       corrupted := make([]byte, len(content))
       copy(corrupted, content)
       corrupted[0] ^= 0xFF
       return corrupted, nil
   }
   ```

2. **Scoped, passed through session options.** `BuggifyChecker` is an interface
   (`Check(site string, probability float64) bool`) defined in
   `internal/sync/session.go` and re-exported from the root as
   `libfossil.BuggifyChecker`. Sync code threads a checker through so each site
   has a stable name:

   ```go
   // BUGGIFY: 5% chance drop the last gimme card.
   if len(gimmes) > 1 && cs.opts.Buggify != nil &&
       cs.opts.Buggify.Check("clone.buildRequest.dropGimme", 0.05) {
       gimmes = gimmes[:len(gimmes)-1]
   }
   ```

`dst.SeededBuggify` in `dst/simulator.go` implements `BuggifyChecker` with a
seeded `*rand.Rand`. `Simulator.New` enables both styles when `SimConfig.Buggify`
is true, wiring independent PRNG streams into the client, the server-side
`MockFossil`, and the global `simio` state.

## Writing new tests

- Prefer table-driven tests in unit code; include success and failure rows.
- If you touch sync, merge, clone, or checkin, add a DST scenario under
  `dst/` and run it across a seed range (see `TestCloneDSTSeedSweep` for the
  pattern).
- Mark any new fault-sensitive code point with BUGGIFY. Use the scoped
  `BuggifyChecker` style when the site lives inside the sync session; reserve
  `simio.Buggify` for code that runs outside an explicit session scope.
- Assert invariants. If the new behavior warrants one, add it to
  `dst/invariants.go` and wire it into `CheckSafety` or a convergence helper.
- Confirm tests pass under both drivers (`make test-drivers`) before
  requesting review.
