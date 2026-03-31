package xsql

import (
	"time"
)

// 日志级别常量
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelErr   = "err"
)

// Config 数据库配置
type Config struct {
	Driver          string        // 数据库驱动：mysql, postgres, sqlite3
	DataSource      string        // 数据源，包含用户名、密码、主机地址、端口、数据库名等信息
	MaxOpenConns    int           // 最大打开连接数
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期
	ConnMaxIdleTime time.Duration // 连接最大空闲时间
	LogLevel        string        // 日志级别
}

// DefaultConfig 返回默认配置 (MySQL)
func DefaultConfig() *Config {
	return &Config{
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: time.Minute * 5,
		LogLevel:        LogLevelInfo,
	}
}

// WithDriver 设置数据库驱动
func (c *Config) WithDriver(v string) *Config {
	c.Driver = v
	return c
}

// WithDataSource 设置数据源
func (c *Config) WithDataSource(v string) *Config {
	c.DataSource = v
	return c
}

// WithMaxOpenConns 设置最大打开连接数
func (c *Config) WithMaxOpenConns(n int) *Config {
	c.MaxOpenConns = n
	return c
}

// WithMaxIdleConns 设置最大空闲连接数
func (c *Config) WithMaxIdleConns(n int) *Config {
	c.MaxIdleConns = n
	return c
}

// WithConnMaxLifetime 设置连接最大生命周期
func (c *Config) WithConnMaxLifetime(d time.Duration) *Config {
	c.ConnMaxLifetime = d
	return c
}

// WithConnMaxIdleTime 设置连接最大空闲时间
func (c *Config) WithConnMaxIdleTime(d time.Duration) *Config {
	c.ConnMaxIdleTime = d
	return c
}

// WithLogLevel 设置日志级别
func (c *Config) WithLogLevel(level string) *Config {
	c.LogLevel = level
	return c
}
