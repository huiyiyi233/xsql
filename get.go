package xsql

import (
	"context"
	"strings"
	"time"
)

// Get 执行查询并返回单个结果
func (db *DB) Get(dest any, query string, args ...any) error {
	start := time.Now()
	err := db.DB.Get(dest, query, args...)
	db.logSQL(query, args, start, err)
	return err
}

// GetContext 执行查询并返回单个结果
func (db *DB) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	start := time.Now()
	err := db.DB.GetContext(ctx, dest, query, args...)
	db.logSQL(query, args, start, err)
	return err
}

// Get 执行查询并返回单个结果
func (t *Tx) Get(dest any, query string, args ...any) error {
	start := time.Now()
	err := t.Tx.Get(dest, query, args...)
	t.logSQL(query, args, start, err)
	return err
}

// GetContext 执行查询并返回单个结果
func (t *Tx) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	start := time.Now()
	err := t.Tx.GetContext(ctx, dest, query, args...)
	t.logSQL(query, args, start, err)
	return err
}

// SelectFields 返回指定字段，未指定则使用默认字段
// 默认字段: "a,b,c"
// SelectFields("a,b,c")           → "a,b,c"
// SelectFields("a,b,c", "a")      → "a"
// SelectFields("a,b,c", "a","b")  → "a,b"
func SelectFields(fields string, s ...string) string {
	switch len(s) {
	case 0:
		return fields
	case 1:
		return s[0]
	default:
		return strings.Join(s, ",")
	}
}
