package app

import (
	"net"
	"os"

	"google.golang.org/grpc"
)

const (
	HOSTNAME = "HOPPER_HOSTNAME"
	USERNAME = "HOPPER_USERNAME"
	PASSWORD = "HOPPER_PASSWORD"
	SECRET   = "HOPPER_SECRET"
)

type App struct {
	srv *grpc.Server
}

func NewApp() (*App, error) {
	secret := os.Getenv(SECRET)
	db, err := NewDB()
	if err != nil {
		return nil, err
	}

	srv := NewServer(db, secret)

	return &App{
		srv,
	}, nil
}

func (a *App) Serve(port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	return a.srv.Serve(listener)
}
