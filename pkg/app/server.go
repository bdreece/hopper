package app

import (
	"github.com/bdreece/hopper/pkg/config"
	. "github.com/bdreece/hopper/pkg/proto/grpc"
	"google.golang.org/grpc"
)

func NewServer(cfg *config.Config) *grpc.Server {
	server := grpc.NewServer()

	RegisterDeviceServiceServer(server, cfg.DeviceService)
	RegisterEventServiceServer(server, cfg.EventService)
	RegisterFirmwareServiceServer(server, cfg.FirmwareService)

	return server
}
