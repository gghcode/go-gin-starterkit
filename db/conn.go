package db

import (
	"github.com/jinzhu/gorm"
)

// Conn is object that has database connection.
type Conn struct {
	db *gorm.DB
}

// NewConn return new instance.
func NewConn(db *gorm.DB) *Conn {
	return &Conn{
		db: db,
	}
}

// GetDB return database connection.
func (conn *Conn) GetDB() *gorm.DB {
	if conn == nil {
		return nil
	}

	return conn.db
}
