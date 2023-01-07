package config

import (
	"log"
	"os"

	"github.com/bdreece/hopper/pkg/proto/grpc"
	"gorm.io/gorm"
)

const (
	HOSTNAME = "HOPPER_HOSTNAME"
	USERNAME = "HOPPER_USERNAME"
	PASSWORD = "HOPPER_PASSWORD"
	SECRET   = "HOPPER_SECRET"
)

type ConfigBuilder struct {
	config *Config
}

func NewConfigBuilder() ConfigBuilder {
	return ConfigBuilder{
		config: new(Config),
	}
}

func (b ConfigBuilder) AddDatabase(db *gorm.DB) ConfigBuilder {
	b.config.DB = db
	return b
}

func (b ConfigBuilder) AddDeviceService(deviceService grpc.DeviceServiceServer) ConfigBuilder {
	b.config.DeviceService = deviceService
	return b
}

func (b ConfigBuilder) AddEventService(eventService grpc.EventServiceServer) ConfigBuilder {
	b.config.EventService = eventService
	return b
}

func (b ConfigBuilder) AddFirmwareService(firmwareService grpc.FirmwareServiceServer) ConfigBuilder {
	b.config.FirmwareService = firmwareService
	return b
}

func (b ConfigBuilder) AddCredentials() ConfigBuilder {
	b.config.Hostname = os.Getenv(SECRET)
	b.config.Username = os.Getenv(HOSTNAME)
	b.config.Password = os.Getenv(USERNAME)
	b.config.Secret = os.Getenv(PASSWORD)
	return b
}

func (b ConfigBuilder) AddPort(port string) ConfigBuilder {
	b.config.Port = port
	return b
}

func (b ConfigBuilder) AddLogger() ConfigBuilder {
	b.config.Logger = log.Default()
	return b
}

func (b ConfigBuilder) Build() *Config {
	return b.config
}
