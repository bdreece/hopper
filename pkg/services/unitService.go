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
	"github.com/bdreece/hopper/pkg/services/utils/iter"
	"gorm.io/gorm"
)

var (
	ErrUnitNotFound = errors.New("unit not found")
	ErrUnitQuery    = errors.New("failed to query units")
)

type UnitService struct {
	db     *gorm.DB
	logger utils.Logger
	grpc.UnimplementedUnitServiceServer
}

func NewUnitService(cfg *config.Config) *UnitService {
	return &UnitService{
		db:     cfg.DB,
		logger: cfg.Logger,

		UnimplementedUnitServiceServer: grpc.UnimplementedUnitServiceServer{},
	}
}

func (s *UnitService) GetUnit(ctx context.Context, in *proto.GetUnitRequest) (*proto.Unit, error) {
	s.logger.Infoln("Querying unit...")

	query := s.db
	switch t := in.Where.(type) {
	case *proto.GetUnitRequest_Uuid:
		s.logger.Infoln("...with UUID")
		query = query.Where("uuid = ?", t.Uuid)

	case *proto.GetUnitRequest_Property:
		s.logger.Infoln("...by property with UUID")
		query = query.
			Joins("inner join type on unit.typeId = type.Id").
			Joins("inner join property on property.typeId = type.Id").
			Where("property.uuid = ?", t.Property.Uuid)

	case *proto.GetUnitRequest_Type:
		s.logger.Infoln("...by type")
		query = query.
			Joins("inner join type on unit.typeId = type.Id")

		switch u := t.Type.Where.(type) {
		case *proto.GetTypeRequest_Uuid:
			s.logger.Infoln("...with UUID")
			query = query.Where("type.uuid = ?", u.Uuid)

		case *proto.GetTypeRequest_Property:
			s.logger.Infoln("...by property with UUID")
			query = query.
				Joins("inner join property on property.typeId = type.Id").
				Where("property.uuid = ?", u.Property.Uuid)
		}
	}

	unit := models.Unit{}
	result := query.First(&unit)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Unit received!")
	return &unit.Unit, nil
}

func (s *UnitService) GetUnits(ctx context.Context, in *proto.GetUnitsRequest) (*proto.Units, error) {
	s.logger.Infoln("Querying units by tenant...")

	query := s.db.Joins("inner join tenant on unit.tenantId = tenant.Id")
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

	units := make([]models.Unit, 0, 1)
	result := query.Scan(&units)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Units received!")
	return &proto.Units{
		Units: iter.Collect(iter.NewMap(
			iter.FromSlice(&units),
			func(in *models.Unit) *proto.Unit {
				return &in.Unit
			},
		)),
	}, nil
}

func (s *UnitService) handleError(err error) error {
	s.logger.Errorf("An error has occurred: %v\n", err)
	return err
}

func (s *UnitService) handleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(ErrUnitNotFound, err)
	}
	err = utils.WrapError(ErrUnitQuery, err)
	return s.handleError(err)
}
