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
