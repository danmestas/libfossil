// Package repo manages Fossil repository database files.
//
// [Create] initializes a new .fossil file with the standard schema,
// admin user, and config entries. [Open] validates an existing file's
// application_id before returning a [Repo] handle.
//
// Both functions accept an optional [simio.Env] for deterministic
// testing (pass nil for production defaults). Close the Repo when done
// to release the SQLite connection.
package repo
