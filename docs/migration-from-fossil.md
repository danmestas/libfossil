# Migration from upstream Fossil

This guide is for users of upstream [Fossil SCM](https://fossil-scm.org) (the C
binary distributed by `fossil-scm.org`) who are evaluating `libfossil` — a
pure-Go implementation of the Fossil repository format and sync protocol.

It is intentionally honest about gaps. If a capability is not listed under
"Command coverage" below, assume it is not implemented yet.

## Repo format compatibility

`libfossil` reads and writes the same on-disk `.fossil` SQLite repository
format as upstream Fossil. A repository created by either tool can be opened
by the other at the file level — there is no import/export step and no data
migration is needed. Both implementations share the same blob table layout,
manifest/control artifact encoding, and delta chain representation.

One caveat: features introduced in very recent upstream Fossil releases
(bleeding-edge artifact card types, new Fossil-specific SQL functions used in
triggers) may not yet be parsed by `libfossil`. Stable Fossil repositories
from the last several years round-trip cleanly.

## Binary name and invocation

The compiled binary is named `libfossil` (not `fossil`) to avoid `$PATH`
collisions when both tools are installed on the same machine:

```
go install github.com/danmestas/libfossil/cmd/libfossil@latest
```

`libfossil` uses a subcommand-grouped syntax: the top-level group is `repo`,
and individual operations live beneath it. Upstream Fossil uses a flat
command style. Examples:

| Upstream Fossil              | libfossil                              |
|------------------------------|----------------------------------------|
| `fossil clone URL repo.fossil` | `libfossil repo clone URL repo.fossil` |
| `fossil ci -m "msg" file.txt`  | `libfossil -R repo.fossil repo ci -m "msg" file.txt` |
| `fossil timeline`              | `libfossil -R repo.fossil repo timeline` |

The `-R <path>` flag selects the repository and is a global flag (before the
`repo` subcommand group).

## Command coverage

The following `repo` subcommands are implemented today. Every entry below
corresponds to a registered Kong command in `cli/repo.go` in the libfossil
source tree.

### Repository creation

- `repo new` — create a new, empty repository file.
- `repo clone` — clone a remote Fossil server over HTTP(S).

### Working tree and checkout

- `repo open` — materialize a checkout in a directory (creates `.fslckout`).
- `repo co` — check out a named version.
- `repo ls` — list files in a version.
- `repo status` — show pending changes in the working directory.
- `repo add` — stage files for addition.
- `repo rm` — stage files for removal.
- `repo rename` — rename a tracked file.
- `repo revert` — undo staging changes.
- `repo ci` — commit staged changes.

### History and inspection

- `repo timeline` — repository history.
- `repo diff` — show changes vs. a version.
- `repo annotate` / `repo blame` — annotate file lines with commit history.
- `repo cat` — output artifact content.
- `repo hash` — hash files (SHA1 or SHA3-256).
- `repo info` — repository statistics.
- `repo resolve` — resolve a symbolic name to a UUID.
- `repo schema` — operate on the synced table schema.
- `repo query` — run ad-hoc SQL against the repository DB.

### Branching and merging

- `repo branch` — list/create branches.
- `repo tag` — tag operations.
- `repo merge` — merge a divergent version into the current checkout.
- `repo conflicts` — list unresolved merge conflicts.
- `repo mark-resolved` — mark a conflict as resolved.
- `repo bisect` — binary-search commits for bug introductions.
- `repo undo` / `repo redo` — undo/redo the last operation.
- `repo stash` — stash working changes.

### Collaboration and admin

- `repo user` — user management.
- `repo invite` — generate an invite token for a user.
- `repo config` — get/set/list repository configuration.
- `repo wiki` — wiki page operations.
- `repo uv` — unversioned file operations.
- `repo extract` — extract files from a version.
- `repo delta` — create/apply delta operations.
- `repo verify` — verify repository integrity.

### Known gaps vs. upstream Fossil CLI

The following upstream commands have no `libfossil` equivalent today:

- `fossil ui` / `fossil server` / `fossil http` — no web UI or HTTP admin
  server is shipped.
- `fossil sync` / `fossil push` / `fossil pull` — there is no CLI verb; sync
  is exposed via the Go API (`Repo.Sync`, `Repo.XferHandler`) but not as a
  subcommand.
- `fossil update` — merging a newer trunk into a checkout has no dedicated
  CLI command; use `repo co` plus `repo merge` as needed.
- `fossil search` — no CLI front-end for full-text search.
- `fossil fts-config` / FTS index management — no CLI command.

## What's not supported yet

Beyond the gaps listed above, several broader features of upstream Fossil are
not yet exposed by `libfossil`:

- **Web UI / HTTP admin** — no equivalent of `fossil ui`. The sync protocol
  handler exists as `Repo.XferHandler()` for programmatic embedding, but the
  HTML UI, wiki viewer, ticket tracker, and forum are not implemented.
- **Interactive conflict resolution** — there is no guided merge tool.
  Conflicts are listed via `repo conflicts` and cleared with
  `repo mark-resolved`, but resolution itself happens in your editor.
- **Full-text search CLI** — `internal/search` contains indexing and query
  primitives, but no CLI command surfaces them. Use the Go API directly.
- **Ticket and forum subsystems** — not implemented.
- **Mirror / export bridges** to Git, svn, and other systems — not implemented.

If any of these are load-bearing for your workflow, stay on upstream Fossil
for now.

## Interop patterns

Where `libfossil` shines is in environments where the upstream C binary is
inconvenient to run:

- **CGo-free CI pipelines** — `libfossil` links statically with a pure-Go
  SQLite driver (`db/driver/modernc`), so clone/sync steps work on stripped
  container images and cross-compiled targets with no system dependencies.
- **Embedding in Go services** — the `Repo` handle can be used directly from
  a Go program to read history, extract artifacts, or accept pushes over a
  custom transport (WebSocket, NATS, libp2p).
- **WASI / browser deployments** — the `db/driver/ncruces` SQLite driver
  compiles to WebAssembly, enabling repo access from edge runtimes.

See `examples/` (planned) for concrete starting points. Until those land,
`cmd/libfossil/main.go` is the shortest working example of wiring the CLI
package into a binary.
