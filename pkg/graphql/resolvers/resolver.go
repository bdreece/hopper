//go:generate go run github.com/99designs/gqlgen generate
package resolvers

import (
	"github.com/bdreece/hopper/pkg/config"
	"github.com/bdreece/hopper/pkg/services/utils"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db     *gorm.DB
	logger utils.Logger
}

func NewResolver(cfg *config.Config) *Resolver {
	return &Resolver{
		db:     cfg.DB,
		logger: cfg.Logger,
	}
}
