package db

import "database/sql"

// Querier is the common interface satisfied by both *DB and *Tx.
// Functions that need to work inside transactions accept Querier
// instead of *DB.
type Querier interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}
