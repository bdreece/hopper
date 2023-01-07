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
