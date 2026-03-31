# xsql

基于 [sqlx](https://github.com/jmoiron/sqlx) 的数据库操作封装库，提供简洁、高效的数据库操作方法。

## 特性

- ✅ 基于 sqlx，兼容标准 database/sql 接口
- ✅ 支持多种数据库驱动（MySQL、PostgreSQL、SQLite3）
- ✅ 内置 SQL 日志记录，支持多级日志控制
- ✅ 完整的 Context 上下文支持
- ✅ 事务处理封装
- ✅ 链式查询构建器
- ✅ 高性能优化（预分配内存、避免反射等）

## 安装

```bash
go get github.com/huiyiyi233/xsql
```

## 快速开始

### 基本使用

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/huiyiyi233/xsql"
)

type User struct {
    ID       int64  `db:"id"`
    Name     string `db:"name"`
    Email    string `db:"email"`
    Age      int    `db:"age"`
}

func main() {
    // 创建配置
    cfg := xsql.DefaultConfig().
        WithDriver("mysql").
        WithDataSource("user:password@tcp(localhost:3306)/testdb").
        WithMaxOpenConns(100).
        WithMaxIdleConns(10).
        WithConnMaxLifetime(time.Hour).
        WithLogLevel(xsql.LogLevelDebug)
    
    // 创建数据库连接
    db, err := xsql.NewDB(cfg)
    if err != nil {
        panic(err)
    }
    defer db.Close()
    
    // 测试连接
    if err := db.Ping(); err != nil {
        panic(err)
    }
}
```

## API 文档

### 配置

#### Config 结构体

```go
type Config struct {
    Driver          string        // 数据库驱动：mysql, postgres, sqlite3
    DataSource      string        // 数据源连接字符串
    MaxOpenConns    int           // 最大打开连接数（默认：100）
    MaxIdleConns    int           // 最大空闲连接数（默认：10）
    ConnMaxLifetime time.Duration // 连接最大生命周期（默认：1 小时）
    ConnMaxIdleTime time.Duration // 连接最大空闲时间（默认：5 分钟）
    LogLevel        string        // 日志级别
}
```

#### 日志级别常量

```go
const (
    LogLevelDebug = "debug"  // 输出所有 SQL
    LogLevelInfo  = "info"   // 仅输出错误 SQL
    LogLevelErr   = "err"    // 仅在出错时输出
)
```

#### 配置方法

```go
// 使用默认配置
cfg := xsql.DefaultConfig()

// 链式设置
cfg.WithDriver("mysql").
    WithDataSource("dsn...").
    WithMaxOpenConns(50).
    WithMaxIdleConns(5).
    WithLogLevel(xsql.LogLevelDebug)
```

### 数据库操作

#### 创建连接

```go
// 从配置创建新连接
db, err := xsql.NewDB(cfg)

// 使用现有的 sqlx.DB 创建
db := xsql.NewSqlxDB(existingSqlxDB, cfg)
```

#### 基础方法

```go
// 关闭连接
db.Close()

// 检查连接
db.Ping()
```

#### Exec - 执行 SQL 语句

```go
// 基础执行
result, err := db.Exec("INSERT INTO users(name, age) VALUES(?, ?)", "John", 25)

// 带 Context
ctx := context.Background()
result, err := db.ExecContext(ctx, "UPDATE users SET age = ? WHERE id = ?", 26, 1)

// 获取影响行数
rowsAffected, _ := result.RowsAffected()
```

#### Get - 查询单条记录

```go
var user User
err := db.Get(&user, "SELECT * FROM users WHERE id = ?", 1)

// 带 Context
err = db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", 1)
```

#### Select - 查询多条记录

```go
var users []User
err := db.Select(&users, "SELECT * FROM users WHERE age > ?", 18)

// 带 Context
err = db.SelectContext(ctx, &users, "SELECT * FROM users WHERE age > ?", 18)
```

#### Query - 返回 sql.Rows

```go
rows, err := db.Query("SELECT * FROM users WHERE age > ?", 18)
if err != nil {
    // 处理错误
}
defer rows.Close()

