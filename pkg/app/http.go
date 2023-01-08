package app

import (
	"net/http"

	"github.com/bdreece/hopper/pkg/config"
)

func NewHttpServer(cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:    cfg.GraphQLPort,
		Handler: cfg.GraphQLHandler,
	}
}
