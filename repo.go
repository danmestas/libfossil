package libfossil

import (
	"github.com/danmestas/libfossil/db"
	"github.com/danmestas/libfossil/internal/repo"
)

// Repo is an opaque handle to a Fossil repository.
type Repo struct {
	inner *repo.Repo
	path  string
}

// Path returns the filesystem path to the repository file.
func (r *Repo) Path() string { return r.path }

// Inner returns the underlying internal repo handle.
// This is exported for use by in-module packages (e.g., cli/) that need
// direct access to the repo DB for raw SQL or internal package calls.
func (r *Repo) Inner() *repo.Repo { return r.inner }

// Close closes the repository and releases resources.
func (r *Repo) Close() error {
	if r.inner == nil {
		return nil
	}
	return r.inner.Close()
}

// DB returns the underlying database handle for raw SQL queries.
// Use this when the high-level Repo methods don't cover your use case.
func (r *Repo) DB() *db.DB { return r.inner.DB() }

// WithTx executes fn within a database transaction.
func (r *Repo) WithTx(fn func(tx *db.Tx) error) error { return r.inner.WithTx(fn) }

// Verify checks repository integrity (blob checksums, delta chains).
func (r *Repo) Verify() error {
	return r.inner.Verify()
}
