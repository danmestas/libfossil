package db

import (
	"database/sql"
	"fmt"
	"os"
)

// DB wraps a SQLite database connection.
type DB struct {
	conn   *sql.DB
	path   string
	driver string
}

// Open opens a SQLite database with the registered driver and default pragmas.
func Open(path string) (*DB, error) {
	return OpenWith(path, OpenConfig{})
}

// OpenWith opens a SQLite database with explicit configuration.
func OpenWith(path string, cfg OpenConfig) (*DB, error) {
	if path == "" {
		panic("db.OpenWith: path must not be empty")
	}
	if registered == nil {
		panic("db.OpenWith: no driver registered — import a driver package (e.g., _ \"github.com/danmestas/libfossil/db/driver/modernc\")")
	}
	driver := cfg.Driver
	if driver == "" {
		driver = registered.Name
	}

	pragmas := DefaultPragmas()
	for k, v := range cfg.Pragmas {
		pragmas[k] = v
	}

	// WASM targets: skip DSN pragmas (ncruces _pragma syntax fails under WASI
	// file locking). Open with nolock=1, then apply safe pragmas via SQL.
	var dsn string
	if wasmClearPragmas {
		suffix := wasmDSNSuffix()
		if suffix != "" {
			dsn = fmt.Sprintf("file:%s?%s", path, suffix)
		} else {
			dsn = path
		}
	} else {
		dsn = registered.BuildDSN(path, pragmas)
	}
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("db.Open(%s): %w", driver, err)
	}

	if wasmClearPragmas {
		// Apply safe pragmas via SQL. Skip journal_mode (WAL not supported on WASI).
		// Only apply known-safe pragmas to prevent SQL injection.
		safePragmas := map[string]bool{
			"busy_timeout": true,
			"foreign_keys": true,
		}
		for k, v := range pragmas {
			if !safePragmas[k] {
				continue
			}
			if _, err := conn.Exec(fmt.Sprintf("PRAGMA %s = %s", k, v)); err != nil {
				conn.Close()
				return nil, fmt.Errorf("db.Open PRAGMA %s: %w", k, err)
			}
		}
	}

	return &DB{conn: conn, path: path, driver: driver}, nil
}

// SqlDB returns the underlying *sql.DB connection.
func (d *DB) SqlDB() *sql.DB {
	if d == nil {
		panic("db.SqlDB: receiver must not be nil")
	}
	return d.conn
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) Path() string {
	return d.path
}

func (d *DB) Driver() string {
	return d.driver
}

func (d *DB) Exec(query string, args ...any) (sql.Result, error) {
	return d.conn.Exec(query, args...)
}

func (d *DB) QueryRow(query string, args ...any) *sql.Row {
	return d.conn.QueryRow(query, args...)
}

func (d *DB) Query(query string, args ...any) (*sql.Rows, error) {
	return d.conn.Query(query, args...)
}

func (d *DB) SetApplicationID(id int32) error {
	_, err := d.conn.Exec(fmt.Sprintf("PRAGMA application_id=%d", id))
	return err
}

func (d *DB) ApplicationID() (int32, error) {
	var id int32
	err := d.conn.QueryRow("PRAGMA application_id").Scan(&id)
	return id, err
}

type Tx struct {
	tx *sql.Tx
}

func (t *Tx) Exec(query string, args ...any) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

func (t *Tx) QueryRow(query string, args ...any) *sql.Row {
	return t.tx.QueryRow(query, args...)
}

func (t *Tx) Query(query string, args ...any) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}

func (d *DB) WithTx(fn func(tx *Tx) error) error {
	if d == nil {
		panic("db.WithTx: receiver must not be nil")
	}
	if fn == nil {
		panic("db.WithTx: fn must not be nil")
	}
	sqlTx, err := d.conn.Begin()
	if err != nil {
		return fmt.Errorf("db.WithTx begin: %w", err)
	}
	defer func() {
		if rbErr := sqlTx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
			fmt.Fprintf(os.Stderr, "db.WithTx: rollback failed: %v\n", rbErr)
		}
	}()
	if err := fn(&Tx{tx: sqlTx}); err != nil {
		return err
	}
	return sqlTx.Commit()
}
