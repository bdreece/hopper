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
	ErrDeviceModelQuery = errors.New("failed to query device models")
)

type DeviceModelResolverService struct {
	query  *gorm.DB
	logger utils.Logger

	*deviceModelResolver
}

func NewDeviceModelResolverService(cfg *config.Config) *DeviceModelResolverService {
	return &DeviceModelResolverService{
		query:  cfg.DB.InnerJoins("deviceModel"),
		logger: cfg.Logger.WithContext("DeviceModel"),
	}
}

func (s *DeviceModelResolverService) Tenant(ctx context.Context, model *models.DeviceModel) (*models.Tenant, error) {
	logger := s.logger.WithContext("Tenant")
	logger.Infoln("Querying device model tenant...")

	tenant := models.Tenant{}
	err := s.query.
		Where("deviceModel.Id = ?", model.ID).
		First(&tenant).
		Error

	if err != nil {
		return nil, s.handleQueryError(err, "Tenant")
	}

	logger.Infoln("Tenant received!")
	return &tenant, nil
}

func (s *DeviceModelResolverService) handleError(err error) error {
	s.logger.Errorf("An error has occurred: %v\n", err)
	return err
}

func (s *DeviceModelResolverService) handleQueryError(err error, field string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(fmt.Errorf("%s not found", field), err)
	}

	err = utils.WrapError(ErrDeviceModelQuery, err)
	return s.handleError(err)
}
