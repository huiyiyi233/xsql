package xsql

import (
	"context"
	"database/sql"
	"time"
)

// Query 执行查询
func (db *DB) Query(query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := db.DB.Query(query, args...)
	db.logSQL(query, args, start, err)
	return rows, err
}

// QueryContext 执行查询
func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := db.DB.QueryContext(ctx, query, args...)
	db.logSQL(query, args, start, err)
	return rows, err
}

// Query 执行查询
func (t *Tx) Query(query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := t.Tx.Query(query, args...)
	t.logSQL(query, args, start, err)
	return rows, err
}

// QueryContext 执行查询
func (t *Tx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := t.Tx.QueryContext(ctx, query, args...)
	t.logSQL(query, args, start, err)
	return rows, err
}
