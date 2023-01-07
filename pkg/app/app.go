package app

import (
	"net"

	"github.com/bdreece/hopper/pkg/config"
	"github.com/bdreece/hopper/pkg/services"
	"google.golang.org/grpc"
)

type App struct {
	srv    *grpc.Server
	config *config.Config
}

func NewApp() (*App, error) {
	builder := config.NewConfigBuilder().
		AddCredentials()

	cfg := builder.Build()

	db, err := NewDB(cfg)
	if err != nil {
		return nil, err
	}

	cfg = builder.
		AddDatabase(db).
		Build()

	cfg = builder.
		AddLogger().
		AddDeviceService(services.NewDeviceService(cfg)).
		AddEventService(services.NewEventService(cfg)).
		AddFirmwareService(services.NewFirmwareService(cfg)).
		AddPort(":8080").
		Build()

	srv := NewServer(cfg)
	return &App{
		srv:    srv,
		config: cfg,
	}, nil
}

func (a *App) Serve() error {
	listener, err := net.Listen("tcp", a.config.Port)
	if err != nil {
		return err
	}

	return a.srv.Serve(listener)
}
