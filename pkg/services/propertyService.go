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
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPropertyNotFound = errors.New("property not found")
	ErrPropertyQuery    = errors.New("failed to query properties")
)

type PropertyService struct {
	db     *gorm.DB
	logger utils.Logger
	grpc.UnimplementedPropertyServiceServer
}

func NewPropertyService(cfg *config.Config) *PropertyService {
	return &PropertyService{
		db:     cfg.DB,
		logger: cfg.Logger.WithContext("PropertyService"),

		UnimplementedPropertyServiceServer: grpc.UnimplementedPropertyServiceServer{},
	}
}

func (s *PropertyService) CreateProperty(ctx context.Context, in *proto.CreatePropertyRequest) (*proto.Property, error) {
	logger := s.logger.WithContext("CreateProperty")
	logger.Infoln("Creating property...")

	property := models.Property{
		Property: proto.Property{
			Uuid:        uuid.NewString(),
			Name:        in.GetName(),
			Description: in.Description,
			TypeId:      in.GetTypeId(),
		},
	}

	logger.Infoln("Saving property to database...")
	result := s.db.Create(&property)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Property saved!")
	return &property.Property, nil
}

func (s *PropertyService) GetProperty(ctx context.Context, in *proto.GetPropertyRequest) (*proto.Property, error) {
	logger := s.logger.WithContext("GetProperty")
	logger.Infoln("Querying property by UUID...")

	property := models.Property{}
	result := s.db.
		Where("uuid = ?", in.Uuid).
		First(&property)

	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Property received!")
	return &property.Property, nil
}

func (s *PropertyService) GetProperties(ctx context.Context, in *proto.GetPropertiesRequest) (*proto.Properties, error) {
	logger := s.logger.WithContext("GetProperties")
	logger.Infoln("Querying properties...")

	query := s.db
	switch t := in.Where.(type) {
	case *proto.GetPropertiesRequest_Device:
		logger.Infoln("...by device with UUID")
		query = query.
			Joins("inner join device on property.deviceId = device.Id").
			Where("device.uuid = ?", t.Device.Uuid)

	case *proto.GetPropertiesRequest_Model:
		logger.Infoln("...by device model")
		query = query.Joins("inner join model on property.modelId = model.Id")

		switch u := t.Model.Where.(type) {
		case *proto.GetDeviceModelRequest_Uuid:
			logger.Infoln("...with UUID")
			query = query.Where("model.uuid = ?", u.Uuid)

		case *proto.GetDeviceModelRequest_Device:
			logger.Infoln("...by device with UUID")
			query = query.
				Joins("inner join device on device.modelId = model.Id").
				Where("device.uuid = ?", u.Device.Uuid)
		}
	}

	properties := make([]models.Property, 0, 1)
	result := query.Scan(&properties)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Properties received!")
	return &proto.Properties{
		Properties: iter.Collect(iter.NewMap(
			iter.FromSlice(&properties),
			func(in *models.Property) *proto.Property {
				return &in.Property
			})),
	}, nil
}

func (s *PropertyService) UpdateProperty(ctx context.Context, in *proto.UpdatePropertyRequest) (*proto.Property, error) {
	logger := s.logger.WithContext("UpdateProperty")
	logger.Infoln("Querying property...")

	property := models.Property{}
	result := s.db.
		Where("uuid = ?", in.Where.Uuid).
		First(&property)

	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Updating property...")
	property.Update(in)
	s.db.Save(&property)

	logger.Infoln("Property updated!")
	return &property.Property, nil
}

func (s *PropertyService) DeleteProperty(ctx context.Context, in *proto.DeletePropertyRequest) (*proto.Property, error) {
	logger := s.logger.WithContext("DeleteProperty")
	logger.Infoln("Querying property with UUID...")

	property := models.Property{}
	result := s.db.
		Where("uuid = ?", in.Where.Uuid).
		First(&property)

	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Deleting property...")
	result = s.db.Delete(&property)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Property deleted!")
	return &property.Property, nil
}

func (s *PropertyService) handleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(ErrPropertyNotFound, err)
	}
	err = utils.WrapError(ErrPropertyQuery, err)
	s.logger.Errorf("An error occurred: %v\n", err)
	return err
}
