package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Connection struct {
	db *gorm.DB
}

func NewConnection(hostname, username, password, database string) (conn *Connection, err error) {
	conn = new(Connection)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, database)
	conn.db, err = gorm.Open(mysql.Open(dsn))
	return
}
