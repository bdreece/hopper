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
