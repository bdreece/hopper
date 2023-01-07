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

package services

import (
	"context"
	"errors"

	"github.com/bdreece/hopper/pkg/config"
	"github.com/bdreece/hopper/pkg/models"
	"github.com/bdreece/hopper/pkg/proto"
	"github.com/bdreece/hopper/pkg/proto/grpc"
	"github.com/bdreece/hopper/pkg/services/utils"
	"gorm.io/gorm"
)

var (
	ErrTenantNotFound = errors.New("Tenant not found")
	ErrTenantQuery    = errors.New("Failed to query tenants")
)

type TenantService struct {
	db     *gorm.DB
	logger utils.Logger

	grpc.UnimplementedTenantServiceServer
}

func NewTenantService(cfg *config.Config) *TenantService {
	return &TenantService{
		db:     cfg.DB,
		logger: cfg.Logger,

		UnimplementedTenantServiceServer: grpc.UnimplementedTenantServiceServer{},
	}
}

func (s *TenantService) GetTenant(ctx context.Context, in *proto.GetTenantRequest) (*proto.Tenant, error) {
	s.logger.Infoln("Querying tenant...")

	query := s.db
	switch t := in.Where.(type) {
	case *proto.GetTenantRequest_Uuid:
		s.logger.Infoln("...with UUID")
		query = query.Where("uuid = ?", t.Uuid)

	case *proto.GetTenantRequest_Device:
		s.logger.Infoln("...by device with UUID")
		query = query.
			Joins("inner join device on device.tenantId = tenant.Id").
			Where("device.uuid = ?", t.Device.Uuid)
	}

	tenant := models.Tenant{}
	result := query.First(&tenant)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Tenant received!")
	return &tenant.Tenant, nil
}

func (s *TenantService) handleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(ErrTenantNotFound, err)
	}
	err = utils.WrapError(ErrTenantQuery, err)
	s.logger.Errorf("An error has occurred: %v\n", err)
	return err
}
