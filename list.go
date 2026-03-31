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
	sqlStr := strings.Builder{}
	sqlStr.Grow(23 + len(lq.table) + 9 + len(lq.query))
	sqlStr.WriteString("SELECT COUNT(*) FROM " + lq.table)
	if lq.query != "" {
		sqlStr.WriteString(" WHERE " + lq.query)
	}
	return lq.GetContext(ctx, count, sqlStr.String(), lq.args...)
}

// List 查询数据
func (lq *ListQuery) List(ctx context.Context, data any) error {
	sqlStr := strings.Builder{}
	sqlStr.Grow(9 + len(lq.fields) + 8 + len(lq.table) + 9 + len(lq.query) + 12 + len(lq.order) + 19)
	sqlStr.WriteString("SELECT " + lq.fields + " FROM " + lq.table)
	if lq.query != "" {
		sqlStr.WriteString(" WHERE " + lq.query)
	}
	if lq.order != "" {
		sqlStr.WriteString(" ORDER BY " + lq.order)
	}
	sqlStr.WriteString(" LIMIT ? OFFSET ?")
	return lq.SelectContext(ctx, data, sqlStr.String(), append(lq.args, lq.limit, lq.offset)...)
}