for rows.Next() {
    var user User
    if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
        // 处理错误
    }
}
```

### 事务处理

```go
err := db.Transaction(func(tx *xsql.Tx) error {
    // 在事务中执行操作
    _, err := tx.Exec("INSERT INTO users(name) VALUES(?)", "Alice")
    if err != nil {
        return err // 自动回滚
    }
    
    _, err = tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE user_id = ?", 1)
    if err != nil {
        return err // 自动回滚
    }
    
    return nil // 成功则提交
})

if err != nil {
    // 处理错误（已回滚）
}
```

### 链式查询构建器

```go
// 创建查询构建器
lq := db.NewList("users")

// 构建查询
var count int64
var users []User

// 统计数量
err := lq.Query("age > ?", 18).Count(ctx, &count)

// 查询列表（带分页和排序）
err = db.NewList("users").
    Fields("id, name, email").
    Query("age > ? AND status = ?", 18, 1).
    Order("created_at DESC").
    Limit(1, 10).  // 第 1 页，每页 10 条
    List(ctx, &users)
```

#### ListQuery 方法

```go
// 设置查询字段
Fields("id, name, email")

// 设置 WHERE 条件
Query("age > ? AND status = ?", 18, 1)

// 设置排序
Order("created_at DESC")

// 设置分页（页码，每页数量）
Limit(1, 10)

// 统计总数
Count(ctx, &count)

// 查询列表
List(ctx, &data)
```

### 工具函数

#### ReplacePlaceholders - 替换占位符

```go
sql := xsql.ReplacePlaceholders("SELECT * FROM $ WHERE id = $", "users", "1")
// 输出：SELECT * FROM users WHERE id = 1
```

#### BuildUpdateClause - 构建 UPDATE SET 子句

```go
data := map[string]any{
    "name":  "John",
    "age":   25,
    "email": "john@example.com",
}

clause, args := xsql.BuildUpdateClause(data)
// clause: "name = ?, age = ?, email = ?"
// args: ["John", 25, "john@example.com"]

// 使用示例
query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", clause)
args = append(args, 1)
db.Exec(query, args...)
```

#### AnyToString - 类型转字符串

```go
str := xsql.AnyToString(42)           // "42"
str := xsql.AnyToString("hello")      // "'hello'"
str := xsql.AnyToString(true)         // "true"
str := xsql.AnyToString(nil)          // "NULL"
str := xsql.AnyToString(3.14)         // "3.14"
```

## 日志输出

SQL 日志格式：

```
[SQL][1.23ms] SELECT * FROM users WHERE id = 1
[SQL][250µs] INSERT INTO users(name) VALUES('Alice')
[SQL][5.67ms] UPDATE users SET age = 26 [ERROR: duplicate key value]
```

日志级别说明：
- `LogLevelDebug`: 输出所有 SQL 语句
- `LogLevelInfo`: 仅输出错误的 SQL
- `LogLevelErr`: 仅在出错时输出

## 性能优化

本库在以下方面进行了性能优化：

1. **内存预分配**: SQL 构建时使用 `strings.Builder.Grow()` 预分配内存
2. **避免反射**: 类型转换使用类型断言而非反射
3. **连接池优化**: 合理配置连接池参数
4. **零拷贝**: 尽可能减少内存复制

## 支持的数据库

- MySQL (使用 `github.com/go-sql-driver/mysql`)
- PostgreSQL (使用 `github.com/lib/pq`)
- SQLite3 (使用 `github.com/mattn/go-sqlite3`)

## 注意事项

1. 确保导入对应的数据库驱动包并注册
2. 使用完 `Rows` 后记得调用 `Close()`
3. 合理使用 Context 控制查询超时
4. 根据实际场景调整连接池参数
5. 生产环境建议使用 `LogLevelInfo` 或 `LogLevelErr`

## License

MIT License
