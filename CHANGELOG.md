# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - Unreleased

Initial open-source release of `libfossil`, a pure-Go implementation of the
Fossil SCM that reads and writes the same `.fossil` SQLite repository format.

### Added

- Repository lifecycle: create new repos and clone from existing ones.
- Working-tree operations: checkout and checkin.
- Timeline traversal over commits and events.
- Merge and rebase primitives.
- Diff and annotate (blame) over tracked content.
- Manifest parsing and content-addressed blob storage.
- Sync protocol client/server for pulling and pushing between repos.
- Observer interfaces for sync and checkout, allowing external hooks into
  both network sync events and working-tree state transitions.
- SQLite driver abstraction with support for both `modernc.org/sqlite` (pure
  Go) and `ncruces/go-sqlite3` (cgo-free, wasm-based) backends.
- Deterministic simulation test harness with BUGGIFY-style fault injection
  for exercising concurrency and failure paths.
- OpenTelemetry observer provided as a separate submodule to keep the core
  dependency footprint small.
- `wasip1/wasm` build target for running `libfossil` under WASI runtimes.

### Notes

- When migrating from a private `go-libfossil` checkout, two artifacts were
  renamed:
  - The committed merge-state file `.edgesync-merge` is now `.libfossil-merge`.
  - The internal SQLite `config` table key prefix `edgesync-ci-lock-` is now
    `libfossil-ci-lock-`. This is not user-visible but is part of the
    on-disk repo state.
