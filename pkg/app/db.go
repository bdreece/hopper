/*
 * hopper - A gRPC API for collecting IoT device event messages
 * Copyright (C) 2022 Brian Reece

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
