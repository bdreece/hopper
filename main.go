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

//go:generate $PROJECT_ROOT/pkg/tools/generate.sh
//go:generate go run github.com/99designs/gqlgen generate
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bdreece/hopper/pkg/app"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func start(a *app.App) {
	if err := a.Serve(); err != nil &&
		err != grpc.ErrServerStopped {
		a.Logger.Errorf("An error occurred: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	a, err := app.NewApp()
	if err != nil {
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go start(a)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	a.Shutdown(ctx, cancel)
}
