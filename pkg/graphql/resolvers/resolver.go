//go:generate go run github.com/99designs/gqlgen generate
package resolvers

import (
	"github.com/bdreece/hopper/pkg/config"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	cfg *config.Config
}

func NewResolver(cfg *config.Config) *Resolver {
	return &Resolver{cfg}
}

func (r *Resolver) Devices() DeviceResolver {
	return NewDeviceResolverService(r.cfg)
}

func (r *Resolver) DeviceModels() DeviceModelResolver {
	return NewDeviceModelResolverService(r.cfg)
}
