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
	ErrTenantNotFound = errors.New("Tenant not found")
	ErrTenantQuery    = errors.New("Failed to query tenants")
)

type TenantService struct {
	db     *gorm.DB
	logger utils.Logger

	grpc.UnimplementedTenantServiceServer
}

func NewTenantService(cfg *config.Config) *TenantService {
	return &TenantService{
		db:     cfg.DB,
		logger: cfg.Logger,

		UnimplementedTenantServiceServer: grpc.UnimplementedTenantServiceServer{},
	}
}

func (s *TenantService) GetTenant(ctx context.Context, in *proto.GetTenantRequest) (*proto.Tenant, error) {
	s.logger.Infoln("Querying tenant...")

	query := s.db
	switch t := in.Where.(type) {
	case *proto.GetTenantRequest_Uuid:
		s.logger.Infoln("...with UUID")
		query = query.Where("uuid = ?", t.Uuid)

	case *proto.GetTenantRequest_Device:
		s.logger.Infoln("...by device with UUID")
		query = query.
			Joins("inner join device on device.tenantId = tenant.Id").
			Where("device.uuid = ?", t.Device.Uuid)
	}

	tenant := models.Tenant{}
	result := query.First(&tenant)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Tenant received!")
	return &tenant.Tenant, nil
}

func (s *TenantService) handleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(ErrTenantNotFound, err)
	}
	err = utils.WrapError(ErrTenantQuery, err)
	s.logger.Errorf("An error has occurred: %v\n", err)
	return err
}
