package services

import (
	"context"
	"errors"
	"log"

	"github.com/bdreece/hopper/pkg/config"
	"github.com/bdreece/hopper/pkg/models"
	"github.com/bdreece/hopper/pkg/proto"
	"github.com/bdreece/hopper/pkg/proto/grpc"
	"gorm.io/gorm"
)

type FirmwareService struct {
	db     *gorm.DB
	logger *log.Logger
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
	var result *gorm.DB = nil
	firmware := models.Firmware{}

	switch t := in.GetWhere().(type) {
	case *proto.GetFirmwareRequest_Uuid:
		result = s.db.
			Where("uuid = ?", t.Uuid).
			First(&firmware)
	case *proto.GetFirmwareRequest_Version:
		result = s.db.
			Where("modelId = ? AND versionMajor = ? AND versionMinor = ? AND versionPatch = ?",
				t.Version.ModelId, t.Version.VersionMajor,
				t.Version.VersionMinor, t.Version.VersionPatch).
			First(&firmware)
	}

	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("Firmware not found")
	}

	return &firmware.Firmware, nil
}
