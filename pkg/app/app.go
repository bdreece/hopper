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

	db, err := NewDB(builder.Build())
	if err != nil {
		return nil, err
	}

	cfg := builder.
		AddDatabase(db).
		AddLogger().
		AddDeviceService(services.NewDeviceService(builder.Build())).
		AddEventService(services.NewEventService(builder.Build())).
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
