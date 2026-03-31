package xsql

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	Number interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	}
)

// ReplacePlaceholders 替换占位符
func ReplacePlaceholders(format string, args ...string) string {

	if len(args) == 0 {
		return format
	}

	// 使用 strings.Builder 获得更好的性能
	var builder strings.Builder
	builder.Grow(len(format) + len(args)*10) // 预分配内存

	for i := 0; i < len(format); i++ {
		if format[i] == '$' && i+1 < len(format) {
			// 检查下一个字符是否是数字
			if format[i+1] >= '1' && format[i+1] <= '9' {
				// 获取数字索引
				idx := int(format[i+1] - '0') // 1-based
				if idx <= len(args) {
					builder.WriteString(args[idx-1])
					i++ // 跳过数字字符
				} else {
					// 索引超出范围，保留原样
					builder.WriteByte(format[i])
				}
			} else {
				// 单个 $ 占位符，按顺序替换
				builder.WriteString(args[0])
			}
		} else {
			builder.WriteByte(format[i])
		}
	}

	return builder.String()
}

// BuildUpdateClause 从字段映射构建 UPDATE 语句的 SET 子句
func BuildUpdateClause(data map[string]any) (clause string, args []any) {
	if len(data) == 0 {
		return "", nil
	}

	// 预分配内存，提高性能
	clauseBuilder := strings.Builder{}
	clauseBuilder.Grow(len(data) * 10) // 预估平均每个字段10字符

	args = make([]any, 0, len(data))

	isFirst := true
	for field, value := range data {
		if !isFirst {
			clauseBuilder.WriteString(", ")
		}

		clauseBuilder.WriteString(field)
		clauseBuilder.WriteString(" = ?")

		args = append(args, value)
		isFirst = false
	}

	return clauseBuilder.String(), args
}

// AnyToString 将任意类型转换为字符串
func AnyToString(d any) string {
	switch v := d.(type) {
	case string:
		return "'" + v + "'"
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return "NULL"
	default:
		// 处理其他类型
		return fmt.Sprintf("%v", v)
	}
}

// LastInsertId 最后插入ID
func LastInsertId[T Number](v int64, err error) (T, error) {
	if err != nil {
		return 0, err
	}
	return T(v), nil
}
