package app

import (
	"fmt"

	"github.com/bdreece/hopper/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DATABASE = "hopperdb"

func NewDB(cfg *config.Config) (db *gorm.DB, err error) {
	logger := cfg.Logger.WithContext("db")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		cfg.Username, cfg.Password, cfg.Hostname, DATABASE)

	logger.Infoln("Opening database...")
	db, err = gorm.Open(mysql.Open(dsn))
	return
}
