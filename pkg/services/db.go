package services

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(hostname, username, password, database string) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, database)
	db, err = gorm.Open(mysql.Open(dsn))
	return
}
