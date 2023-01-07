package app

import (
	"fmt"

	"github.com/bdreece/hopper/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DATABASE = "hopperdb"

func NewDB(config *config.Config) (db *gorm.DB, err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.Hostname, DATABASE)
	db, err = gorm.Open(mysql.Open(dsn))
	return
}
