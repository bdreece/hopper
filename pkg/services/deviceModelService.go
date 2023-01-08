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
	ErrDeviceModelNotFound = errors.New("device model not found")
	ErrDeviceModelQuery    = errors.New("failed to query device model")
)

type DeviceModelService struct {
	db     *gorm.DB
	logger utils.Logger
	grpc.UnimplementedDeviceModelServiceServer
}

func NewDeviceModelService(cfg *config.Config) *DeviceModelService {
	return &DeviceModelService{
		db:     cfg.DB,
		logger: cfg.Logger.WithContext("DeviceModelService"),

		UnimplementedDeviceModelServiceServer: grpc.UnimplementedDeviceModelServiceServer{},
	}
}

func (s *DeviceModelService) GetDeviceModel(ctx context.Context, in *proto.GetDeviceModelRequest) (*proto.DeviceModel, error) {
	logger := s.logger.WithContext("GetDeviceModel")
	logger.Infoln("Querying device model...")

	query := s.db
	switch t := in.Where.(type) {
	case *proto.GetDeviceModelRequest_Device:
		logger.Infoln("...by device")
		query = query.
			Joins("inner join deviceModel on device.modelId = deviceModel.Id").
			Where("device.Uuid = ?", t.Device.GetUuid())

	case *proto.GetDeviceModelRequest_Uuid:
		logger.Infoln("...by UUID")
		query = query.Where("uuid = ?", t.Uuid)
	}

	deviceModel := models.DeviceModel{}
	result := query.First(&deviceModel)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Device model received!")
	return &deviceModel.DeviceModel, nil
}

func (s *DeviceModelService) GetDeviceModels(ctx context.Context, in *proto.GetDeviceModelsRequest) (*proto.DeviceModels, error) {
	logger := s.logger.WithContext("GetDeviceModels")
	logger.Infoln("Querying device models...")

	query := s.db.Joins("inner join tenant on deviceModel.tenantId = tenant.Id")
	switch t := in.Where.Tenant.Where.(type) {
	case *proto.GetTenantRequest_Uuid:
		logger.Infoln("...with UUID")
		query = query.Where("tenant.uuid = ?", t.Uuid)

	case *proto.GetTenantRequest_Device:
		logger.Infoln("...by device with UUID")
		query = query.
			Joins("inner join device on device.tenantId = tenant.Id").
			Where("device.uuid = ?", t.Device.Uuid)
	}

	deviceModels := make([]models.DeviceModel, 0, 1)
	result := query.First(&deviceModels)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Device models received!")
	return &proto.DeviceModels{
		Models: iter.Collect(iter.NewMap(
			iter.FromSlice(&deviceModels),
			func(in *models.DeviceModel) *proto.DeviceModel {
				return &in.DeviceModel
			})),
	}, nil
}

func (s *DeviceModelService) handleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(ErrDeviceModelNotFound, err)
	}
	err = utils.WrapError(ErrDeviceModelQuery, err)
	s.logger.Errorf("An error occurred: %v\n", err)
	return err
}
