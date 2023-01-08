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
	ErrFirmwareNotFound = errors.New("firmware not found")
	ErrFirmwareQuery    = errors.New("failed to query firmwares")
)

type FirmwareService struct {
	db     *gorm.DB
	logger utils.Logger
	grpc.UnimplementedFirmwareServiceServer
}

func NewFirmwareService(cfg *config.Config) *FirmwareService {
	return &FirmwareService{
		db:     cfg.DB,
		logger: cfg.Logger.WithContext("FirmwareService"),

		UnimplementedFirmwareServiceServer: grpc.UnimplementedFirmwareServiceServer{},
	}
}

func (s *FirmwareService) GetFirmware(ctx context.Context, in *proto.GetFirmwareRequest) (*proto.Firmware, error) {
	logger := s.logger.WithContext("GetFirmware")
	logger.Infoln("Querying firmware...")

	query := s.db
	switch t := in.Where.(type) {
	case *proto.GetFirmwareRequest_Uuid:
		logger.Infoln("...by UUID")
		query = query.Where("uuid = ?", t.Uuid)

	case *proto.GetFirmwareRequest_Version:
		logger.Infoln("...by device model")
		query = query.Joins("inner join deviceModel on firmware.modelId = deviceModel.Id")

		switch u := t.Version.Model.Where.(type) {
		case *proto.GetDeviceModelRequest_Uuid:
			logger.Infoln("...with UUID")
			query = query.Where("deviceModel.uuid = ?", u.Uuid)

		case *proto.GetDeviceModelRequest_Device:
			logger.Infoln("...by device with UUID")
			query = query.
				Joins("inner join device on device.modelId = deviceModel.Id").
				Where("device.uuid = ?", u.Device.Uuid)
		}

		logger.Infoln("...with version")
		query = query.Where("versionMajor = ? and versionMinor = ? and versionPatch = ?",
			t.Version.VersionMajor, t.Version.VersionMinor, t.Version.VersionPatch)
	}

	firmware := models.Firmware{}
	result := query.First(&firmware)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Firmware received")
	return &firmware.Firmware, nil
}

func (s *FirmwareService) handleError(err error) error {
	s.logger.Errorf("An error occurred: %v\n", err)
	return err
}

func (s *FirmwareService) handleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(ErrFirmwareNotFound, err)
	}
	err = utils.WrapError(ErrFirmwareQuery, err)
	return s.handleError(err)
}
