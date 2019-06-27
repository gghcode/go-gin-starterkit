package db

import (
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

// Conn is object that has database connection.
type Conn struct {
	db *gorm.DB
}

// NewConn return new instance.
func NewConn(config config.Configuration) (*Conn, error) {
	db, err := gorm.Open(config.Postgres.Driver,
		"host="+config.Postgres.Host+
			" port="+config.Postgres.Port+
			" user="+config.Postgres.User+
			" dbname="+config.Postgres.Name+
			" password="+config.Postgres.Password+
			" sslmode=disable")

	if err != nil {
		return nil, errors.Wrap(err, "db connect failed...")
	}

	return &Conn{
		db: db,
	}, nil
}

// GetDB return database connection.
func (conn *Conn) GetDB() *gorm.DB {
	if conn == nil {
		return nil
	}

	return conn.db
}

// Close close db session.
func (conn *Conn) Close() error {
	return conn.GetDB().Close()
}
