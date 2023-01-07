package services

import (
	"context"
	"errors"

	"github.com/bdreece/hopper/pkg/config"
	"github.com/bdreece/hopper/pkg/models"
	"github.com/bdreece/hopper/pkg/proto"
	"github.com/bdreece/hopper/pkg/proto/grpc"
	"github.com/bdreece/hopper/pkg/services/utils"
	"github.com/bdreece/hopper/pkg/services/utils/iter"
	"gorm.io/gorm"
)

var (
	ErrCreateEvent     = errors.New("Failed creating event")
	ErrEventNotFound   = errors.New("Event not found")
	ErrEventQuery      = errors.New("Failed to query event")
	ErrMissingDeviceId = errors.New("Missing device ID")
)

type EventService struct {
	db     *gorm.DB
	logger utils.Logger
	grpc.UnimplementedEventServiceServer
}

func NewEventService(cfg *config.Config) *EventService {
	return &EventService{
		db:                              cfg.DB,
		logger:                          cfg.Logger,
		UnimplementedEventServiceServer: grpc.UnimplementedEventServiceServer{},
	}
}

func (s *EventService) CreateEvents(ctx context.Context, in *proto.CreateEventsRequest) (*proto.Events, error) {
	s.logger.Infoln("Resolving device ID from context...")
	deviceId, ok := ctx.Value("deviceId").(uint32)
	if !ok {
		return nil, s.handleError(ErrMissingDeviceId)
	}

	s.logger.Infoln("Creating events...")
	events := make([]models.Event, 0, 1)
	eventMsgs := make([]*proto.Event, 0, 1)

	for _, eventRequest := range in.Events {
		event := models.NewEvent(deviceId, eventRequest)
		eventMsgs = append(eventMsgs, &event.Event)
		events = append(events, event)
	}

	s.logger.Infof("Saving %d events to database...", len(events))
	result := s.db.Create(&events)
	if result.Error != nil {
		result.Error = utils.WrapError(ErrCreateEvent, result.Error)
		return nil, s.handleError(result.Error)
	}

	s.logger.Infoln("Events saved!")
	return &proto.Events{
		Events: eventMsgs,
	}, nil
}

func (s *EventService) GetEvent(ctx context.Context, in *proto.GetEventRequest) (*proto.Event, error) {
	s.logger.Infoln("Querying event...")

	query := s.db
	switch t := in.Where.(type) {
	case *proto.GetEventRequest_Uuid:
		s.logger.Infoln("...by UUID")
		query = query.
			Where("uuid = ?", t.Uuid)
	case *proto.GetEventRequest_DeviceTimestamp:
		s.logger.Infoln("...by device and timestamp")
		query = query.
			Joins("inner join device on event.deviceId = device.Id").
			Where("device.uuid = ? AND event.timestamp = ?",
				t.DeviceTimestamp.Device.Uuid,
				t.DeviceTimestamp.Timestamp)
	}

	event := models.Event{}
	result := query.First(&event)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Event received!")
	return &event.Event, nil
}

func (s *EventService) GetEvents(ctx context.Context, in *proto.GetEventsRequest) (*proto.Events, error) {
	s.logger.Infoln("Querying events by device UUID")
	events := make([]models.Event, 0, 1)
	result := s.db.
		Joins("inner join device on event.deviceId = device.Id").
		Where("device.deviceUuid = ?", in.Where.Device.Uuid).
		Scan(&events)

	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	s.logger.Infoln("Events received!")
	return &proto.Events{
		Events: iter.Collect(iter.NewMap(
			iter.FromSlice(&events),
			func(in *models.Event) *proto.Event {
				return &in.Event
			},
		)),
	}, nil
}

func (s *EventService) handleError(err error) error {
	s.logger.Errorf("An error occurred: %v", err)
	return err
}

func (s *EventService) handleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(ErrEventNotFound, err)
	}
	err = utils.WrapError(ErrEventQuery, err)
	return s.handleError(err)
}
