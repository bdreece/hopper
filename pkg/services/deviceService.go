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

type DeviceService struct {
	db     *gorm.DB
	secret string
	grpc.UnimplementedDeviceServiceServer
}

func NewDeviceService(cfg *config.Config) *DeviceService {
	return &DeviceService{
		db:                               cfg.DB,
		secret:                           cfg.Secret,
		UnimplementedDeviceServiceServer: grpc.UnimplementedDeviceServiceServer{},
	}
}

func (s *DeviceService) AuthDevice(ctx context.Context, in *proto.AuthDeviceRequest) (*proto.AuthDeviceResponse, error) {
	device := models.Device{}

	id, tenantId, err := utils.DecodeApiKey(in.ApiKey, s.secret)
	if err != nil {
		return nil, err
	}

	result := s.db.
		Where("id = ?", id).
		First(&device)

	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		device = models.Device{
			Device: proto.Device{
				Uuid:     uuid.NewString(),
				TenantId: uint32(*tenantId),
			},
		}
		s.db.Save(&device)
	}

	jwt, expiration, err := utils.CreateToken(fmt.Sprint(id), fmt.Sprint(tenantId), s.secret)
	if err != nil {
		return nil, err
	}

	return &proto.AuthDeviceResponse{
		Token:      *jwt,
		Expiration: timestamppb.New(*expiration),
	}, nil
}

func (s *DeviceService) GetDevice(ctx context.Context, in *proto.GetDeviceRequest) (*proto.Device, error) {
	device := models.Device{}
	result := s.db.
		Where("uuid = ?", in.Uuid).
		First(&device)

	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("Device not found!")
	}

	return &device.Device, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, in *proto.UpdateDeviceRequest) (*proto.Device, error) {
	device := models.Device{}
	result := s.db.
		Where("uuid = ?", in.Where.Uuid).
		First(&device)

	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("Device not found!")
	}

	device.Update(in)
	s.db.Save(&device)

	return &device.Device, nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, in *proto.DeleteDeviceRequest) (*proto.Device, error) {
	device := models.Device{}
	result := s.db.
		Where("uuid = ?", in.GetUuid()).
		Delete(&device)

	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("Failed to delete device")
	}

	return &device.Device, nil
}
