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
	"fmt"

	"github.com/bdreece/hopper/pkg/config"
	"github.com/bdreece/hopper/pkg/models"
	"github.com/bdreece/hopper/pkg/proto"
	"github.com/bdreece/hopper/pkg/proto/grpc"
	"github.com/bdreece/hopper/pkg/services/utils"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

var (
	ErrCreateToken    = errors.New("Failed to create access token")
	ErrDeviceNotFound = errors.New("Device not found")
	ErrDeviceQuery    = errors.New("Failed to query device")
	ErrDecodeApiKey   = errors.New("Failed to decode API key")
)

type DeviceService struct {
	db     *gorm.DB
	logger utils.Logger
	secret string

	grpc.UnimplementedDeviceServiceServer
}

func NewDeviceService(cfg *config.Config) *DeviceService {
	return &DeviceService{
		db:     cfg.DB,
		logger: cfg.Logger.WithContext("DeviceService"),
		secret: cfg.Secret,

		UnimplementedDeviceServiceServer: grpc.UnimplementedDeviceServiceServer{},
	}
}

func (s *DeviceService) AuthDevice(ctx context.Context, in *proto.AuthDeviceRequest) (*proto.AuthDeviceResponse, error) {
	logger := s.logger.WithContext("AuthDevice")
	logger.Infoln("Decoding API key...")
	id, tenantId, err := utils.DecodeApiKey(in.ApiKey, s.secret)
	if err != nil {
		err = utils.WrapError(ErrDecodeApiKey, err)
		return nil, s.handleError(err)
	}

	logger.Infoln("Querying device...")
	device := models.Device{}
	result := s.db.
		Where("id = ?", id).
		First(&device)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, s.handleQueryError(result.Error)
		}

		logger.Warnln("Device does not exist!")
		logger.Infoln("Creating new device...")
		device = models.Device{
			Device: proto.Device{
				Uuid:     uuid.NewString(),
				TenantId: uint32(*tenantId),
			},
		}

		result = s.db.Save(&device)
		if result.Error != nil {
			return nil, s.handleQueryError(result.Error)
		}

		logger.Infoln("Device created!")
	}

	logger.Infoln("Creating access token...")
	jwt, expiration, err := utils.CreateToken(fmt.Sprint(id), fmt.Sprint(tenantId), s.secret)
	if err != nil {
		err = utils.WrapError(ErrCreateToken, err)
		return nil, s.handleError(err)
	}

	logger.Infoln("Device authenticated!")
	return &proto.AuthDeviceResponse{
		Token:      *jwt,
		Expiration: timestamppb.New(*expiration),
	}, nil
}

func (s *DeviceService) GetDevice(ctx context.Context, in *proto.GetDeviceRequest) (*proto.Device, error) {
	logger := s.logger.WithContext("GetDevice")
	logger.Infoln("Querying device by UUID...")
	device := models.Device{}
	result := s.db.
		Where("uuid = ?", in.Uuid).
		First(&device)

	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Device received!")
	return &device.Device, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, in *proto.UpdateDeviceRequest) (*proto.Device, error) {
	logger := s.logger.WithContext("UpdateDevice")
	logger.Infoln("Querying device by UUID...")

	device := models.Device{}
	result := s.db.
		Where("uuid = ?", in.Where.Uuid).
		First(&device)

	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Updating device...")
	device.Update(in)
	s.db.Save(&device)

	logger.Infoln("Device updated!")
	return &device.Device, nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, in *proto.DeleteDeviceRequest) (*proto.Device, error) {
	logger := s.logger.WithContext("DeleteDevice")
	logger.Infoln("Querying device by UUID...")

	device := models.Device{}
	result := s.db.
		Where("uuid = ?", in.GetUuid()).
		First(&device)

	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Deleting device...")
	result = s.db.Delete(&device)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Device deleted")
	return &device.Device, nil
}

func (s *DeviceService) handleError(err error) error {
	s.logger.Errorf("An error occurred: %v", err)
	return err
}

func (s *DeviceService) handleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(ErrDeviceNotFound, err)
	}
	err = utils.WrapError(ErrDeviceQuery, err)
	return s.handleError(err)
}
