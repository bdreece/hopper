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
	"github.com/bdreece/hopper/pkg/config"
	. "github.com/bdreece/hopper/pkg/proto/grpc"
	"google.golang.org/grpc"
)

func NewServer(cfg *config.Config) *grpc.Server {
	logger := cfg.Logger.WithContext("server")
	server := grpc.NewServer()

	logger.Infoln("Registering services...")
	RegisterDeviceServiceServer(server, cfg.DeviceService)
	RegisterDeviceModelServiceServer(server, cfg.DeviceModelService)
	RegisterEventServiceServer(server, cfg.EventService)
	RegisterFirmwareServiceServer(server, cfg.FirmwareService)
	RegisterPropertyServiceServer(server, cfg.PropertyService)
	RegisterTenantServiceServer(server, cfg.TenantService)
	RegisterTypeServiceServer(server, cfg.TypeService)
	RegisterUnitServiceServer(server, cfg.UnitService)

	logger.Infoln("Services registered")
	return server
}
