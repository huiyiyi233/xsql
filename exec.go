package xsql

import (
	"context"
	"database/sql"
	"time"
)

// Exec 执行 SQL 语句
func (db *DB) Exec(query string, args ...any) (sql.Result, error) {
	start := time.Now()
	result, err := db.DB.Exec(query, args...)
	db.logSQL(query, args, start, err)
	return result, err
}

// ExecContext 执行 SQL 语句
func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	result, err := db.DB.ExecContext(ctx, query, args...)
	db.logSQL(query, args, start, err)
	return result, err
}

// Exec 执行 SQL 语句
func (t *Tx) Exec(query string, args ...any) (sql.Result, error) {
	start := time.Now()
	result, err := t.Tx.Exec(query, args...)
	t.logSQL(query, args, start, err)
	return result, err
}

// ExecContext 执行 SQL 语句
func (t *Tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	result, err := t.Tx.ExecContext(ctx, query, args...)
	t.logSQL(query, args, start, err)
	return result, err
}
