package services

import (
	"context"
	"errors"

	"github.com/bdreece/hopper/pkg/config"
	. "github.com/bdreece/hopper/pkg/models"
	pb "github.com/bdreece/hopper/pkg/proto"
	"github.com/bdreece/hopper/pkg/proto/grpc"
	"gorm.io/gorm"
)

type EventService struct {
	grpc.UnimplementedEventServiceServer
	db *gorm.DB
}

func NewEventService(cfg *config.Config) *EventService {
	return &EventService{
		UnimplementedEventServiceServer: grpc.UnimplementedEventServiceServer{},
		db:                              cfg.DB,
	}
}

func (e *EventService) CreateEvents(ctx context.Context, in *pb.CreateEventsRequest) (*pb.Events, error) {
	deviceId, ok := ctx.Value("deviceId").(uint32)
	if !ok {
		return nil, errors.New("Missing device ID")
	}

	events := make([]Event, 1)
	eventModels := make([]*pb.Event, 1)
	for _, eventRequest := range in.Events {
		event := NewEvent(deviceId, eventRequest)
		eventModels = append(eventModels, &event.Event)
		events = append(events, event)
	}

	result := e.db.Create(&events)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected != int64(len(eventModels)) {
		return nil, errors.New("Failed to create all events")
	}

	return &pb.Events{
		Events: eventModels,
	}, nil
}

func (e *EventService) GetEvent(ctx context.Context, in *pb.GetEventRequest) (*pb.Event, error) {
	var result *gorm.DB = nil
	event := &Event{}

	switch t := in.Where.(type) {
	case *pb.GetEventRequest_Uuid:
		result = e.db.
			Where("uuid = ?", t.Uuid).
			First(event)
	case *pb.GetEventRequest_DeviceTimestamp:
		result = e.db.
			Where("deviceId = ? AND timestamp = ?", t.DeviceTimestamp.DeviceId, t.DeviceTimestamp.Timestamp).
			First(event)
	}

	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("Event not found")
	}

	return &event.Event, nil
}

func (e *EventService) GetEvents(ctx context.Context, in *pb.GetEventsRequest) (*pb.Events, error) {
	events := make([]Event, 0, 1)
	result := e.db.
		Joins("inner join device on event.deviceId = device.Id").
		Where("device.deviceUuid = ?", in.Where.Device.Uuid).
		Scan(&events)

	if result.Error != nil {
		return nil, result.Error
	}

	eventModels := make([]*pb.Event, 1)
	for i := range events {
		eventModels = append(eventModels, &events[i].Event)
	}

	return &pb.Events{
		Events: eventModels,
	}, nil
}
