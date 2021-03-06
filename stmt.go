package nap

import (
	"context"
	"database/sql"
)

// Stmt is an aggregate prepared statement.
// It holds a prepared statement for each underlying physical db.
type Stmt struct {
	db    *DB
	stmts []*sql.Stmt
}

// Close closes the statement by concurrently closing all underlying
// statements concurrently, returning the first non nil error.
func (s *Stmt) Close() error {
	return scatter(len(s.stmts), func(i int) error {
		return s.stmts[i].Close()
	})
}

// Exec executes a prepared statement with the given arguments
// and returns a Result summarizing the effect of the statement.
// Exec uses the master as the underlying physical db.
func (s *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	return s.Master().Exec(args...)
}

// ExecContext executes a prepared statement with the given arguments
// and returns a Result summarizing the effect of the statement.
// Exec uses the master as the underlying physical db.
func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return s.Master().ExecContext(ctx, args...)
}

// Query executes a prepared query statement with the given
// arguments and returns the query results as a *sql.Rows.
// Query uses a slave as the underlying physical db.
func (s *Stmt) Query(args ...interface{}) (*sql.Rows, error) {
	return s.Slave().Query(args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
// QueryContext uses a slave as the physical db.
func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	return s.Slave().QueryContext(ctx, args...)
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error
// will be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *sql.Row's Scan scans the first selected row and discards the rest.
// QueryRow uses a slave as the underlying physical db.
func (s *Stmt) QueryRow(args ...interface{}) *sql.Row {
	return s.Slave().QueryRow(args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always return a non-nil value.
// Errors are deferred until Row's Scan method is called.
// QueryRowContext uses a slave as the physical db.
func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	return s.Slave().QueryRowContext(ctx, args...)
}

// Master returns the master stmt physical database
func (s *Stmt) Master() *sql.Stmt {
	return s.stmts[0]
}

// Slave returns one of the stmt physical databases which is a slave
func (s *Stmt) Slave() *sql.Stmt {
	return s.stmts[s.db.slave(len(s.db.pdbs))]
}
