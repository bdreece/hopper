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
	ErrFirmwareNotFound = errors.New("Firmware not found")
	ErrFirmwareQuery    = errors.New("Failed to query firmwares")
)

type FirmwareService struct {
	db     *gorm.DB
	logger utils.Logger
	grpc.UnimplementedFirmwareServiceServer
}

func NewFirmwareService(cfg *config.Config) *FirmwareService {
	return &FirmwareService{
		db:                                 cfg.DB,
		logger:                             cfg.Logger,
		UnimplementedFirmwareServiceServer: grpc.UnimplementedFirmwareServiceServer{},
	}
}

func (s *FirmwareService) GetFirmware(ctx context.Context, in *proto.GetFirmwareRequest) (*proto.Firmware, error) {
	s.logger.Infoln("Querying firmware...")

	query := s.db
	switch t := in.Where.(type) {
	case *proto.GetFirmwareRequest_Uuid:
		s.logger.Infoln("...by UUID")
		query = query.Where("uuid = ?", t.Uuid)

	case *proto.GetFirmwareRequest_Version:
		s.logger.Infoln("...by device model")
		query = query.Joins("inner join deviceModel on firmware.modelId = deviceModel.Id")

		switch u := t.Version.Model.Where.(type) {
		case *proto.GetDeviceModelRequest_Uuid:
			s.logger.Infoln("...with UUID")
			query = query.Where("deviceModel.uuid = ?", u.Uuid)

		case *proto.GetDeviceModelRequest_Device:
			s.logger.Infoln("...by device with UUID")
			query = query.
				Joins("inner join device on device.modelId = deviceModel.Id").
				Where("device.uuid = ?", u.Device.Uuid)
		}

		s.logger.Infoln("...with version")
		query = query.Where("versionMajor = ? AND versionMinor = ? AND versionPatch = ?",
			t.Version.VersionMajor, t.Version.VersionMinor, t.Version.VersionPatch)
	}

	firmware := models.Firmware{}
	result := query.First(&firmware)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Firmware received")
	return &firmware.Firmware, nil
}

func (s *FirmwareService) handleError(err error) error {
	s.logger.Errorf("An error occurred: %v", err)
	return err
}

func (s *FirmwareService) handleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(ErrFirmwareNotFound, err)
	}
	err = utils.WrapError(ErrFirmwareQuery, err)
	return s.handleError(err)
}
