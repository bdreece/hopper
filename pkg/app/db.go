package app

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DATABASE = "hopperdb"

func NewDB() (db *gorm.DB, err error) {
	hostname := os.Getenv(HOSTNAME)
	username := os.Getenv(USERNAME)
	password := os.Getenv(PASSWORD)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, DATABASE)
	db, err = gorm.Open(mysql.Open(dsn))
	return
}
