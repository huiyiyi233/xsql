package xsql

import (
	"context"
	"database/sql"
	"strings"
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

// Update 更新数据
func (db *DB) Update(ctx context.Context, data map[string]any, table, where string, args ...any) (sql.Result, error) {
	sqlStr, argValue := update(data, table, where, args...)
	return db.ExecContext(ctx, sqlStr, argValue...)
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

// Update 更新数据
func (t *Tx) Update(ctx context.Context, data map[string]any, table, where string, args ...any) (sql.Result, error) {
	sqlStr, argValue := update(data, table, where, args...)
	return t.ExecContext(ctx, sqlStr, argValue...)
}

// update 更新数据
func update(data map[string]any, table string, where string, args ...any) (string, []any) {

	// 预分配内存，提高性能
	queryBuilder := strings.Builder{}
	// 估算: "UPDATE " + table + " SET " + 字段部分 + WHERE部分
	estimatedLength := 8 + len(table) + 5 + len(data)*15 + len(where) + 7
	queryBuilder.Grow(estimatedLength)

	// 保存 SET 子句的值
	setValues := make([]any, 0, len(data))

	// 构建 SET 子句
	queryBuilder.WriteString("UPDATE ")
	queryBuilder.WriteString(table)
	queryBuilder.WriteString(" SET ")

	i := 0
	for field, value := range data {
		if i > 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(field)
		queryBuilder.WriteString(" = ?")
		setValues = append(setValues, value)
		i++
	}

	// 添加 WHERE 子句
	if where != "" {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(where)
	}

	// 合并参数：先 SET 子句的参数，后 WHERE 子句的参数
	allArgs := make([]any, 0, len(setValues)+len(args))
	allArgs = append(allArgs, setValues...)
	if len(args) > 0 {
		allArgs = append(allArgs, args...)
	}

	return queryBuilder.String(), allArgs
}
