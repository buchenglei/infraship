package database

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlConfig struct {
	DSN             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	DevMode         bool
}

type MysqlConnector struct {
	conf MysqlConfig
	once sync.Once

	db    *gorm.DB
	sqlDB *sql.DB
}

func NewMysqlConnector(conf MysqlConfig) (*MysqlConnector, error) {
	return &MysqlConnector{
		conf: conf,
	}, nil
}

func (m *MysqlConnector) Ping(ctx context.Context) error {
	return m.sqlDB.Ping()
}

func (m *MysqlConnector) Connect(ctx context.Context) (*gorm.DB, error) {
	var err error
	m.once.Do(func() {
		var (
			db    *gorm.DB
			sqlDB *sql.DB
		)

		db, err = gorm.Open(mysql.Open(m.conf.DSN), &gorm.Config{})
		if err != nil {
			return
		}

		if m.conf.DevMode {
			db = db.Debug()
		}

		sqlDB, err = db.DB()
		if err != nil {
			return
		}
		// SetMaxIdleConns 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(10)
		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(20)
		// SetConnMaxLifetime 设置了连接可复用的最大时间。
		sqlDB.SetConnMaxLifetime(time.Hour)

		if db == nil {
			panic("gorm init fail")
		}
		m.db = db
		m.sqlDB = sqlDB
	})
	if err != nil {
		return nil, err
	}

	return m.db, nil

}

func (m *MysqlConnector) Close(ctx context.Context) error {
	if _sql, err := m.db.DB(); err != nil && _sql != nil {
		_sql.Close()
	}

	return nil
}
