package services

import (
	"context"
	"errors"
	"fmt"

	. "github.com/bdreece/hopper/pkg/models"
	pb "github.com/bdreece/hopper/pkg/proto"
	"github.com/bdreece/hopper/pkg/proto/grpc"
	"github.com/bdreece/hopper/pkg/services/utils"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type DeviceService struct {
	grpc.UnimplementedDeviceServiceServer
	db     *gorm.DB
	secret string
}

func NewDeviceService(db *gorm.DB, secret string) *DeviceService {
	return &DeviceService{
		grpc.UnimplementedDeviceServiceServer{},
		db,
		secret,
	}
}

func (s *DeviceService) AuthDevice(ctx context.Context, in *pb.AuthDeviceRequest) (*pb.AuthDeviceResponse, error) {
	device := &Device{}

	id, tenantId, err := utils.DecodeApiKey(in.ApiKey, s.secret)
	if err != nil {
		return nil, err
	}

	result := s.db.Where("id = ?", id).First(device)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		device = &Device{
			Device: pb.Device{
				Uuid:     uuid.NewString(),
				TenantId: uint32(*tenantId),
			},
		}
		s.db.Save(device)
	}

	jwt, expiration, err := utils.CreateToken(fmt.Sprint(id), fmt.Sprint(tenantId), s.secret)
	if err != nil {
		return nil, err
	}

	return &pb.AuthDeviceResponse{
		Token:      *jwt,
		Expiration: timestamppb.New(*expiration),
	}, nil
}

func (s *DeviceService) GetDevice(ctx context.Context, in *pb.GetDeviceRequest) (*pb.Device, error) {
	device := &Device{}
	result := s.db.Where("uuid = ?", in.Uuid).First(device)

	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("Device not found!")
	}

	return &device.Device, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, in *pb.UpdateDeviceRequest) (*pb.Device, error) {
	device := &Device{}
	result := s.db.Where("uuid = ?", in.Where.Uuid).First(device)

	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("Device not found!")
	}

	device.Update(in)
	s.db.Save(device)

	return &device.Device, nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, in *pb.DeleteDeviceRequest) (*pb.Device, error) {
	device := &Device{}
	result := s.db.Where("uuid = ?", in.GetUuid()).Delete(device)

	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("Failed to delete device")
	}

	return &device.Device, nil
}
