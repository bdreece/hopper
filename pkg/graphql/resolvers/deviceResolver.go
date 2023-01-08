package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/bdreece/hopper/pkg/config"
	"github.com/bdreece/hopper/pkg/models"
	"github.com/bdreece/hopper/pkg/utils"
	"gorm.io/gorm"
)

var (
	ErrDeviceQuery = errors.New("failed to query device field")
)

type DeviceResolverService struct {
	query  *gorm.DB
	logger utils.Logger

	*deviceResolver
}

func NewDeviceResolverService(cfg *config.Config) *DeviceResolverService {
	return &DeviceResolverService{
		query:  cfg.DB.InnerJoins("device"),
		logger: cfg.Logger.WithContext("DeviceResolver"),
	}
}

func (s *DeviceResolverService) Firmware(ctx context.Context, device *models.Device) (*models.Firmware, error) {
	logger := s.logger.WithContext("Firmware")

	logger.Infoln("Querying device firmware...")
	firmware := models.Firmware{}
	err := s.query.
		Where("device.id = ?", device.ID).
		First(&firmware).
		Error

	if err != nil {
		return nil, s.handleError(err, "Firmware")
	}

	logger.Infoln("Firmware received!")
	return &firmware, nil
}

func (s *DeviceResolverService) Tenant(ctx context.Context, device *models.Device) (*models.Tenant, error) {
	logger := s.logger.WithContext("Tenant")
	logger.Infoln("Querying device tenant...")

	tenant := models.Tenant{}
	err := s.query.
		Where("device.id = ?", device.ID).
		First(&tenant).
		Error

	if err != nil {
		return nil, s.handleError(err, "Tenant")
	}

	logger.Infoln("Tenant received!")
	return &tenant, nil
}

func (s *DeviceResolverService) Model(ctx context.Context, device *models.Device) (*models.DeviceModel, error) {
	logger := s.logger.WithContext("Model")
	logger.Infoln("Querying device model...")

	model := models.DeviceModel{}
	err := s.query.
		Where("device.id = ?", device.ID).
		First(&model).
		Error

	if err != nil {
		return nil, s.handleError(err, "Model")
	}

	return &model, nil
}

func (s *DeviceResolverService) handleError(err error, field string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(fmt.Errorf("%s not found", field), err)
	}
	err = utils.WrapError(ErrDeviceQuery, err)
	s.logger.Errorf("An error has occurred: %v\n", err)
	return err
}
