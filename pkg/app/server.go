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
	pb "github.com/bdreece/hopper/pkg/proto/grpc"
	"google.golang.org/grpc"
)

func NewServer(cfg *config.Config) *grpc.Server {
	logger := cfg.Logger.WithContext("server")
	server := grpc.NewServer()

	logger.Infoln("Registering services...")
	pb.RegisterDeviceServiceServer(server, cfg.DeviceService)
	pb.RegisterDeviceModelServiceServer(server, cfg.DeviceModelService)
	pb.RegisterEventServiceServer(server, cfg.EventService)
	pb.RegisterFirmwareServiceServer(server, cfg.FirmwareService)
	pb.RegisterPropertyServiceServer(server, cfg.PropertyService)
	pb.RegisterTenantServiceServer(server, cfg.TenantService)
	pb.RegisterTypeServiceServer(server, cfg.TypeService)
	pb.RegisterUnitServiceServer(server, cfg.UnitService)

	logger.Infoln("Services registered")
	return server
}
