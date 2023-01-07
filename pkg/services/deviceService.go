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
		logger: cfg.Logger,
		secret: cfg.Secret,

		UnimplementedDeviceServiceServer: grpc.UnimplementedDeviceServiceServer{},
	}
}

func (s *DeviceService) AuthDevice(ctx context.Context, in *proto.AuthDeviceRequest) (*proto.AuthDeviceResponse, error) {
	s.logger.Infoln("Authenticating device...")

	s.logger.Infoln("Decoding API key...")
	id, tenantId, err := utils.DecodeApiKey(in.ApiKey, s.secret)
	if err != nil {
		err = utils.WrapError(ErrDecodeApiKey, err)
		return nil, s.handleError(err)
	}

	s.logger.Infoln("Querying device...")
	device := models.Device{}
	result := s.db.
		Where("id = ?", id).
		First(&device)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, s.handleQueryError(result.Error)
		}

		s.logger.Warnln("Device does not exist!")
		s.logger.Infoln("Creating new device...")
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

		s.logger.Infoln("Device created!")
	}

	s.logger.Infoln("Creating access token...")
	jwt, expiration, err := utils.CreateToken(fmt.Sprint(id), fmt.Sprint(tenantId), s.secret)
	if err != nil {
		err = utils.WrapError(ErrCreateToken, err)
		return nil, s.handleError(err)
	}

	s.logger.Infoln("Device authenticated!")
	return &proto.AuthDeviceResponse{
		Token:      *jwt,
		Expiration: timestamppb.New(*expiration),
	}, nil
}

func (s *DeviceService) GetDevice(ctx context.Context, in *proto.GetDeviceRequest) (*proto.Device, error) {
	s.logger.Infoln("Querying device by UUID...")
	device := models.Device{}
	result := s.db.
		Where("uuid = ?", in.Uuid).
		First(&device)

	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Device received!")
	return &device.Device, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, in *proto.UpdateDeviceRequest) (*proto.Device, error) {
	s.logger.Infoln("Querying device by UUID...")
	device := models.Device{}
	result := s.db.
		Where("uuid = ?", in.Where.Uuid).
		First(&device)

	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Updating device...")
	device.Update(in)
	s.db.Save(&device)

	s.logger.Infoln("Device updated!")
	return &device.Device, nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, in *proto.DeleteDeviceRequest) (*proto.Device, error) {
	s.logger.Infoln("Querying device by UUID...")
	device := models.Device{}
	result := s.db.
		Where("uuid = ?", in.GetUuid()).
		First(&device)

	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Deleting device...")
	result = s.db.Delete(&device)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Device deleted")
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
