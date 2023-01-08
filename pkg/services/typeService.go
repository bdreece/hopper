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
	"github.com/bdreece/hopper/pkg/utils"
	"github.com/bdreece/hopper/pkg/utils/iter"
	"gorm.io/gorm"
)

var (
	ErrTypeNotFound = errors.New("type not found")
	ErrTypeQuery    = errors.New("failed to query types")
)

type TypeService struct {
	db     *gorm.DB
	logger utils.Logger

	grpc.UnimplementedTypeServiceServer
}

func NewTypeService(cfg *config.Config) *TypeService {
	return &TypeService{
		db:     cfg.DB,
		logger: cfg.Logger,
	}
}

func (s *TypeService) GetType(ctx context.Context, in *proto.GetTypeRequest) (*proto.Type, error) {
	s.logger.Infoln("Querying type...")

	query := s.db
	switch t := in.Where.(type) {
	case *proto.GetTypeRequest_Uuid:
		s.logger.Infoln("...with UUID")
		query = query.
			Where("uuid = ?", t.Uuid)

	case *proto.GetTypeRequest_Property:
		s.logger.Infoln("...by property with UUID")
		query = query.
			Joins("inner join property on property.typeId = type.Id").
			Where("property.uuid = ?", t.Property.Uuid)
	}

	typeModel := models.Type{}
	result := query.First(&typeModel)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Type received!")
	return &typeModel.Type, nil
}

func (s *TypeService) GetTypes(ctx context.Context, in *proto.GetTypesRequest) (*proto.Types, error) {
	s.logger.Infoln("Querying types by tenant...")

	query := s.db.Joins("inner join tenant on type.tenantId = tenant.Id")
	switch t := in.Where.Tenant.Where.(type) {
	case *proto.GetTenantRequest_Uuid:
		s.logger.Infoln("...with UUID")
		query = query.Where("tenant.uuid = ?", t.Uuid)

	case *proto.GetTenantRequest_Device:
		s.logger.Infoln("...by device with UUID")
		query = query.
			Joins("inner join device on device.tenantId = tenant.Id").
			Where("device.uuid = ?", t.Device.Uuid)
	}

	types := make([]models.Type, 0, 1)
	result := query.Scan(&types)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	return &proto.Types{
		Types: iter.Collect(iter.NewMap(
			iter.FromSlice(&types),
			func(in *models.Type) *proto.Type {
				return &in.Type
			},
		)),
	}, nil
}

func (s *TypeService) handleError(err error) error {
	s.logger.Errorf("An error has occurred: %v\n", err)
	return err
}

func (s *TypeService) handleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(ErrTypeNotFound, err)
	}
	err = utils.WrapError(ErrTypeQuery, err)
	return s.handleError(err)
}
