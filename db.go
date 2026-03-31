package xsql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// DB 数据库连接封装
type DB struct {
	*sqlx.DB
	config *Config
}

type Tx struct {
	*sqlx.Tx
	config *Config
}

// NewDB 创建新的数据库连接
func NewDB(cfg *Config) (*DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	sqlxDB, err := sqlx.Connect(cfg.Driver, cfg.DataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 设置连接池配置
	sqlxDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlxDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlxDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlxDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return &DB{
		DB:     sqlxDB,
		config: cfg,
	}, nil
}

// NewSqlxDB 使用现有的 sqlx.DB 创建
func NewSqlxDB(sqlxDB *sqlx.DB, cfg *Config) *DB {
	return &DB{
		DB:     sqlxDB,
		config: cfg,
	}
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	return db.DB.Close()
}

// Ping 检查数据库连接
func (db *DB) Ping() error {
	return db.DB.Ping()
}

// Transaction 执行事务，失败则回滚
func (db *DB) Transaction(f func(tx *Tx) error) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	err = f(&Tx{Tx: tx, config: db.config})
	if err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}
