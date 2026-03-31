package xsql

import (
	"context"
	"time"
)

// Select 执行查询并返回多个结果
func (db *DB) Select(dest any, query string, args ...any) error {
	start := time.Now()
	err := db.DB.Select(dest, query, args...)
	db.logSQL(query, args, start, err)
	return err
}

// SelectContext 执行查询并返回多个结果
func (db *DB) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	start := time.Now()
	err := db.DB.SelectContext(ctx, dest, query, args...)
	db.logSQL(query, args, start, err)
	return err
}

// Select 执行查询并返回多个结果
func (t *Tx) Select(dest any, query string, args ...any) error {
	start := time.Now()
	err := t.Tx.Select(dest, query, args...)
	t.logSQL(query, args, start, err)
	return err
}

// SelectContext 执行查询并返回多个结果
func (t *Tx) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	start := time.Now()
	err := t.Tx.SelectContext(ctx, dest, query, args...)
	t.logSQL(query, args, start, err)
	return err
}
