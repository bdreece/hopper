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
	"context"
	"errors"
	"net"

	"github.com/bdreece/hopper/pkg/config"
	"github.com/bdreece/hopper/pkg/graphql"
	"github.com/bdreece/hopper/pkg/services"
	"github.com/bdreece/hopper/pkg/services/utils"
	"google.golang.org/grpc"
)

var (
	ErrAppStartup = errors.New("failed to start application")
)

type App struct {
	Logger utils.Logger

	server *grpc.Server
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
		return nil, handleError(err, logger)
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
		AddGraphQLServer(graphql.NewGraphQLServer(cfg)).
		AddGraphQLPort(":8080").
		AddGrpcPort(":8081").
		Build()

	logger.Infoln("Creating server...")
	srv := NewServer(cfg)

	logger.Infoln("Application built!")
	return &App{
		Logger: logger,
		server: srv,
		config: cfg,
	}, nil
}

func (a *App) Serve() error {
	a.Logger.Infof("gRPC listening on port %s\n", a.config.GrpcPort)
	a.Logger.Infof("GraphQL listening on port %s\n", a.config.GraphQLPort)
	listener, err := net.Listen("tcp", a.config.GrpcPort)
	if err != nil {
		return handleError(err, a.Logger)
	}

	a.Logger.Infoln("Starting application...")
	return a.server.Serve(listener)
}

func (a *App) Shutdown(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		a.server.GracefulStop()
		cancel()
	}()
	<-ctx.Done()
	a.server.Stop()
}

func handleError(err error, logger utils.Logger) error {
	err = utils.WrapError(ErrAppStartup, err)
	logger.Errorf("An error occurred: %v\n", err)
	return err
}
