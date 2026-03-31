package xsql

import (
	"context"
	"strings"
)

type (
	ListQuery struct {
		*DB
		table  string
		fields string
		query  string
		args   []any
		order  string
		limit  int
		offset int
	}
)

func (db *DB) NewList(table string) *ListQuery {
	return &ListQuery{
		DB:     db,
		table:  table,
		fields: "*",
	}
}

// Fields 查询字段
func (lq *ListQuery) Fields(fields string) *ListQuery {
	lq.fields = fields
	return lq
}

// Query 查询条件
func (lq *ListQuery) Query(query string, args ...any) *ListQuery {
	lq.query = query
	lq.args = args
	return lq
}

// Order 排序
func (lq *ListQuery) Order(order string) *ListQuery {
	lq.order = order
	return lq
}

// Limit 分页
func (lq *ListQuery) Limit(page, limit int) *ListQuery {
	lq.limit = limit
	lq.offset = (page - 1) * limit
	return lq
}

// Count 统计数据
func (lq *ListQuery) Count(ctx context.Context, count *int64) error {
	var sqlStr strings.Builder
	// 预估算："SELECT COUNT(*) FROM " + table + " WHERE " + query
	sqlStr.Grow(20 + len(lq.table) + len(lq.query))
	sqlStr.WriteString("SELECT COUNT(*) FROM ")
	sqlStr.WriteString(lq.table)
	if lq.query != "" {
		sqlStr.WriteString(" WHERE ")
		sqlStr.WriteString(lq.query)
	}
	return lq.GetContext(ctx, count, sqlStr.String(), lq.args...)
}

// List 查询数据
func (lq *ListQuery) List(ctx context.Context, data any) error {
	var sqlStr strings.Builder
	// 预估算："SELECT " + fields + " FROM " + table + " WHERE " + query + " ORDER BY " + order + " LIMIT ? OFFSET ?"
	whereLen := 0
	if lq.query != "" {
		whereLen = 8 + len(lq.query)
	}
	orderLen := 0
	if lq.order != "" {
		orderLen = 11 + len(lq.order)
	}
	sqlStr.Grow(9 + len(lq.fields) + 7 + len(lq.table) + whereLen + orderLen + 18)
	sqlStr.WriteString("SELECT ")
	sqlStr.WriteString(lq.fields)
	sqlStr.WriteString(" FROM ")
	sqlStr.WriteString(lq.table)
	if lq.query != "" {
		sqlStr.WriteString(" WHERE ")
		sqlStr.WriteString(lq.query)
	}
	if lq.order != "" {
		sqlStr.WriteString(" ORDER BY ")
		sqlStr.WriteString(lq.order)
	}
	sqlStr.WriteString(" LIMIT ? OFFSET ?")
	return lq.SelectContext(ctx, data, sqlStr.String(), append(lq.args, lq.limit, lq.offset)...)
}
