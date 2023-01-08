package app

import (
	"net/http"

	"github.com/bdreece/hopper/pkg/config"
)

func NewHttpServer(cfg *config.Config) *http.Server {
	return &http.Server{
		Handler: cfg.GraphQLHandler,
	}
}
