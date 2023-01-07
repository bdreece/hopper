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

package config

import (
	"github.com/bdreece/hopper/pkg/proto/grpc"
	"github.com/bdreece/hopper/pkg/services/utils"
	"gorm.io/gorm"
)

type Config struct {
	DB       *gorm.DB
	Logger   utils.Logger
	Port     string
	Hostname string
	Username string
	Password string
	Secret   string

	DeviceService      grpc.DeviceServiceServer
	DeviceModelService grpc.DeviceModelServiceServer
	EventService       grpc.EventServiceServer
	FirmwareService    grpc.FirmwareServiceServer
	PropertyService    grpc.PropertyServiceServer
	TenantService      grpc.TenantServiceServer
	TypeService        grpc.TypeServiceServer
	UnitService        grpc.UnitServiceServer
}
