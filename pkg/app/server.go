package app

import (
	pb "github.com/bdreece/hopper/pkg/proto/grpc"
	"github.com/bdreece/hopper/pkg/services"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB, secret string) *grpc.Server {
	server := grpc.NewServer()

	deviceService := services.NewDeviceService(db, secret)
	pb.RegisterDeviceServiceServer(server, deviceService)

	eventService := services.NewEventService(db)
	pb.RegisterEventServiceServer(server, eventService)

	return server
}
