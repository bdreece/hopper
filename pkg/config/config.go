package config

import (
	"log"

	"github.com/bdreece/hopper/pkg/proto/grpc"
	"gorm.io/gorm"
)

type Config struct {
	DeviceService grpc.DeviceServiceServer
	EventService  grpc.EventServiceServer
	Logger        *log.Logger
	DB            *gorm.DB
	Port          string
	Hostname      string
	Username      string
	Password      string
	Secret        string
}
