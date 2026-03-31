package xsql

import (
	"fmt"
	"strings"
	"time"
)

// logSQL 打印 SQL 日志
func (db *DB) logSQL(query string, args []any, start time.Time, err error) {
	logSQL(query, args, start, err, db.config.LogLevel)
}

func (t *Tx) logSQL(query string, args []any, start time.Time, err error) {
	logSQL(query, args, start, err, t.config.LogLevel)
}

func logSQL(query string, args []any, start time.Time, err error, logLevel string) {
	// 先判断日志级别，避免不必要的内存分配
	shouldLog := logLevel == LogLevelDebug ||
		(logLevel == LogLevelInfo && err != nil) ||
		(logLevel == LogLevelErr && err != nil)
	if !shouldLog {
		return
	}

	var s strings.Builder
	// 预估算容量：基础长度 + 查询长度 + 参数估计长度
	s.Grow(16 + len(query) + len(args)*10)

	// 添加时间戳和错误信息
	s.WriteString("[SQL][")
	s.WriteString(formatDuration(time.Since(start)))
	s.WriteString("] ")

	// 处理参数替换
	argIndex := 0
	lastPos := 0
	for i := 0; i < len(query); i++ {
		if query[i] == '?' && argIndex < len(args) {
			s.WriteString(query[lastPos:i])
			s.WriteString(AnyToString(args[argIndex]))
			argIndex++
			lastPos = i + 1
		}
	}
	s.WriteString(query[lastPos:])

	// 如果有错误，追加错误信息
	if err != nil {
		s.WriteString(" [ERROR: ")
		s.WriteString(err.Error())
		s.WriteString("]")
	}
	fmt.Println(s.String())
}

// 优化时间格式（更易读）
func formatDuration(d time.Duration) string {
	switch {
	case d < time.Millisecond:
		return fmt.Sprintf("%dµs", d.Microseconds())
	case d < time.Second:
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	default:
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}
