package app

import (
	"errors"
	"net"

	"github.com/bdreece/hopper/pkg/config"
	"github.com/bdreece/hopper/pkg/services"
	"github.com/bdreece/hopper/pkg/services/utils"
	"google.golang.org/grpc"
)

var (
	ErrAppStartup = errors.New("Failed to start application")
)

type App struct {
	server *grpc.Server
	logger utils.Logger
	config *config.Config
}

func NewApp() (*App, error) {
	builder := config.
		NewConfigBuilder().
		AddCredentials().
		AddLogger()

	cfg := builder.Build()
	logger := cfg.Logger.WithContext("app")
	logger.Infoln("Building application...")

	db, err := NewDB(cfg)
	if err != nil {
		err = utils.WrapError(ErrAppStartup, err)
		logger.Errorf("An error occurred: %v\n", err)
		return nil, err
	}

	logger.Infoln("Injecting database connection...")
	cfg = builder.
		AddDatabase(db).
		Build()

	logger.Infoln("Injecting services...")
	cfg = builder.
		AddDeviceService(services.NewDeviceService(cfg)).
		AddDeviceModelService(services.NewDeviceModelService(cfg)).
		AddEventService(services.NewEventService(cfg)).
		AddFirmwareService(services.NewFirmwareService(cfg)).
		AddPropertyService(services.NewPropertyService(cfg)).
		AddTenantService(services.NewTenantService(cfg)).
		AddTypeService(services.NewTypeService(cfg)).
		AddPort(":8080").
		Build()

	logger.Infoln("Creating server...")
	srv := NewServer(cfg)

	logger.Infoln("Application built!")
	return &App{
		server: srv,
		logger: logger,
		config: cfg,
	}, nil
}

func (a *App) Serve() error {
	a.logger.Infof("Listening on port %s\n", a.config.Port)
	listener, err := net.Listen("tcp", a.config.Port)
	if err != nil {
		err = utils.WrapError(ErrAppStartup, err)
		a.logger.Errorf("An error occurred: %v\n", err)
		return err
	}

	a.logger.Infoln("Starting application...")
	return a.server.Serve(listener)
}
